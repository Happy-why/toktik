package service

import (
	"context"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-common/oss"
	"github.com/Happy-Why/toktik-common/utils"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/video"
	"github.com/Happy-Why/toktik-video/internal/dao"
	"github.com/Happy-Why/toktik-video/internal/dao/mysql"
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
	respRepo    repo.HandlerResp
	transaction mysql.Transaction
	ossRepo     oss.OSS
}

func NewVideoService() *VideoServiceImpl {
	return &VideoServiceImpl{
		cache:       dao.Group.Rdb,
		videoRepo:   mysql.NewVideoDao(),
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
				FavoriteCount:   &userInfo.TotalFavorited,
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
	// ⑤ 构建 video实例，存入数据库
	videoInfo := &auto.Video{
		Model:         gorm.Model{ID: uint(videoId)},
		UserId:        uint(req.UserId),
		Title:         req.Title,
		PlayURL:       videoURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	if err = vs.videoRepo.CreateVideo(ctx, videoInfo); err != nil {
		zap.L().Error("vs.ossRepo.UploadByteFile coverURL err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(myerr.SubmitVideoErr, err.Error(), &video.VideoPublishResponse{}), nil
	}
	return vs.respRepo.VideoPublishResponse(errcode.StatusOK, model.MsgNil, &video.VideoPublishResponse{}), nil
}

func (vs *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {

	return nil, nil
}

func (vs *VideoServiceImpl) FavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) FavoriteAction(ctx context.Context, req *video.FavoriteActionRequest) (resp *video.FavoriteActionResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) CommentAction(ctx context.Context, req *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) CommentList(ctx context.Context, req *video.CommentListRequest) (resp *video.CommentListResponse, err error) {
	return nil, nil
}
