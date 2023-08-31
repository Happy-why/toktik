package service

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
	"toktik-common/errcode"
	"toktik-common/oss"
	"toktik-common/utils"
	"toktik-rpc/kitex_gen/user"
	"toktik-rpc/kitex_gen/video"
	"toktik-video/internal/dao/mysql"
	"toktik-video/internal/dao/redis"
	"toktik-video/internal/global"
	"toktik-video/internal/model"
	"toktik-video/internal/model/auto"
	"toktik-video/internal/repo"
	"toktik-video/pkg/myerr"
	"toktik-video/pkg/rpc/client"
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
	// ① 先查询缓存，按投稿时间查询最多30个video，按投稿时间降序
	feedKey := auto.CreatePublishKey()
	nextTime, videoIds, err := vs.rClient.ZGetVideoIds(ctx, feedKey, *req.LatestTime)
	if videoIds == nil {
		return vs.respRepo.VideoFeedResponse(myerr.CanNotSearchVideo, model.MsgNil, &video.VideoFeedResponse{}), nil
	}
	if err != nil {
		zap.L().Error("vs.rClient.ZGetVideoIds err:", zap.Error(err))
		return vs.respRepo.VideoFeedResponse(errcode.ErrRedis, err.Error(), &video.VideoFeedResponse{}), nil
	}
	videoInfos := make([]*auto.Video, len(videoIds))
	userIds := make([]int64, len(videoIds))
	for i, videoId := range videoIds {
		// ② 根据 video_id 从缓存中查询 video_info
		videoInfo, err := vs.GetVideoInfo(ctx, videoId)
		if err != nil {
			return vs.respRepo.VideoFeedResponse(errcode.ErrDB, err.Error(), &video.VideoFeedResponse{}), nil
		}
		// 获取每个视频作者的user_id , video_info_user_id
		userIds[i] = int64(videoInfo.UserId)
		videoInfos[i] = videoInfo
	}
	// ③ 获取 nextTime
	resp.NextTime = &nextTime
	// ④ 根据user_ids，查询出所有用户信息
	userListResp, err := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{TargetId: userIds})
	if userListResp == nil {
		zap.L().Error("client.UserClient.GetUserList 返回空指针")
		return vs.respRepo.VideoFeedResponse(errcode.ErrServer, model.MsgNil, &video.VideoFeedResponse{}), nil
	}
	if userListResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.UserClient.GetUserList err:", zap.Error(err))
		return vs.respRepo.VideoFeedResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &video.VideoFeedResponse{}), nil
	}
	userListMap := make(map[uint]*user.User)
	for _, userInfo := range userListResp.UserList {
		userListMap[uint(userInfo.Id)] = userInfo
	}
	for _, v := range videoInfos {
		isFavorite := false
		if isLogin {
			isFavorite, err = vs.IsFavoriteVideo(ctx, req.UserId, int64(v.ID))
			if err != nil {
				return vs.respRepo.VideoFeedResponse(errcode.ErrDB, err.Error(), &video.VideoFeedResponse{}), nil
			}
		}
		resp.VideoList = append(resp.VideoList, &video.Video{
			Id:            int64(v.ID),
			Author:        userListMap[v.UserId],
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
	// 对于 用户和视频信息，缓存放到redis中，查询用户和视频信息先查缓存，缓存没有查数据库然后添加缓存
	// 一般像 用户和视频信息这样的 读多写少的信息，可以使用延迟双删策略，在更新信息操作时候进行缓存，进行缓存，但是本项目中，用户和视频信息都带有计数数据
	// 计数数据多为 读写相当，需要先更新redis然后定时更新数据库。
	// 在本项目中，可以将具体信息存一张表，对计数信息再存一张表，这样可以用两种策略对信息进行缓存。
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
	timeNow := time.Now()
	videoInfo := &auto.Video{
		BaseModel:     auto.BaseModel{ID: uint(videoId), CreatedAt: timeNow, UpdatedAt: timeNow},
		UserId:        uint(req.UserId),
		Title:         req.Title,
		PlayURL:       videoURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	// ⑦ 开启事务，将video实例存入数据库
	if err = vs.transaction.Action(func(conn mysql.DbConn) error {
		if err = vs.videoRepo.CreateVideo(ctx, conn, videoInfo); err != nil {
			zap.L().Error("vs.videoRepo.CreateVideo err:", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return vs.respRepo.VideoPublishResponse(myerr.SubmitVideoErr, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// ⑥ 将 video 信息存入 redis，并对 user 的 work_count+1
	videoKey := auto.CreateVideoKey(uint(videoId))
	videoMap := auto.CreateMapVideoInfo(videoInfo)
	if err = vs.rClient.HSetVideoInfo(ctx, videoKey, videoMap); err != nil {
		zap.L().Error("vs.rClient.HSetVideoInfo err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(errcode.ErrDB, err.Error(), &video.VideoPublishResponse{}), nil
	}
	addWorkCountResp, _ := client.UserClient.AddUserWorkCount(ctx, &user.AddUserWorkCountRequest{UserId: req.UserId})
	if addWorkCountResp == nil {
		zap.L().Error("client.UserClient.AddUserWorkCount 返回空指针")
		return vs.respRepo.VideoPublishResponse(errcode.ErrServer, model.MsgNil, &video.VideoPublishResponse{}), nil
	}
	if addWorkCountResp.StatusCode != 0 {
		zap.L().Error("client.UserClient.AddUserWorkCount err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(errcode.CreateErr(addWorkCountResp.StatusCode, model.MsgNil), addWorkCountResp.StatusMsg, &video.VideoPublishResponse{}), nil
	}
	// ⑦ 将video_id作为member，时间作为 score使用 zset 添加到 redis中
	publishKey := auto.CreatePublishKey()
	if err = vs.rClient.PublishVideo(ctx, publishKey, float64(timeNow.Unix()), videoId); err != nil {
		zap.L().Error("vs.rClient.PublishVideo err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(errcode.ErrRedis, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// ⑧ 删除用户的视频缓存，user_video::video_ids
	userVideoKey := auto.CreateUserVideoKey(uint(req.UserId))
	if err = vs.rClient.SDelUserVideo(ctx, userVideoKey); err != nil {
		zap.L().Error("vs.rClient.SDelUserVideo err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(errcode.ErrRedis, err.Error(), &video.VideoPublishResponse{}), nil
	}

	return vs.respRepo.VideoPublishResponse(errcode.StatusOK, model.MsgNil, &video.VideoPublishResponse{}), nil
}

func (vs *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	resp = new(video.PublishListResponse)
	// 1.处理业务逻辑
	// 2.处理业务
	// ① 根据 user_id 查询当前用户的所有信息
	userIndexResp, _ := client.UserClient.UserIndex(ctx, &user.UserIndexRequest{UserId: req.UserId, Token: "", MyUserId: req.UserId})
	if userIndexResp == nil {
		zap.L().Error("client.UserClient.UserIndex 返回空指针")
		return vs.respRepo.PublishListResponse(errcode.ErrServer, model.MsgNil, &video.PublishListResponse{}), nil
	}
	if userIndexResp.StatusCode != 0 {
		zap.L().Error("client.UserClient.UserIndex err:", zap.Error(err))
		return vs.respRepo.PublishListResponse(errcode.CreateErr(userIndexResp.StatusCode, model.MsgNil), userIndexResp.StatusMsg, &video.PublishListResponse{}), nil
	}

	// ② 根据 user_id 在 video表中查询 当前用户发表的视频信息，和视频id
	// 去缓存查询 user 的 video_ids
	userVideoKey := auto.CreateUserVideoKey(uint(req.UserId))
	videoIds, err := vs.rClient.SGetUserVideoIds(ctx, userVideoKey)
	if err != nil {
		zap.L().Error("vs.rClient.SGetUserVideoIds err:", zap.Error(err))
		return vs.respRepo.PublishListResponse(errcode.ErrRedis, err.Error(), &video.PublishListResponse{}), nil
	}
	if videoIds == nil {
		// 缓存没查到，去数据库查
		videoIds, err = vs.videoRepo.GetVideoIdsByUserId(ctx, req.UserId)
		if err != nil {
			zap.L().Error("vs.videoRepo.GetVideoIdsByUserId err:", zap.Error(err))
			return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
		}
		// 添加缓存
		if err = vs.rClient.SAddUserVideoIds(ctx, userVideoKey, videoIds); err != nil {
			zap.L().Error("vs.rClient.SAddUserVideoIds err:", zap.Error(err))
			return vs.respRepo.PublishListResponse(errcode.ErrRedis, err.Error(), &video.PublishListResponse{}), nil
		}
	}
	// 扩容
	resp.VideoList = make([]*video.Video, len(videoIds))
	// ③ 查找 videoInfos ，
	for i, videoId := range videoIds {
		// 根据 video_id 从缓存中查询 video_info
		videoInfo, err := vs.GetVideoInfo(ctx, videoId)
		if err != nil {
			return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
		}
		// ④ 判断视频是否点赞
		isFavorite, err := vs.IsFavoriteVideo(ctx, req.UserId, videoId)
		if err != nil {
			return vs.respRepo.PublishListResponse(errcode.ErrRedis, err.Error(), &video.PublishListResponse{}), nil
		}
		resp.VideoList[i] = &video.Video{
			Id:            videoId,
			Author:        userIndexResp.User,
			PlayUrl:       videoInfo.PlayURL,
			CoverUrl:      videoInfo.CoverURL,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoInfo.Title,
		}
	}
	return vs.respRepo.PublishListResponse(errcode.StatusOK, model.MsgNil, resp), nil

}

func (vs *VideoServiceImpl) FavoriteAction(ctx context.Context, req *video.FavoriteActionRequest) (resp *video.FavoriteActionResponse, err error) {
	videoKey := auto.CreateVideoKey(uint(req.VideoId))
	favKey := auto.CreateFavKey(uint(req.UserId))
	// 1.校验业务逻辑
	// ① 判断该视频是否存在,去redis查询video_info，并获取作者user_id，没有查到就去查数据库，数据库中没有返回错误，查到开启线程池添加缓存
	videoInfo, err := vs.GetVideoInfo(ctx, req.VideoId)
	if err == gorm.ErrRecordNotFound {
		return vs.respRepo.FavoriteActionResponse(myerr.VideoNotExist, err.Error(), &video.FavoriteActionResponse{}), nil
	}
	if err != nil {
		return vs.respRepo.FavoriteActionResponse(errcode.ErrDB, err.Error(), &video.FavoriteActionResponse{}), nil
	}
	// ② 去redis，判断是否对该视频有点赞记录,点赞记录没有过期时间
	exist, err := vs.rClient.IsFavRecordExist(ctx, favKey, req.VideoId)
	if err != nil {
		zap.L().Error("vs.rClient.IsFavRecordExist err:", zap.Error(err))
		return vs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &video.FavoriteActionResponse{}), nil
	}
	// 2.处理业务
	switch req.ActionType {
	case model.FAVORITE:
		if exist {
			return vs.respRepo.FavoriteActionResponse(myerr.AlreadyFavorite, model.MsgNil, &video.FavoriteActionResponse{}), nil
		}
		// ① 将 点赞关系 添加到 redis中
		err = vs.rClient.CreateFavorite(ctx, favKey, req.VideoId)
		if err != nil {
			zap.L().Error("vs.rClient.CreateFavorite err:", zap.Error(err))
			return vs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// 直接对redis中 user_info 和 video_info 的 count进行更改
		// ② 对视频 的 favorite_count+1
		err = vs.rClient.AddVideoFavoriteCount(ctx, videoKey)
		if err != nil {
			zap.L().Error("vs.rClient.AddVideoFavoriteCount err:", zap.Error(err))
			return vs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// ③ 对当前用户的 favorite_count+1，对视频的作者的 total_favorite+1
		updateUserFavCntResp, _ := client.UserClient.UpdateUserFavoriteCount(ctx, &user.UpdateUserFavoriteCountRequest{
			UserId:     req.UserId,
			AuthorId:   int64(videoInfo.UserId),
			ActionType: model.FAVORITE,
		})
		if updateUserFavCntResp == nil {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount 返回空指针")
			return vs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &video.FavoriteActionResponse{}), nil
		}
		if updateUserFavCntResp.StatusCode != 0 {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount err:", zap.Error(err))
			return vs.respRepo.FavoriteActionResponse(errcode.CreateErr(updateUserFavCntResp.StatusCode, model.MsgNil), updateUserFavCntResp.StatusMsg, &video.FavoriteActionResponse{}), nil
		}
	case model.CANCELFAVORITE:
		if !exist {
			return vs.respRepo.FavoriteActionResponse(myerr.IsNotFavorite, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// ① 将 点赞关系 从redis中删除
		err = vs.rClient.DelFavorite(ctx, favKey, req.VideoId)
		if err != nil {
			zap.L().Error("vs.rClient.DelFavorite err:", zap.Error(err))
			return vs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// 直接对redis中 user_info 和 video_info 的 count进行更改
		// ② 对视频 的 favorite_count-1
		err = vs.rClient.SubVideoFavoriteCount(ctx, videoKey)
		if err != nil {
			zap.L().Error("vs.rClient.SubVideoFavoriteCount err:", zap.Error(err))
			return vs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &video.FavoriteActionResponse{}), nil
		}
		// ③ 对当前用户的 favorite_count-1，对视频的作者的 total_favorite-1
		updateUserFavCntResp, _ := client.UserClient.UpdateUserFavoriteCount(ctx, &user.UpdateUserFavoriteCountRequest{
			UserId:     req.UserId,
			AuthorId:   int64(videoInfo.UserId),
			ActionType: model.CANCELFAVORITE,
		})
		if updateUserFavCntResp == nil {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount 返回空指针")
			return vs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &video.FavoriteActionResponse{}), nil
		}
		if updateUserFavCntResp.StatusCode != 0 {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount err:", zap.Error(err))
			return vs.respRepo.FavoriteActionResponse(errcode.CreateErr(updateUserFavCntResp.StatusCode, model.MsgNil), updateUserFavCntResp.StatusMsg, &video.FavoriteActionResponse{}), nil
		}
	}
	return vs.respRepo.FavoriteActionResponse(errcode.StatusOK, model.MsgNil, &video.FavoriteActionResponse{}), nil
}

func (vs *VideoServiceImpl) FavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	resp = new(video.FavoriteListResponse)
	// 1.校验业务逻辑
	// 2.处理业务
	// ① 根据 user_id 查询当前用户的所有信息
	userIndexResp, err := client.UserClient.UserIndex(ctx, &user.UserIndexRequest{UserId: req.UserId, Token: "", MyUserId: req.UserId})
	if userIndexResp == nil {
		zap.L().Error("client.UserClient.UserIndex 返回空指针")
		return vs.respRepo.FavoriteListResponse(errcode.ErrServer, model.MsgNil, &video.FavoriteListResponse{}), nil
	}
	if userIndexResp.StatusCode != 0 {
		zap.L().Error("client.UserClient.UserIndex err:", zap.Error(err))
		return vs.respRepo.FavoriteListResponse(errcode.CreateErr(userIndexResp.StatusCode, model.MsgNil), userIndexResp.StatusMsg, &video.FavoriteListResponse{}), nil
	}
	// ② 根据 user_id 在 redis中查询 点赞的video_ids
	favoriteKey := auto.CreateFavKey(uint(req.UserId))
	videoIds, err := vs.rClient.GetFavoriteVideoIds(ctx, favoriteKey)
	if err != nil {
		zap.L().Error("vs.rClient.GetFavoriteVideoIds err:", zap.Error(err))
		return vs.respRepo.FavoriteListResponse(errcode.ErrRedis, err.Error(), &video.FavoriteListResponse{}), nil
	}
	if videoIds == nil {
		return vs.respRepo.FavoriteListResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	// ③ 根据 video_id 在 redis中查询 video_info
	for _, videoId := range videoIds {
		videoInfo, err := vs.GetVideoInfo(ctx, videoId)
		if err != nil {
			return vs.respRepo.FavoriteListResponse(errcode.ErrDB, err.Error(), &video.FavoriteListResponse{}), nil
		}
		isFavorite, err := vs.IsFavoriteVideo(ctx, req.UserId, videoId)
		if err != nil {
			return vs.respRepo.FavoriteListResponse(errcode.ErrDB, err.Error(), &video.FavoriteListResponse{}), nil
		}
		resp.VideoList = append(resp.VideoList, &video.Video{
			Id:            int64(videoInfo.ID),
			Author:        userIndexResp.User,
			PlayUrl:       videoInfo.PlayURL,
			CoverUrl:      videoInfo.CoverURL,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoInfo.Title,
		})
	}
	// 4.返回数据
	return vs.respRepo.FavoriteListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) CommentAction(ctx context.Context, req *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	resp = new(video.CommentActionResponse)
	videoKey := auto.CreateVideoKey(uint(req.VideoId))
	// 1.校验业务逻辑
	// ① 判断该视频是否存在,去redis查询video_info，并获取作者user_id，没有查到就去查数据库，数据库中没有返回错误，查到开启线程池添加缓存
	//videoInfo, err := vs.GetVideoInfo(ctx, req.VideoId)
	//if err != nil {
	//	return vs.respRepo.CommentActionResponse(errcode.ErrDB, err.Error(), &video.CommentActionResponse{}), nil
	//}
	// 2.处理业务
	switch req.ActionType {
	case model.COMMENT:
		// ① 将 评论信息添加到数据库中
		// 构建评论实例
		commentId, _ := global.SnowFlake.NextId()
		timeNow := time.Now()
		commentInfo := &auto.Comment{
			BaseModel: auto.BaseModel{ID: uint(commentId), CreatedAt: timeNow, UpdatedAt: timeNow},
			VideoId:   uint(req.VideoId),
			UserId:    uint(req.UserId),
			Content:   *req.CommentText,
		}
		err = vs.videoRepo.CreateComment(ctx, commentInfo)
		if err != nil {
			zap.L().Error("vs.videoRepo.CreateComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrDB, err.Error(), &video.CommentActionResponse{}), nil
		}
		// 直接对redis中video_info 的 comment_count进行更改
		// ② 对视频 的 comment_count+1
		err = vs.rClient.AddVideoCommentCount(ctx, videoKey)
		if err != nil {
			zap.L().Error("vs.rClient.AddVideoCommentCount err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &video.CommentActionResponse{}), nil
		}
		// ③ 将comment加入 redis，用zset
		commentKey := auto.CreateCommentKey(req.VideoId)
		content := auto.CreateCommentValue(commentId, req.UserId, *req.CommentText)
		if err = vs.rClient.ZSAddCommentInfo(ctx, commentKey, float64(timeNow.Unix()), content); err != nil {
			zap.L().Error("vs.rClient.ZSAddCommentInfo err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &video.CommentActionResponse{}), nil
		}
		// ④ 获取用户信息
		userIndexResp, _ := client.UserClient.UserIndex(ctx, &user.UserIndexRequest{
			UserId:   req.UserId,
			Token:    "",
			MyUserId: req.UserId,
		})
		if userIndexResp == nil {
			zap.L().Error("client.UserClient.UserIndex 返回空指针")
			return vs.respRepo.CommentActionResponse(errcode.ErrServer, model.MsgNil, &video.CommentActionResponse{}), nil
		}
		if userIndexResp.StatusCode != 0 {
			zap.L().Error("client.UserClient.UserIndex err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.CreateErr(userIndexResp.StatusCode, model.MsgNil), userIndexResp.StatusMsg, &video.CommentActionResponse{}), nil
		}
		// 模型转换，格式化时间
		createDate := timeNow.Format("01-02")
		resp = &video.CommentActionResponse{
			Comment: &video.Comment{
				Id:         commentId,
				User:       userIndexResp.User,
				Content:    *req.CommentText,
				CreateDate: createDate,
			},
		}

	case model.CANCELCOMMENT:
		// ① 从数据库中删除评论数据
		commentInfo := &auto.Comment{
			BaseModel: auto.BaseModel{ID: uint(*req.CommentId)},
		}
		// 删除缓存
		commentKey := auto.CreateCommentKey(req.VideoId)
		if err = vs.rClient.DelComment(ctx, commentKey); err != nil {
			zap.L().Error("vs.rClient.DelComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &video.CommentActionResponse{}), nil
		}
		err = vs.videoRepo.DeleteComment(ctx, commentInfo)
		if err != nil {
			zap.L().Error("vs.videoRepo.DeleteComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrDB, err.Error(), &video.CommentActionResponse{}), nil
		}
		// 延迟 删除缓存
		time.Sleep(time.Millisecond * 50)
		if err = vs.rClient.DelComment(ctx, commentKey); err != nil {
			zap.L().Error("vs.rClient.DelComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &video.CommentActionResponse{}), nil
		}
		// 直接对redis中 video_info 的 comment_count进行更改
		// ② 对视频 的 comment_count-1
		err = vs.rClient.SubVideoCommentCount(ctx, videoKey)
		if err != nil {
			zap.L().Error("vs.rClient.SubVideoCommentCount err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &video.CommentActionResponse{}), nil
		}
		resp = &video.CommentActionResponse{}
	}
	return vs.respRepo.CommentActionResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) CommentList(ctx context.Context, req *video.CommentListRequest) (resp *video.CommentListResponse, err error) {
	resp = new(video.CommentListResponse)
	// 1.校验业务逻辑
	// 2.处理业务，读评论多，写评论相对没那么多
	// ① 去 redis 直接查评论信息
	commentKey := auto.CreateCommentKey(req.VideoId)
	userIds, commentList, err := vs.rClient.ZGetCommentList(ctx, commentKey)
	if err != nil {
		zap.L().Error("client.UserClient.ZGetCommentList err:", zap.Error(err))
		return vs.respRepo.CommentListResponse(errcode.ErrServer, model.MsgNil, &video.CommentListResponse{}), nil
	}
	// 无缓存,去数据库
	userIds, err = vs.videoRepo.GetCommentAuthorIds(ctx, req.VideoId)
	if err != nil {
		zap.L().Error("client.UserClient.GetCommentAuthorIds err:", zap.Error(err))
		return vs.respRepo.CommentListResponse(errcode.ErrDB, model.MsgNil, &video.CommentListResponse{}), nil
	}
	if userIds == nil {
		return vs.respRepo.CommentListResponse(errcode.StatusOK, model.MsgNil, &video.CommentListResponse{}), nil
	}
	commentList, err = vs.videoRepo.GetCommentList(ctx, req.VideoId)
	if err != nil {
		zap.L().Error("client.UserClient.GetCommentList err:", zap.Error(err))
		return vs.respRepo.CommentListResponse(errcode.ErrDB, model.MsgNil, &video.CommentListResponse{}), nil
	}
	// 添加缓存，消息队列
	for i, v := range commentList {
		// ③ 将comment加入 redis，用zset
		commentKey = auto.CreateCommentKey(req.VideoId)
		content := auto.CreateCommentValue(int64(v.ID), userIds[i], v.Content)
		if err = vs.rClient.ZSAddCommentInfo(ctx, commentKey, float64(v.CreatedAt.Unix()), content); err != nil {
			zap.L().Error("vs.rClient.ZSAddCommentInfo err:", zap.Error(err))
			return vs.respRepo.CommentListResponse(errcode.ErrRedis, err.Error(), &video.CommentListResponse{}), nil
		}
	}
	// ② 查出所有 user_id的详细信息
	userListResp, _ := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{
		UserId:   req.UserId,
		TargetId: userIds,
	})
	if userListResp == nil {
		zap.L().Error("client.UserClient.GetUserList 返回空指针")
		return vs.respRepo.CommentListResponse(errcode.ErrServer, model.MsgNil, &video.CommentListResponse{}), nil
	}
	if userListResp.StatusCode != 0 {
		zap.L().Error("client.UserClient.GetUserList err:", zap.Error(err))
		return vs.respRepo.CommentListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &video.CommentListResponse{}), nil
	}
	// 3.模型转换并返回
	for i, comment := range commentList {
		resp.CommentList = append(resp.CommentList, &video.Comment{
			Id:         int64(comment.ID),
			User:       userListResp.UserList[i],
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		})
	}
	// 4.返回数据
	return vs.respRepo.CommentListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) IsFavoriteVideo(ctx context.Context, userId, videoId int64) (bool, error) {
	// 查缓存
	favKey := auto.CreateFavKey(uint(userId))
	exist, err := vs.rClient.IsFavRecordExist(ctx, favKey, videoId)
	if err != nil {
		zap.L().Error("vs.rClient.IsFavRecordExist err:", zap.Error(err))
		return false, err
	}
	return exist, nil
}

func (vs *VideoServiceImpl) GetVideoInfo(ctx context.Context, videoId int64) (*auto.Video, error) {
	videoKey := auto.CreateVideoKey(uint(videoId))
	videoInfo, err := vs.rClient.HGetVideoInfo(ctx, videoKey)
	if err != nil {
		zap.L().Error("vs.rClient.HGetVideoInfo err:", zap.Error(err))
		return nil, err
	}
	if videoInfo == nil {
		// 没有缓存，查数据库
		videoInfo, err = vs.videoRepo.GetVideoInfoByVideoId(ctx, videoId)
		if err != nil {
			zap.L().Error("vs.rClient.GetVideoInfoByVideoId err:", zap.Error(err))
			return nil, err
		}
		// 添加缓存
		if err = vs.rClient.HSetVideoInfo(ctx, videoKey, auto.CreateMapVideoInfo(videoInfo)); err != nil {
			zap.L().Error("vs.rClient.HSetVideoInfo err:", zap.Error(err))
			return nil, err
		}
	}
	return videoInfo, nil
}
