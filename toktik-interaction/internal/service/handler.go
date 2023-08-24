package service

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-interaction/internal/dao/mysql"
	"github.com/Happy-Why/toktik-interaction/internal/dao/redis"
	"github.com/Happy-Why/toktik-interaction/internal/model"
	"github.com/Happy-Why/toktik-interaction/internal/model/auto"
	"github.com/Happy-Why/toktik-interaction/internal/repo"
	"github.com/Happy-Why/toktik-interaction/pkg/myerr"
	"github.com/Happy-Why/toktik-interaction/pkg/rpc/client"
	inter "github.com/Happy-Why/toktik-rpc/kitex_gen/interaction"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
	cache           repo.Cache
	interactionRepo repo.InteractionRepo
	respRepo        repo.HandlerResp
	transaction     mysql.Transaction
}

func NewInteractionService() *InteractionServiceImpl {
	return &InteractionServiceImpl{
		cache:           redis.GetRdbCache(),
		interactionRepo: mysql.NewInteractionDao(),
		respRepo:        repo.NewHandlerResps(),
		transaction:     mysql.NewTransaction(),
	}
}

func (is *InteractionServiceImpl) FollowSB(ctx context.Context, req *inter.FollowActionRequest) (resp *inter.FollowActionResponse, err error) {
	if req.MyUserId == req.ToUserId {
		return is.respRepo.FollowSBResponse(myerr.NotCanFollowSelf, model.MsgNil, &inter.FollowActionResponse{}), nil
	}
	// 逻辑校验，在 relation 表中查询是否存在 关注关系
	relationID, exist, err := is.interactionRepo.IsRelationExist(ctx, req.MyUserId, req.ToUserId)
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("is.interactionRepo.IsRelationExist err:", zap.Error(err))
		return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
	}
	relationInfo := &auto.Relation{
		Model:    gorm.Model{},
		UserId:   uint(req.MyUserId),
		TargetId: uint(req.ToUserId),
	}
	// 进行业务
	switch req.ActionType {
	case model.Follow:
		if exist {
			return is.respRepo.FollowSBResponse(myerr.UserAlreadyFollowed, model.MsgNil, &inter.FollowActionResponse{}), nil
		}
		// 进行事务，关注后，将当前用户的follow_count 和 对方用户的follower_count 增1
		if err = is.transaction.Action(func(conn mysql.DbConn) error {
			if err = is.interactionRepo.FollowUserAction(ctx, conn, relationInfo); err != nil {
				zap.L().Error("is.interactionRepo.FollowUserAction err:", zap.Error(err))
				return err
			}
			//if err = is.interactionRepo.AddFollowCount(ctx, conn, uint(req.MyUserId)); err != nil {
			//	zap.L().Error("is.interactionRepo.AddFollowCount err:", zap.Error(err))
			//	return err
			//}
			//if err = is.interactionRepo.AddFollowerCount(ctx, conn, uint(req.ToUserId)); err != nil {
			//	zap.L().Error("is.interactionRepo.AddFollowerCount err:", zap.Error(err))
			//	return err
			//}
			return nil
		}); err != nil {
			return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
		}
		addFollowResp, _ := client.UserClient.AddFollowCount(ctx, &user.AddFollowCountRequest{UserId: req.MyUserId, TargetId: req.ToUserId})
		if addFollowResp == nil {
			fmt.Printf("addFpllowResp:%#v\n", addFollowResp)
			zap.L().Error("client.InteractionClient.IsFollowTarget 返回空指针")
			return is.respRepo.FollowSBResponse(errcode.ErrServer, model.MsgNil, &inter.FollowActionResponse{}), nil
		}
		if addFollowResp.StatusCode != model.RpcSuccess {
			zap.L().Error(" client.InteractionClient.IsFollowTarget err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.CreateErr(addFollowResp.StatusCode, model.MsgNil), addFollowResp.StatusMsg, &inter.FollowActionResponse{}), nil
		}
	case model.CancelFollow:
		if !exist {
			return is.respRepo.FollowSBResponse(myerr.UserNotFollowed, model.MsgNil, &inter.FollowActionResponse{}), nil
		}
		relationInfo.ID = relationID
		// 进行事务，关注后，将当前用户的follow_count 和 对方用户的follower_count 增1
		if err = is.transaction.Action(func(conn mysql.DbConn) error {
			if err = is.interactionRepo.CancelFollowUser(ctx, conn, relationInfo); err != nil {
				zap.L().Error("is.interactionRepo.CancelFollowUser err:", zap.Error(err))
				return nil
			}
			//if err = is.interactionRepo.SubFollowCount(ctx, conn, uint(req.MyUserId)); err != nil {
			//	zap.L().Error("is.interactionRepo.SubFollowCount err:", zap.Error(err))
			//	return err
			//}
			//if err = is.interactionRepo.SubFollowerCount(ctx, conn, uint(req.ToUserId)); err != nil {
			//	zap.L().Error("is.interactionRepo.SubFollowerCount err:", zap.Error(err))
			//	return err
			//}
			return nil
		}); err != nil {
			return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
		}
		subFollowResp, _ := client.UserClient.SubFollowCount(ctx, &user.SubFollowCountRequest{UserId: req.MyUserId, TargetId: req.ToUserId})
		if subFollowResp == nil {
			zap.L().Error("client.InteractionClient.IsFollowTarget 返回空指针")
			return is.respRepo.FollowSBResponse(errcode.ErrServer, model.MsgNil, &inter.FollowActionResponse{}), nil
		}
		if subFollowResp.StatusCode != model.RpcSuccess {
			zap.L().Error(" client.InteractionClient.IsFollowTarget err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.CreateErr(subFollowResp.StatusCode, model.MsgNil), subFollowResp.StatusMsg, &inter.FollowActionResponse{}), nil
		}
	default:
		return is.respRepo.FollowSBResponse(myerr.OperationErr, model.MsgNil, &inter.FollowActionResponse{}), nil
	}
	// 模型转换，返回参数
	return is.respRepo.FollowSBResponse(errcode.StatusOK, model.MsgNil, &inter.FollowActionResponse{}), nil
}

func (is *InteractionServiceImpl) FollowList(ctx context.Context, req *inter.FollowListRequest) (resp *inter.FollowListResponse, err error) {
	// 没有什么逻辑需要校验
	// 直接进行业务
	// 在 relation表中查询
	userIDs, err := is.interactionRepo.GetFollowIDs(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
		return is.respRepo.FollowListResponse(errcode.ErrDB, err.Error(), &inter.FollowListResponse{}), nil
	}
	if len(userIDs) == 0 {
		return is.respRepo.FollowListResponse(errcode.StatusOK, model.MsgNil, &inter.FollowListResponse{}), nil
	}
	userListResp, err := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{TargetId: userIDs})
	if userListResp == nil {
		zap.L().Error("client.UserClient.GetUserList 返回空指针")
		return is.respRepo.FollowListResponse(errcode.ErrServer, model.MsgNil, &inter.FollowListResponse{}), nil
	}
	if err != nil {
		zap.L().Error("client.UserClient.GetUserList err:", zap.Error(err))
		return is.respRepo.FollowListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &inter.FollowListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FollowListResponse)
	for _, v := range userListResp.UserList {
		resp.UserList = append(resp.UserList, v)
	}
	return is.respRepo.FollowListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
func (is *InteractionServiceImpl) FansList(ctx context.Context, req *inter.FansListRequest) (resp *inter.FansListResponse, err error) {
	userIDs, err := is.interactionRepo.GetFansIDs(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
		return is.respRepo.FansListResponse(errcode.ErrDB, err.Error(), &inter.FansListResponse{}), nil
	}
	if len(userIDs) == 0 {
		return is.respRepo.FansListResponse(errcode.StatusOK, model.MsgNil, &inter.FansListResponse{}), nil
	}
	userListResp, err := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{TargetId: userIDs})
	if userListResp == nil {
		zap.L().Error("client.UserClient.GetUserList 返回空指针")
		return is.respRepo.FansListResponse(errcode.ErrServer, model.MsgNil, &inter.FansListResponse{}), nil
	}
	if err != nil {
		zap.L().Error("client.UserClient.GetUserList err:", zap.Error(err))
		return is.respRepo.FansListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &inter.FansListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FansListResponse)
	for _, v := range userListResp.UserList {
		resp.UserList = append(resp.UserList, v)
	}
	return is.respRepo.FansListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) FriendList(ctx context.Context, req *inter.FriendListRequest) (resp *inter.FriendListResponse, err error) {
	// 没什么逻辑
	userIDs, err := is.interactionRepo.GetFriendIDs(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
		return is.respRepo.FriendListResponse(errcode.ErrDB, err.Error(), &inter.FriendListResponse{}), nil
	}
	if len(userIDs) == 0 {
		return is.respRepo.FriendListResponse(errcode.StatusOK, model.MsgNil, &inter.FriendListResponse{}), nil
	}
	userListResp, err := client.UserClient.GetUserList(ctx, &user.GetUserListRequest{TargetId: userIDs})
	if userListResp == nil {
		zap.L().Error("client.UserClient.GetUserList 返回空指针")
		return is.respRepo.FriendListResponse(errcode.ErrServer, model.MsgNil, &inter.FriendListResponse{}), nil
	}
	if err != nil {
		zap.L().Error("client.UserClient.GetUserList err:", zap.Error(err))
		return is.respRepo.FriendListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &inter.FriendListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FriendListResponse)
	for _, v := range userListResp.UserList {
		resp.UserList = append(resp.UserList, v)
	}
	return is.respRepo.FriendListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) IsFollowTarget(ctx context.Context, req *inter.IsFollowTargetRequest) (resp *inter.IsFollowTargetResponse, err error) {
	// 处理业务
	// 到relation表中 查询 userId 是否关注 targetId
	exist, err := is.interactionRepo.IsFollowUser(ctx, req.UserId, req.TargetId)
	if err != nil {
		zap.L().Error("is.interactionRepo.IsFollowUser err:", zap.Error(err))
		return is.respRepo.IsFollowTargetResponse(errcode.ErrDB, err.Error(), &inter.IsFollowTargetResponse{}), nil
	}
	// 模型转换
	resp = &inter.IsFollowTargetResponse{
		Exist: exist,
	}
	return is.respRepo.IsFollowTargetResponse(errcode.StatusOK, model.MsgNil, &inter.IsFollowTargetResponse{}), nil
}
