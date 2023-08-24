package service

import (
	"context"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-common/oss"
	"github.com/Happy-Why/toktik-common/utils"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/video"
	"github.com/Happy-Why/toktik-video/internal/dao/mysql"
	"github.com/Happy-Why/toktik-video/internal/dao/redis"
	"github.com/Happy-Why/toktik-video/internal/global"
	"github.com/Happy-Why/toktik-video/internal/model"
	"github.com/Happy-Why/toktik-video/internal/model/auto"
	"github.com/Happy-Why/toktik-video/internal/repo"
	"github.com/Happy-Why/toktik-video/pkg/myerr"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	cache       repo.Cache
	videoRepo   repo.VideoRepo
	rClient     repo.RClientRepo
	respRepo    repo.HandlerResp
	transaction mysql.Transaction
	ossRepo     oss.OSS
}

func NewVideoService() *VideoServiceImpl {
	return &VideoServiceImpl{
		cache:       redis.GetRdbCache(),
		videoRepo:   mysql.NewVideoDao(),
		rClient:     redis.NewVideoRClient(),
		respRepo:    repo.NewHandlerResps(),
		transaction: mysql.NewTransaction(),
		ossRepo:     global.OSS,
	}
}

func (vs *VideoServiceImpl) VideoFeed(ctx context.Context, req *video.VideoFeedRequest) (resp *video.VideoFeedResponse, err error) {
	resp = new(video.VideoFeedResponse)
	// 1.处理业务逻辑
	// 判断是否登录
	isLogin := req.UserId != 0
	// 2.业务处理
	// ① 按投稿时间查询最多30个video，按投稿时间降序
	videoInfos, err := vs.videoRepo.GetVideosByTime(ctx, *req.LatestTime*1e9)
	if err != nil {
		zap.L().Error("vs.videoRepo.GetVideosByTime err:", zap.Error(err))
		return vs.respRepo.VideoFeedResponse(errcode.ErrDB, err.Error(), &video.VideoFeedResponse{}), nil
	}
	// 获取next_time
	time := videoInfos[len(videoInfos)-1].CreatedAt.Unix()
	resp.NextTime = &time
	for _, v := range videoInfos {
		// ② 根据userId查询出 userInfo
		userInfo, err := vs.videoRepo.GetUserInfoByID(ctx, int64(v.UserId))
		if err != nil {
			zap.L().Error("vs.videoRepo.GetUserInfoByID err:", zap.Error(err))
			return vs.respRepo.VideoFeedResponse(errcode.ErrDB, err.Error(), &video.VideoFeedResponse{}), nil
		}
		isFollow := false
		isFavorite := false
		// ③ 若登录判断是否关注对方 //TODO 查询缓存
		if isLogin {
			isFollow, err = vs.videoRepo.IsFollowTargetUser(ctx, req.UserId, int64(v.UserId))
			if err != nil {
				zap.L().Error("vs.videoRepo.IsFollowTargetUser err:", zap.Error(err))
				return vs.respRepo.VideoFeedResponse(errcode.ErrDB, err.Error(), &video.VideoFeedResponse{}), nil
			}
			// ④ 判断该视频是否点赞
			isFavorite, _ = vs.videoRepo.IsFavoriteVideo(ctx, req.UserId, int64(v.ID))
			if err != nil {
				zap.L().Error("vs.videoRepo.IsFavoriteVideo err:", zap.Error(err))
				return vs.respRepo.VideoFeedResponse(errcode.ErrDB, err.Error(), &video.VideoFeedResponse{}), nil
			}
		}

		resp.VideoList = append(resp.VideoList, &video.Video{
			Id: int64(v.ID),
			Author: &video.User{
				Id:              int64(userInfo.ID),
				Name:            userInfo.Username,
				FollowCount:     &userInfo.FollowCount,
				FollowerCount:   &userInfo.FollowerCount,
				IsFollow:        isFollow,
				Avatar:          &userInfo.Avatar,
				BackgroundImage: &userInfo.BackgroundImage,
				Signature:       &userInfo.Signature,
				TotalFavorited:  &userInfo.TotalFavorited,
				WorkCount:       &userInfo.WorkCount,
				FavoriteCount:   &userInfo.FavoriteCount,
			},
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		})
	}

	return vs.respRepo.VideoFeedResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) VideoPublish(ctx context.Context, req *video.VideoPublishRequest) (resp *video.VideoPublishResponse, err error) {
	// 1.处理业务逻辑
	// bucket_name:why_bucket 视频目录：video/user_id/video_name 封面目录：cover/user_id/cover_name
	// video_name：video_id + title.mp4  cover：video_id + title.jpg
	// 使用雪花算法获得video_id，根据video_id和title获得视频文件和视频封面的URL

	// 2.处理业务,上传视频和视频封面
	videoId, _ := global.SnowFlake.NextId()
	// ① 构建 objectKey,也就是 目录+name
	videoDirectories := []string{model.VIDEO, strconv.FormatInt(req.UserId, 10), strconv.FormatInt(videoId, 10) + req.Title}
	videoObjectKey := vs.ossRepo.CreateObjectKey(model.MP4, videoDirectories...)
	coverDirectories := []string{model.COVER, strconv.FormatInt(req.UserId, 10), strconv.FormatInt(videoId, 10) + req.Title}
	coverObjectKey := vs.ossRepo.CreateObjectKey(model.JPG, coverDirectories...)
	// ② 上传视频文件,并返回 视频文件的 URL
	videoURL, err := vs.ossRepo.UploadByteFile(videoObjectKey, req.Data)
	if err != nil {
		zap.L().Error("vs.ossRepo.UploadByteFile videoURL err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(myerr.UploadVideoErr, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// ③ 根据视频的URL，获取 视频的第三帧作为封面
	reader, err := utils.ReadFrameAsJpeg(videoURL, 3)
	if err != nil {
		zap.L().Error("utils.ReadFrameAsJpeg err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(myerr.GetCoverErr, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// ④ 上传封面文件,并返回 封面文件的 URL
	coverURL, err := vs.ossRepo.UploadFile(coverObjectKey, reader)
	if err != nil {
		zap.L().Error("vs.ossRepo.UploadByteFile coverURL err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(myerr.UploadCoverErr, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// ⑤ 构建 video实例
	videoInfo := &auto.Video{
		Model:         gorm.Model{ID: uint(videoId)},
		UserId:        uint(req.UserId),
		Title:         req.Title,
		PlayURL:       videoURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	// ⑥ 将 video 信息存入 redis，并对 user 的 work_count+1
	videoKey := auto.NewVideoKey(uint(videoId))
	videoMap := auto.NewMapVideoInfo(videoInfo)
	if err = vs.rClient.HSetVideoInfo(ctx, videoKey, videoMap); err != nil {
		zap.L().Error("vs.rClient.HSetVideoInfo err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(errcode.ErrDB, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// 开启事务，将video实例存入数据库，并对 user 的 work_count+1
	if err = vs.transaction.Action(func(conn mysql.DbConn) error {
		if err = vs.videoRepo.CreateVideo(ctx, conn, videoInfo); err != nil {
			zap.L().Error("vs.videoRepo.CreateVideo err:", zap.Error(err))
			return err
		}
		if err = vs.videoRepo.AddUserWorkCount(ctx, conn, req.UserId); err != nil {
			zap.L().Error("vs.videoRepo.AddUserWorkCount err:", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return vs.respRepo.VideoPublishResponse(myerr.SubmitVideoErr, err.Error(), &video.VideoPublishResponse{}), nil
	}

	return vs.respRepo.VideoPublishResponse(errcode.StatusOK, model.MsgNil, &video.VideoPublishResponse{}), nil
}

func (vs *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	resp = new(video.PublishListResponse)
	// 1.处理业务逻辑
	// 2.处理业务
	// ① 根据 user_id 查询当前用户的所有信息
	userInfo, err := vs.videoRepo.GetUserInfoByID(ctx, req.UserId)
	if err != nil {
		zap.L().Error("vs.videoRepo.GetUserInfoByID err:", zap.Error(err))
		return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
	}
	// ② 根据 user_id 在 video表中查询 当前用户发表的视频
	videoInfos, err := vs.videoRepo.GetVideosByUserId(ctx, req.UserId)
	if err != nil {
		zap.L().Error("vs.videoRepo.GetVideosByUserId err:", zap.Error(err))
		return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
	}
	for _, v := range videoInfos {
		// ④ 判断该视频是否点赞
		isFavorite, _ := vs.videoRepo.IsFavoriteVideo(ctx, req.UserId, int64(v.ID))
		if err != nil {
			zap.L().Error("vs.videoRepo.IsFavoriteVideo err:", zap.Error(err))
			return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
		}
		resp.VideoList = append(resp.VideoList, &video.Video{
			Id: int64(v.ID),
			Author: &video.User{
				Id:              int64(userInfo.ID),
				Name:            userInfo.Username,
				FollowCount:     &userInfo.FollowCount,
				FollowerCount:   &userInfo.FollowerCount,
				IsFollow:        false,
				Avatar:          &userInfo.Avatar,
				BackgroundImage: &userInfo.BackgroundImage,
				Signature:       &userInfo.Signature,
				TotalFavorited:  &userInfo.TotalFavorited,
				WorkCount:       &userInfo.WorkCount,
				FavoriteCount:   &userInfo.FavoriteCount,
			},
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		})
	}
	return vs.respRepo.PublishListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) FavoriteAction(ctx context.Context, req *video.FavoriteActionRequest) (resp *video.FavoriteActionResponse, err error) {
	// 1.校验业务逻辑
	// ① 判断该视频是否存在
	ok, err := vs.videoRepo.IsVideoExist(ctx, req.VideoId)
	if err != nil {
		zap.L().Error("vs.videoRepo.IsVideoExist err:", zap.Error(err))
		return vs.respRepo.FavoriteActionResponse(errcode.ErrDB, err.Error(), &video.FavoriteActionResponse{}), nil
	}
	if !ok {
		return vs.respRepo.FavoriteActionResponse(myerr.VideoNotExist, err.Error(), &video.FavoriteActionResponse{}), nil
	}
	// ② 判断是否对该视频有点赞记录
	exist, err := vs.videoRepo.IsFavoriteRecordExist(ctx, req.UserId, req.VideoId)
	if err != nil {
		zap.L().Error("vs.videoRepo.IsFavoriteRecordExist err:", zap.Error(err))
		return vs.respRepo.FavoriteActionResponse(errcode.ErrDB, err.Error(), &video.FavoriteActionResponse{}), nil
	}
	// 2.处理业务
	//TODO 使用缓存计数
	// 构建一条favorite数据
	favoriteInfo := &auto.Favorite{
		Model:   gorm.Model{},
		UserId:  uint(req.UserId),
		VideoId: uint(req.VideoId),
	}
	switch req.ActionType {
	case model.FAVORITE:
		if exist {
			return vs.respRepo.FavoriteActionResponse(myerr.AlreadyFavorite, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// 使用事务，在 favorite表中添加点赞记录
		// 对当前用户的 favorite_count+1，对视频的作者的 total_favorite+1，对视频的 favorite_count+1
		if err = vs.transaction.Action(func(conn mysql.DbConn) error {
			if err = vs.videoRepo.CreateFavoriteRecord(ctx, conn, favoriteInfo); err != nil {
				zap.L().Error("vs.videoRepo.CreateFavoriteRecord err:", zap.Error(err))
				return err
			}
			if err = vs.videoRepo.AddUserFavoriteCount(ctx, conn, req.UserId); err != nil {
				zap.L().Error("vs.videoRepo.AddUserFavoriteCount err:", zap.Error(err))
				return err
			}
			if err = vs.videoRepo.AddVideoFavoriteCount(ctx, conn, req.VideoId); err != nil {
				zap.L().Error("vs.videoRepo.AddVideoFavoriteCount err:", zap.Error(err))
				return err
			}
			return nil
		}); err != nil {
			return vs.respRepo.FavoriteActionResponse(errcode.ErrDB, err.Error(), &video.FavoriteActionResponse{}), nil
		}
	case model.CANCELFAVORITE:
		if !exist {
			return vs.respRepo.FavoriteActionResponse(myerr.IsNotFavorite, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// 使用事务，在 favorite表中删除点赞记录,
		// 对当前用户的 favorite_count-1，对视频的作者的 total_favorite-1，对视频的 favorite_count-1
		if err = vs.transaction.Action(func(conn mysql.DbConn) error {
			if err = vs.videoRepo.DeleteFavoriteRecord(ctx, conn, favoriteInfo); err != nil {
				zap.L().Error("vs.videoRepo.DeleteFavoriteRecord err:", zap.Error(err))
				return err
			}
			if err = vs.videoRepo.SubUserFavoriteCount(ctx, conn, req.UserId); err != nil {
				zap.L().Error("vs.videoRepo.SubUserFavoriteCount err:", zap.Error(err))
				return err
			}
			if err = vs.videoRepo.SubVideoFavoriteCount(ctx, conn, req.VideoId); err != nil {
				zap.L().Error("vs.videoRepo.SubVideoFavoriteCount err:", zap.Error(err))
				return err
			}
			return nil
		}); err != nil {
			return vs.respRepo.FavoriteActionResponse(errcode.ErrDB, err.Error(), &video.FavoriteActionResponse{}), nil
		}
	}
	return vs.respRepo.FavoriteActionResponse(errcode.StatusOK, model.MsgNil, &video.FavoriteActionResponse{}), nil
}

func (vs *VideoServiceImpl) FavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	resp = new(video.FavoriteListResponse)
	// 1.校验业务逻辑
	// 2.处理业务
	// ① 根据 user_id 查询当前用户的所有信息
	userInfo, err := vs.videoRepo.GetUserInfoByID(ctx, req.UserId)
	if err != nil {
		zap.L().Error("vs.videoRepo.GetUserInfoByID err:", zap.Error(err))
		return vs.respRepo.FavoriteListResponse(errcode.ErrDB, err.Error(), &video.FavoriteListResponse{}), nil
	}
	// ② 根据 user_id 在 video表中查询 当前用户点赞的视频
	videoInfos, err := vs.videoRepo.GetFavoriteVideosByUserId(ctx, req.UserId)
	if err != nil {
		zap.L().Error("vs.videoRepo.GetVideosByUserId err:", zap.Error(err))
		return vs.respRepo.FavoriteListResponse(errcode.ErrDB, err.Error(), &video.FavoriteListResponse{}), nil
	}
	// 3.模型转换
	for _, v := range videoInfos {
		resp.VideoList = append(resp.VideoList, &video.Video{
			Id: int64(v.ID),
			Author: &video.User{
				Id:              int64(userInfo.ID),
				Name:            userInfo.Username,
				FollowCount:     &userInfo.FollowCount,
				FollowerCount:   &userInfo.FollowerCount,
				IsFollow:        false,
				Avatar:          &userInfo.Avatar,
				BackgroundImage: &userInfo.BackgroundImage,
				Signature:       &userInfo.Signature,
				TotalFavorited:  &userInfo.TotalFavorited,
				WorkCount:       &userInfo.WorkCount,
				FavoriteCount:   &userInfo.TotalFavorited,
			},
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    true,
			Title:         v.Title,
		})
	}
	// 4.返回数据
	return vs.respRepo.FavoriteListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) CommentAction(ctx context.Context, req *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) CommentList(ctx context.Context, req *video.CommentListRequest) (resp *video.CommentListResponse, err error) {
	return nil, nil
}
