package service

import (
	"context"
	"go.uber.org/zap"
	"time"
	"toktik-comment/internal/cache"
	"toktik-comment/internal/dao/mysql"
	"toktik-comment/internal/global"
	"toktik-comment/internal/model"
	"toktik-comment/internal/model/auto"
	"toktik-comment/internal/repo"
	"toktik-comment/pkg/myerr"
	"toktik-comment/pkg/rpc/client"
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/comment"
	"toktik-rpc/kitex_gen/user"
	"toktik-rpc/kitex_gen/video"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	videoRepo   repo.VideoRepo
	rCache      repo.RCacheRepo
	respRepo    repo.HandlerResp
	transaction mysql.Transaction
}

func NewVideoService() *VideoServiceImpl {
	return &VideoServiceImpl{
		videoRepo:   mysql.NewCommentDao(),
		rCache:      cache.NewCommentCache(),
		respRepo:    repo.NewHandlerResps(),
		transaction: mysql.NewTransaction(),
	}
}

func (vs *VideoServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	resp = new(comment.CommentActionResponse)
	// 1.校验业务逻辑
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
			return vs.respRepo.CommentActionResponse(errcode.ErrDB, err.Error(), &comment.CommentActionResponse{}), nil
		}
		// 直接对redis中video_info 的 comment_count进行更改
		// ② 对视频 的 comment_count+1
		addCommentCntResp, _ := client.VideoClient.AddVideoCommentCount(ctx, &video.AddVideoCommentCountRequest{VideoId: req.VideoId})
		if addCommentCntResp == nil {
			zap.L().Error("client.VideoClient.AddVideoCommentCount 返回空指针")
			return vs.respRepo.CommentActionResponse(errcode.ErrServer, model.MsgNil, &comment.CommentActionResponse{}), nil
		}
		if addCommentCntResp.StatusCode != 0 {
			zap.L().Error("client.VideoClient.AddVideoCommentCount err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.CreateErr(addCommentCntResp.StatusCode, model.MsgNil), addCommentCntResp.StatusMsg, &comment.CommentActionResponse{}), nil
		}
		// ③ 将comment加入 redis，用zset
		commentKey := auto.CreateCommentKey(req.VideoId)
		content := auto.CreateCommentValue(commentId, req.UserId, *req.CommentText)
		if err = vs.rCache.ZSAddCommentInfo(ctx, commentKey, float64(timeNow.Unix()), content); err != nil {
			zap.L().Error("vs.rCache.ZSAddCommentInfo err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &comment.CommentActionResponse{}), nil
		}
		// ④ 获取用户信息
		userIndexResp, _ := client.UserClient.UserIndex(ctx, &user.UserIndexRequest{
			UserId:   req.UserId,
			Token:    "",
			MyUserId: req.UserId,
		})
		if userIndexResp == nil {
			zap.L().Error("client.UserClient.UserIndex 返回空指针")
			return vs.respRepo.CommentActionResponse(errcode.ErrServer, model.MsgNil, &comment.CommentActionResponse{}), nil
		}
		if userIndexResp.StatusCode != 0 {
			zap.L().Error("client.UserClient.UserIndex err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.CreateErr(userIndexResp.StatusCode, model.MsgNil), userIndexResp.StatusMsg, &comment.CommentActionResponse{}), nil
		}
		// 模型转换，格式化时间
		createDate := timeNow.Format("01-02")
		resp = &comment.CommentActionResponse{
			Comment: &comment.Comment{
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
		if err = vs.rCache.DelComment(ctx, commentKey); err != nil {
			zap.L().Error("vs.rCache.DelComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &comment.CommentActionResponse{}), nil
		}
		err = vs.videoRepo.DeleteComment(ctx, commentInfo)
		if err != nil {
			zap.L().Error("vs.videoRepo.DeleteComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrDB, err.Error(), &comment.CommentActionResponse{}), nil
		}
		// 延迟 删除缓存
		time.Sleep(time.Millisecond * 50)
		if err = vs.rCache.DelComment(ctx, commentKey); err != nil {
			zap.L().Error("vs.rCache.DelComment err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.ErrRedis, err.Error(), &comment.CommentActionResponse{}), nil
		}
		// 直接对redis中 video_info 的 comment_count进行更改
		// ② 对视频 的 comment_count-1
		subCommentCntResp, _ := client.VideoClient.SubVideoCommentCount(ctx, &video.SubVideoCommentCountRequest{VideoId: req.VideoId})
		if subCommentCntResp == nil {
			zap.L().Error("client.VideoClient.SubVideoCommentCount 返回空指针")
			return vs.respRepo.CommentActionResponse(errcode.ErrServer, model.MsgNil, &comment.CommentActionResponse{}), nil
		}
		if subCommentCntResp.StatusCode != 0 {
			zap.L().Error("client.VideoClient.SubVideoCommentCount err:", zap.Error(err))
			return vs.respRepo.CommentActionResponse(errcode.CreateErr(subCommentCntResp.StatusCode, model.MsgNil), subCommentCntResp.StatusMsg, &comment.CommentActionResponse{}), nil
		}
		resp = &comment.CommentActionResponse{}

	default:
		return vs.respRepo.CommentActionResponse(myerr.OperationErr, model.MsgNil, &comment.CommentActionResponse{}), nil
	}
	return vs.respRepo.CommentActionResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (vs *VideoServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	resp = new(comment.CommentListResponse)
	// 1.校验业务逻辑
	// 2.处理业务，读评论多，写评论相对没那么多
	// ① 去 redis 直接查评论信息
	commentKey := auto.CreateCommentKey(req.VideoId)
	userIds, commentList, err := vs.rCache.ZGetCommentList(ctx, commentKey)
	if err != nil {
		zap.L().Error("client.UserClient.ZGetCommentList err:", zap.Error(err))
		return vs.respRepo.CommentListResponse(errcode.ErrServer, model.MsgNil, &comment.CommentListResponse{}), nil
	}
	if userIds == nil {
		zap.L().Info("评论 没有被缓存！！！！查数据库！！！！")
		// 无缓存,去数据库
		userIds, err = vs.videoRepo.GetCommentAuthorIds(ctx, req.VideoId)
		if err != nil {
			zap.L().Error("client.UserClient.GetCommentAuthorIds err:", zap.Error(err))
			return vs.respRepo.CommentListResponse(errcode.ErrDB, model.MsgNil, &comment.CommentListResponse{}), nil
		}
		if userIds == nil {
			return vs.respRepo.CommentListResponse(errcode.StatusOK, model.MsgNil, &comment.CommentListResponse{}), nil
		}
		commentList, err = vs.videoRepo.GetCommentList(ctx, req.VideoId)
		if err != nil {
			zap.L().Error("client.UserClient.GetCommentList err:", zap.Error(err))
			return vs.respRepo.CommentListResponse(errcode.ErrDB, model.MsgNil, &comment.CommentListResponse{}), nil
		}
		// 添加缓存，消息队列,可以批量存入
		commentKey = auto.CreateCommentKey(req.VideoId)
		for i, v := range commentList {
			// ③ 将comment加入 redis，用zset
			content := auto.CreateCommentValue(int64(v.ID), userIds[i], v.Content)
			if err = vs.rCache.ZSAddCommentInfo(ctx, commentKey, float64(v.CreatedAt.Unix()), content); err != nil {
				zap.L().Error("vs.rCache.ZSAddCommentInfo err:", zap.Error(err))
				return vs.respRepo.CommentListResponse(errcode.ErrRedis, err.Error(), &comment.CommentListResponse{}), nil
			}
		}
	}

	// ② 查出所有 user_id的详细信息
	userListResp, _ := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{
		UserId:   req.UserId,
		TargetId: userIds,
	})
	if userListResp == nil {
		zap.L().Error("client.UserClient.GetUserList 返回空指针")
		return vs.respRepo.CommentListResponse(errcode.ErrServer, model.MsgNil, &comment.CommentListResponse{}), nil
	}
	if userListResp.StatusCode != 0 {
		zap.L().Error("client.UserClient.GetUserList err:", zap.Error(err))
		return vs.respRepo.CommentListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &comment.CommentListResponse{}), nil
	}
	// 3.模型转换并返回
	resp.CommentList = make([]*comment.Comment, len(commentList))
	for i, v := range commentList {
		resp.CommentList[i] = &comment.Comment{
			Id:         int64(v.ID),
			User:       userListResp.UserList[i],
			Content:    v.Content,
			CreateDate: v.CreatedAt.Format("01-02"),
		}
	}
	// 4.返回数据
	return vs.respRepo.CommentListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
