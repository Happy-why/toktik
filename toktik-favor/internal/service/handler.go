package service

import (
	"context"
	"go.uber.org/zap"
	"toktik-common/errcode"
	"toktik-favor/internal/cache"
	"toktik-favor/internal/dao/mysql"
	"toktik-favor/internal/model"
	"toktik-favor/internal/model/auto"
	"toktik-favor/internal/repo"
	"toktik-favor/pkg/myerr"
	"toktik-favor/pkg/rpc/client"
	"toktik-rpc/kitex_gen/favor"
	"toktik-rpc/kitex_gen/user"
	"toktik-rpc/kitex_gen/video"
)

// FavorServiceImpl implements the last service interface defined in the IDL.
type FavorServiceImpl struct {
	favorRepo   repo.FavorRepo
	rCache      repo.RCacheRepo
	respRepo    repo.HandlerResp
	transaction mysql.Transaction
}

func NewFavorService() *FavorServiceImpl {
	return &FavorServiceImpl{
		favorRepo:   mysql.NewFavorDao(),
		rCache:      cache.NewFavorCache(),
		respRepo:    repo.NewHandlerResps(),
		transaction: mysql.NewTransaction(),
	}
}

func (fs *FavorServiceImpl) FavoriteAction(ctx context.Context, req *favor.FavoriteActionRequest) (resp *favor.FavoriteActionResponse, err error) {
	favKey := auto.CreateFavKey(uint(req.UserId))
	// 1.校验业务逻辑
	// ① 判断该视频是否存在,去redis查询video_info，并获取作者user_id，没有查到就去查数据库，数据库中没有返回错误，查到开启线程池添加缓存
	// 去 video 服务，查找video_info
	getVideoInfoResp, _ := client.VideoClient.GetVideoInfo(ctx, &video.GetVideoInfoRequest{VideoId: req.VideoId})
	if getVideoInfoResp == nil {
		zap.L().Error("client.VideoClient.GetVideoInfo 返回空指针")
		return fs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &favor.FavoriteActionResponse{}), nil
	}
	if getVideoInfoResp.StatusCode != 0 {
		zap.L().Error("client.VideoClient.GetVideoInfo err:", zap.Error(err))
		return fs.respRepo.FavoriteActionResponse(errcode.CreateErr(getVideoInfoResp.StatusCode, model.MsgNil), getVideoInfoResp.StatusMsg, &favor.FavoriteActionResponse{}), nil
	}
	// ② 去redis，判断是否对该视频有点赞记录,点赞记录没有过期时间
	exist, err := fs.rCache.IsFavRecordExist(ctx, favKey, req.VideoId)
	if err != nil {
		zap.L().Error("fs.rCache.IsFavRecordExist err:", zap.Error(err))
		return fs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &favor.FavoriteActionResponse{}), nil
	}
	// 2.处理业务
	switch req.ActionType {
	case model.FAVORITE:
		if exist {
			return fs.respRepo.FavoriteActionResponse(myerr.AlreadyFavorite, model.MsgNil, &favor.FavoriteActionResponse{}), nil
		}
		// ① 将 点赞关系 添加到 redis中
		err = fs.rCache.CreateFavorite(ctx, favKey, req.VideoId)
		if err != nil {
			zap.L().Error("fs.rCache.CreateFavorite err:", zap.Error(err))
			return fs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &favor.FavoriteActionResponse{}), nil
		}
		// 直接对redis中 user_info 和 video_info 的 count进行更改
		// ② 对视频 的 favorite_count+1
		addFavCntResp, _ := client.VideoClient.AddVideoFavoriteCount(ctx, &video.AddVideoFavoriteCountRequest{VideoId: req.VideoId})
		if addFavCntResp == nil {
			zap.L().Error("client.VideoClient.AddVideoFavoriteCount 返回空指针")
			return fs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &favor.FavoriteActionResponse{}), nil
		}
		if addFavCntResp.StatusCode != 0 {
			zap.L().Error("client.VideoClient.AddVideoFavoriteCount err:", zap.Error(err))
			return fs.respRepo.FavoriteActionResponse(errcode.CreateErr(addFavCntResp.StatusCode, model.MsgNil), addFavCntResp.StatusMsg, &favor.FavoriteActionResponse{}), nil
		}
		// ③ 对当前用户的 favorite_count+1，对视频的作者的 total_favorite+1
		updateUserFavCntResp, _ := client.UserClient.UpdateUserFavoriteCount(ctx, &user.UpdateUserFavoriteCountRequest{
			UserId:     req.UserId,
			AuthorId:   getVideoInfoResp.VideoInfo.Author.Id,
			ActionType: model.FAVORITE,
		})
		if updateUserFavCntResp == nil {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount 返回空指针")
			return fs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &favor.FavoriteActionResponse{}), nil
		}
		if updateUserFavCntResp.StatusCode != 0 {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount err:", zap.Error(err))
			return fs.respRepo.FavoriteActionResponse(errcode.CreateErr(updateUserFavCntResp.StatusCode, model.MsgNil), updateUserFavCntResp.StatusMsg, &favor.FavoriteActionResponse{}), nil
		}
	case model.CANCELFAVORITE:
		if !exist {
			return fs.respRepo.FavoriteActionResponse(myerr.IsNotFavorite, err.Error(), &favor.FavoriteActionResponse{}), nil
		}
		// ① 将 点赞关系 从redis中删除
		err = fs.rCache.DelFavorite(ctx, favKey, req.VideoId)
		if err != nil {
			zap.L().Error("fs.rCache.DelFavorite err:", zap.Error(err))
			return fs.respRepo.FavoriteActionResponse(errcode.ErrRedis, err.Error(), &favor.FavoriteActionResponse{}), nil
		}
		// 直接对redis中 user_info 和 video_info 的 count进行更改
		// ② 对视频 的 favorite_count-1
		subFavCntResp, _ := client.VideoClient.SubVideoFavoriteCount(ctx, &video.SubVideoFavoriteCountRequest{VideoId: req.VideoId})
		if subFavCntResp == nil {
			zap.L().Error("client.VideoClient.SubVideoFavoriteCount 返回空指针")
			return fs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &favor.FavoriteActionResponse{}), nil
		}
		if subFavCntResp.StatusCode != 0 {
			zap.L().Error("client.VideoClient.SubVideoFavoriteCount err:", zap.Error(err))
			return fs.respRepo.FavoriteActionResponse(errcode.CreateErr(subFavCntResp.StatusCode, model.MsgNil), subFavCntResp.StatusMsg, &favor.FavoriteActionResponse{}), nil
		}
		// ③ 对当前用户的 favorite_count-1，对视频的作者的 total_favorite-1
		updateUserFavCntResp, _ := client.UserClient.UpdateUserFavoriteCount(ctx, &user.UpdateUserFavoriteCountRequest{
			UserId:     req.UserId,
			AuthorId:   getVideoInfoResp.VideoInfo.Author.Id,
			ActionType: model.CANCELFAVORITE,
		})
		if updateUserFavCntResp == nil {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount 返回空指针")
			return fs.respRepo.FavoriteActionResponse(errcode.ErrServer, model.MsgNil, &favor.FavoriteActionResponse{}), nil
		}
		if updateUserFavCntResp.StatusCode != 0 {
			zap.L().Error("client.UserClient.UpdateUserFavoriteCount err:", zap.Error(err))
			return fs.respRepo.FavoriteActionResponse(errcode.CreateErr(updateUserFavCntResp.StatusCode, model.MsgNil), updateUserFavCntResp.StatusMsg, &favor.FavoriteActionResponse{}), nil
		}
	default:
		return fs.respRepo.FavoriteActionResponse(myerr.OperationErr, model.MsgNil, &favor.FavoriteActionResponse{}), nil
	}
	return fs.respRepo.FavoriteActionResponse(errcode.StatusOK, model.MsgNil, &favor.FavoriteActionResponse{}), nil
}

func (fs *FavorServiceImpl) FavoriteList(ctx context.Context, req *favor.FavoriteListRequest) (resp *favor.FavoriteListResponse, err error) {
	resp = new(favor.FavoriteListResponse)
	// 1.校验业务逻辑
	// 2.处理业务
	// ② 根据 user_id 在 redis中查询 点赞的video_ids
	favoriteKey := auto.CreateFavKey(uint(req.UserId))
	videoIds, err := fs.rCache.GetFavoriteVideoIds(ctx, favoriteKey)
	if err != nil {
		zap.L().Error("fs.rCache.GetFavoriteVideoIds err:", zap.Error(err))
		return fs.respRepo.FavoriteListResponse(errcode.ErrRedis, err.Error(), &favor.FavoriteListResponse{}), nil
	}
	if videoIds == nil {
		//
		return fs.respRepo.FavoriteListResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	// ③ 根据 video_ids 查询到 每个 video_info
	getManyVideoInfosResp, _ := client.VideoClient.GetManyVideoInfos(ctx, &video.GetManyVideoInfosRequest{
		UserId:   req.UserId,
		VideoIds: videoIds,
	})
	if getManyVideoInfosResp == nil {
		zap.L().Error("client.VideoClient.GetManyVideoInfos 返回空指针")
		return fs.respRepo.FavoriteListResponse(errcode.ErrServer, model.MsgNil, &favor.FavoriteListResponse{}), nil
	}
	if getManyVideoInfosResp.StatusCode != 0 {
		zap.L().Error("client.VideoClient.GetManyVideoInfos err:", zap.Error(err))
		return fs.respRepo.FavoriteListResponse(errcode.CreateErr(getManyVideoInfosResp.StatusCode, model.MsgNil), getManyVideoInfosResp.StatusMsg, &favor.FavoriteListResponse{}), nil
	}
	resp.VideoList = make([]*video.Video, len(videoIds))
	for i, videoInfo := range getManyVideoInfosResp.VideoInfos {
		resp.VideoList[i] = videoInfo
	}
	// 4.返回数据
	return fs.respRepo.FavoriteListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (fs *FavorServiceImpl) IsFavoriteVideos(ctx context.Context, req *favor.IsFavoriteVideosRequest) (resp *favor.IsFavoriteVideosResponse, err error) {
	resp = new(favor.IsFavoriteVideosResponse)
	resp.ManyIsFavorite = make([]bool, len(req.VideoIds))
	// 查缓存
	favKey := auto.CreateFavKey(uint(req.UserId))
	for i, videoId := range req.VideoIds {
		exist, err := fs.rCache.IsFavRecordExist(ctx, favKey, videoId)
		if err != nil {
			zap.L().Error("fs.rCache.IsFavRecordExist err:", zap.Error(err))
			return fs.respRepo.IsFavoriteVideoResponse(errcode.ErrRedis, err.Error(), &favor.IsFavoriteVideosResponse{}), nil
		}
		resp.ManyIsFavorite[i] = exist
	}
	return fs.respRepo.IsFavoriteVideoResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
