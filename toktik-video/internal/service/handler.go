package service

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"time"
	"toktik-common/errcode"
	"toktik-common/oss"
	"toktik-common/utils"
	"toktik-rpc/kitex_gen/favor"
	"toktik-rpc/kitex_gen/user"
	"toktik-rpc/kitex_gen/video"
	"toktik-video/internal/cache"
	"toktik-video/internal/dao/mysql"
	"toktik-video/internal/global"
	"toktik-video/internal/model"
	"toktik-video/internal/model/auto"
	"toktik-video/internal/repo"
	"toktik-video/pkg/myerr"
	"toktik-video/pkg/rpc/client"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	videoRepo   repo.VideoRepo
	rCache      repo.RCacheRepo
	respRepo    repo.HandlerResp
	transaction mysql.Transaction
	ossRepo     oss.OSS
}

func NewVideoService() *VideoServiceImpl {
	return &VideoServiceImpl{
		videoRepo:   mysql.NewVideoDao(),
		rCache:      cache.NewVideoCache(),
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
	nextTime, videoIds, err := vs.rCache.ZGetVideoIds(ctx, feedKey, *req.LatestTime)
	if videoIds == nil {
		return vs.respRepo.VideoFeedResponse(myerr.CanNotSearchVideo, model.MsgNil, &video.VideoFeedResponse{}), nil
	}
	if err != nil {
		zap.L().Error("vs.rCache.ZGetVideoIds err:", zap.Error(err))
		return vs.respRepo.VideoFeedResponse(errcode.ErrRedis, err.Error(), &video.VideoFeedResponse{}), nil
	}
	videoInfos := make([]*auto.Video, len(videoIds))
	userIds := make([]int64, len(videoIds))
	for i, videoId := range videoIds {
		// ② 根据 video_id 从缓存中查询 video_info
		videoInfo, err := vs.GetVideoInfoSelf(ctx, videoId)
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
	manyIsFavorite := make([]bool, len(videoIds))
	if isLogin {
		// ⑤ 查询出用户 对每个视频是否点赞
		IsFavoriteResp, _ := client.FavorClient.IsFavoriteVideos(ctx, &favor.IsFavoriteVideosRequest{
			UserId:   req.UserId,
			VideoIds: videoIds,
		})
		if IsFavoriteResp == nil {
			zap.L().Error("client.FavorClient.IsFavoriteVideos 返回空指针")
			return vs.respRepo.VideoFeedResponse(errcode.ErrServer, model.MsgNil, &video.VideoFeedResponse{}), nil
		}
		if IsFavoriteResp.StatusCode != model.RpcSuccess {
			zap.L().Error("client.FavorClient.IsFavoriteVideos err:", zap.Error(err))
			return vs.respRepo.VideoFeedResponse(errcode.CreateErr(IsFavoriteResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &video.VideoFeedResponse{}), nil
		}
		manyIsFavorite = IsFavoriteResp.ManyIsFavorite
	}

	for i, v := range videoInfos {
		resp.VideoList = append(resp.VideoList, &video.Video{
			Id:            int64(v.ID),
			Author:        userListMap[v.UserId],
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    manyIsFavorite[i],
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
	if err = vs.rCache.HSetVideoInfo(ctx, videoKey, videoMap); err != nil {
		zap.L().Error("vs.rCache.HSetVideoInfo err:", zap.Error(err))
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
	if err = vs.rCache.PublishVideo(ctx, publishKey, float64(timeNow.Unix()), videoId); err != nil {
		zap.L().Error("vs.rCache.PublishVideo err:", zap.Error(err))
		return vs.respRepo.VideoPublishResponse(errcode.ErrRedis, err.Error(), &video.VideoPublishResponse{}), nil
	}
	// ⑧ 删除用户的视频缓存，user_video::video_ids
	userVideoKey := auto.CreateUserVideoKey(uint(req.UserId))
	if err = vs.rCache.SDelUserVideo(ctx, userVideoKey); err != nil {
		zap.L().Error("vs.rCache.SDelUserVideo err:", zap.Error(err))
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
	videoIds, err := vs.rCache.SGetUserVideoIds(ctx, userVideoKey)
	if err != nil {
		zap.L().Error("vs.rCache.SGetUserVideoIds err:", zap.Error(err))
		return vs.respRepo.PublishListResponse(errcode.ErrRedis, err.Error(), &video.PublishListResponse{}), nil
	}
	if videoIds == nil {
		// 缓存没查到，去数据库查
		videoIds, err = vs.videoRepo.GetVideoIdsByUserId(ctx, req.UserId)
		if err != nil {
			zap.L().Error("vs.videoRepo.GetVideoIdsByUserId err:", zap.Error(err))
			return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
		}
		if len(videoIds) == 0 {
			return vs.respRepo.PublishListResponse(errcode.StatusOK, model.MsgNil, resp), nil
		}
		// 添加缓存
		if err = vs.rCache.SAddUserVideoIds(ctx, userVideoKey, videoIds); err != nil {
			zap.L().Error("vs.rCache.SAddUserVideoIds err:", zap.Error(err))
			return vs.respRepo.PublishListResponse(errcode.ErrRedis, err.Error(), &video.PublishListResponse{}), nil
		}
	}
	// ③ 判断 用户对每个视频是否点赞
	IsFavoriteResp, _ := client.FavorClient.IsFavoriteVideos(ctx, &favor.IsFavoriteVideosRequest{
		UserId:   req.UserId,
		VideoIds: videoIds,
	})
	if IsFavoriteResp == nil {
		zap.L().Error("client.FavorClient.IsFavoriteVideos 返回空指针")
		return vs.respRepo.PublishListResponse(errcode.ErrServer, model.MsgNil, &video.PublishListResponse{}), nil
	}
	if IsFavoriteResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.FavorClient.IsFavoriteVideos err:", zap.Error(err))
		return vs.respRepo.PublishListResponse(errcode.CreateErr(IsFavoriteResp.StatusCode, model.MsgNil), IsFavoriteResp.StatusMsg, &video.PublishListResponse{}), nil
	}
	resp.VideoList = make([]*video.Video, len(videoIds))
	// ③ 查找 videoInfos ，
	for i, videoId := range videoIds {
		// 根据 video_id 从缓存中查询 video_info
		videoKey := auto.CreateVideoKey(uint(videoId))
		videoInfo, err := vs.rCache.HGetVideoInfo(ctx, videoKey)
		if err != nil {
			zap.L().Error("vs.rCache.HGetVideoInfo err:", zap.Error(err))
			return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
		}
		if videoInfo == nil {
			// 没有缓存，查数据库
			videoInfo, err = vs.videoRepo.GetVideoInfoByVideoId(ctx, videoId)
			if err != nil {
				zap.L().Error("vs.rCache.GetVideoInfoByVideoId err:", zap.Error(err))
				return vs.respRepo.PublishListResponse(errcode.ErrDB, err.Error(), &video.PublishListResponse{}), nil
			}
			if videoInfo == nil {
				return vs.respRepo.PublishListResponse(myerr.VideoNotExist, model.MsgNil, &video.PublishListResponse{}), nil
			}
			// 添加缓存
			if err = vs.rCache.HSetVideoInfo(ctx, videoKey, auto.CreateMapVideoInfo(videoInfo)); err != nil {
				zap.L().Error("vs.rCache.HSetVideoInfo err:", zap.Error(err))
				return vs.respRepo.PublishListResponse(errcode.ErrRedis, err.Error(), &video.PublishListResponse{}), nil
			}
		}
		resp.VideoList[i] = &video.Video{
			Id:            videoId,
			Author:        userIndexResp.User,
			PlayUrl:       videoInfo.PlayURL,
			CoverUrl:      videoInfo.CoverURL,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    IsFavoriteResp.ManyIsFavorite[i],
			Title:         videoInfo.Title,
		}
	}
	return vs.respRepo.PublishListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) GetVideoInfo(ctx context.Context, req *video.GetVideoInfoRequest) (resp *video.GetVideoInfoResponse, err error) {
	resp = new(video.GetVideoInfoResponse)
	videoInfo, err := vs.GetVideoInfoSelf(ctx, req.VideoId)
	if err != nil {
		zap.L().Error("vs.GetVideoInfoSelf err:", zap.Error(err))
		return vs.respRepo.GetVideoInfoResponse(errcode.ErrDB, err.Error(), &video.GetVideoInfoResponse{}), nil
	}
	// 通过 user_id 去查 user_info
	userIndexResp, _ := client.UserClient.UserIndex(ctx, &user.UserIndexRequest{
		UserId:   int64(videoInfo.UserId),
		Token:    "",
		MyUserId: req.UserId,
	})
	if userIndexResp == nil {
		zap.L().Error("client.UserClient.UserIndex 返回空指针")
		return vs.respRepo.GetVideoInfoResponse(errcode.ErrServer, model.MsgNil, &video.GetVideoInfoResponse{}), nil
	}
	if userIndexResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.UserClient.UserIndex err:", zap.Error(err))
		return vs.respRepo.GetVideoInfoResponse(errcode.CreateErr(userIndexResp.StatusCode, model.MsgNil), userIndexResp.StatusMsg, &video.GetVideoInfoResponse{}), nil
	}
	// 判断是否对该视频点过赞
	isFavorVideoResp, _ := client.FavorClient.IsFavoriteVideos(ctx, &favor.IsFavoriteVideosRequest{
		UserId:   req.UserId,
		VideoIds: []int64{req.VideoId},
	})
	if isFavorVideoResp == nil {
		zap.L().Error("client.FavorClient.IsFavoriteVideos 返回空指针")
		return vs.respRepo.GetVideoInfoResponse(errcode.ErrServer, model.MsgNil, &video.GetVideoInfoResponse{}), nil
	}
	if isFavorVideoResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.FavorClient.IsFavoriteVideos err:", zap.Error(err))
		return vs.respRepo.GetVideoInfoResponse(errcode.CreateErr(isFavorVideoResp.StatusCode, model.MsgNil), isFavorVideoResp.StatusMsg, &video.GetVideoInfoResponse{}), nil
	}
	resp.VideoInfo = &video.Video{
		Id:            int64(videoInfo.ID),
		Author:        userIndexResp.User,
		PlayUrl:       videoInfo.PlayURL,
		CoverUrl:      videoInfo.CoverURL,
		FavoriteCount: videoInfo.FavoriteCount,
		CommentCount:  videoInfo.CommentCount,
		IsFavorite:    isFavorVideoResp.ManyIsFavorite[0],
		Title:         videoInfo.Title,
	}
	return vs.respRepo.GetVideoInfoResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) GetManyVideoInfos(ctx context.Context, req *video.GetManyVideoInfosRequest) (resp *video.GetManyVideoInfosResponse, err error) {
	resp = new(video.GetManyVideoInfosResponse)
	// 判断是否对每个 video 点过赞
	isFavorVideosResp, _ := client.FavorClient.IsFavoriteVideos(ctx, &favor.IsFavoriteVideosRequest{
		UserId:   req.UserId,
		VideoIds: req.VideoIds,
	})
	if isFavorVideosResp == nil {
		zap.L().Error("client.FavorClient.IsFavoriteVideos 返回空指针")
		return vs.respRepo.GetManyVideoInfosResponse(errcode.ErrServer, model.MsgNil, &video.GetManyVideoInfosResponse{}), nil
	}
	if isFavorVideosResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.FavorClient.IsFavoriteVideos err:", zap.Error(err))
		return vs.respRepo.GetManyVideoInfosResponse(errcode.CreateErr(isFavorVideosResp.StatusCode, model.MsgNil), isFavorVideosResp.StatusMsg, &video.GetManyVideoInfosResponse{}), nil
	}
	// 先获取所有video_info
	resp.VideoInfos = make([]*video.Video, len(req.VideoIds))
	authorIds := make([]int64, len(req.VideoIds))
	for i, videoId := range req.VideoIds {
		videoInfo, err := vs.GetVideoInfoSelf(ctx, videoId)
		if err != nil {
			zap.L().Error("vs.GetVideoInfoSelf err:", zap.Error(err))
			return vs.respRepo.GetManyVideoInfosResponse(errcode.ErrDB, err.Error(), &video.GetManyVideoInfosResponse{}), nil
		}
		resp.VideoInfos[i] = &video.Video{
			Id: int64(videoInfo.ID),
			//Author:        getUserListResp.UserList[i],
			PlayUrl:       videoInfo.PlayURL,
			CoverUrl:      videoInfo.CoverURL,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    isFavorVideosResp.ManyIsFavorite[i],
			Title:         videoInfo.Title,
		}
		// 根据 video_id 查询到每个 video 的作者 ，方便下面添加作者信息
		authorIds[i] = int64(videoInfo.UserId)
	}

	getUserListResp, _ := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{
		UserId:   req.UserId,
		TargetId: authorIds,
	})
	if getUserListResp == nil {
		zap.L().Error("client.UserClient.GetUserListRequest 返回空指针")
		return vs.respRepo.GetManyVideoInfosResponse(errcode.ErrServer, model.MsgNil, &video.GetManyVideoInfosResponse{}), nil
	}
	if getUserListResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.UserClient.GetUserListRequest err:", zap.Error(err))
		return vs.respRepo.GetManyVideoInfosResponse(errcode.CreateErr(getUserListResp.StatusCode, model.MsgNil), getUserListResp.StatusMsg, &video.GetManyVideoInfosResponse{}), nil
	}
	for i, authorInfo := range getUserListResp.UserList {
		resp.VideoInfos[i].Author = authorInfo
	}
	return vs.respRepo.GetManyVideoInfosResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) AddVideoFavoriteCount(ctx context.Context, req *video.AddVideoFavoriteCountRequest) (resp *video.AddVideoFavoriteCountResponse, err error) {
	// 添加缓存中 video_info 中的 favorite_count
	videoKey := auto.CreateVideoKey(uint(req.VideoId))
	if err = vs.rCache.AddVideoFavoriteCount(ctx, videoKey); err != nil {
		zap.L().Error("vs.rCache.AddVideoFavoriteCount err:", zap.Error(err))
		return vs.respRepo.AddVideoFavoriteCountResponse(errcode.ErrRedis, err.Error(), &video.AddVideoFavoriteCountResponse{}), nil
	}
	return vs.respRepo.AddVideoFavoriteCountResponse(errcode.StatusOK, model.MsgNil, &video.AddVideoFavoriteCountResponse{}), nil
}

func (vs *VideoServiceImpl) SubVideoFavoriteCount(ctx context.Context, req *video.SubVideoFavoriteCountRequest) (resp *video.SubVideoFavoriteCountResponse, err error) {
	// 减少缓存中 video_info 中的 favorite_count
	videoKey := auto.CreateVideoKey(uint(req.VideoId))
	if err = vs.rCache.SubVideoFavoriteCount(ctx, videoKey); err != nil {
		zap.L().Error("vs.rCache.SubVideoFavoriteCount err:", zap.Error(err))
		return vs.respRepo.SubVideoFavoriteCountResponse(errcode.ErrRedis, err.Error(), &video.SubVideoFavoriteCountResponse{}), nil
	}
	return vs.respRepo.SubVideoFavoriteCountResponse(errcode.StatusOK, model.MsgNil, &video.SubVideoFavoriteCountResponse{}), nil
}

func (vs *VideoServiceImpl) AddVideoCommentCount(ctx context.Context, req *video.AddVideoCommentCountRequest) (resp *video.AddVideoCommentCountResponse, err error) {
	// 添加缓存中 video_info 中的 comment_count
	videoKey := auto.CreateVideoKey(uint(req.VideoId))
	if err = vs.rCache.AddVideoCommentCount(ctx, videoKey); err != nil {
		zap.L().Error("vs.rCache.AddVideoCommentCount err:", zap.Error(err))
		return vs.respRepo.AddVideoCommentCountResponse(errcode.ErrRedis, err.Error(), &video.AddVideoCommentCountResponse{}), nil
	}
	return vs.respRepo.AddVideoCommentCountResponse(errcode.StatusOK, model.MsgNil, &video.AddVideoCommentCountResponse{}), nil
}

func (vs *VideoServiceImpl) SubVideoCommentCount(ctx context.Context, req *video.SubVideoCommentCountRequest) (resp *video.SubVideoCommentCountResponse, err error) {
	// 添加缓存中 video_info 中的 comment_count
	videoKey := auto.CreateVideoKey(uint(req.VideoId))
	if err = vs.rCache.SubVideoCommentCount(ctx, videoKey); err != nil {
		zap.L().Error("vs.rCache.SubVideoCommentCount err:", zap.Error(err))
		return vs.respRepo.SubVideoCommentCountResponse(errcode.ErrRedis, err.Error(), &video.SubVideoCommentCountResponse{}), nil
	}
	return vs.respRepo.SubVideoCommentCountResponse(errcode.StatusOK, model.MsgNil, &video.SubVideoCommentCountResponse{}), nil
}

func (vs *VideoServiceImpl) GetVideoInfoSelf(ctx context.Context, videoId int64) (*auto.Video, error) {
	videoKey := auto.CreateVideoKey(uint(videoId))
	videoInfo, err := vs.rCache.HGetVideoInfo(ctx, videoKey)
	if err != nil {
		zap.L().Error("vs.rCache.HGetVideoInfo err:", zap.Error(err))
		return nil, err
	}
	if videoInfo == nil {
		// 没有缓存，查数据库
		videoInfo, err = vs.videoRepo.GetVideoInfoByVideoId(ctx, videoId)
		if err != nil {
			zap.L().Error("vs.rCache.GetVideoInfoByVideoId err:", zap.Error(err))
			return nil, err
		}
		if videoInfo == nil {
			return nil, err
		}
		// 添加缓存
		if err = vs.rCache.HSetVideoInfo(ctx, videoKey, auto.CreateMapVideoInfo(videoInfo)); err != nil {
			zap.L().Error("vs.rCache.HSetVideoInfo err:", zap.Error(err))
			return nil, err
		}
	}
	return videoInfo, nil
}
