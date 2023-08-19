package service

import (
	"context"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-interaction/internal/dao"
	"github.com/Happy-Why/toktik-interaction/internal/dao/mysql"
	"github.com/Happy-Why/toktik-interaction/internal/model"
	"github.com/Happy-Why/toktik-interaction/internal/model/auto"
	"github.com/Happy-Why/toktik-interaction/internal/repo"
	"github.com/Happy-Why/toktik-interaction/pkg/myerr"
	inter "github.com/Happy-Why/toktik-rpc/kitex_gen/interaction"
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
		cache:           dao.Group.Rdb,
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
	// TODO 将 “热用户” follow_count 和 follower_count 缓存在redis中，保持缓存一致性，或者先保存到redis中，然后定时任务存放到 mysql中
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
			if err = is.interactionRepo.AddFollowCount(ctx, conn, uint(req.MyUserId)); err != nil {
				zap.L().Error("is.interactionRepo.AddFollowCount err:", zap.Error(err))
				return err
			}
			if err = is.interactionRepo.AddFollowerCount(ctx, conn, uint(req.ToUserId)); err != nil {
				zap.L().Error("is.interactionRepo.AddFollowerCount err:", zap.Error(err))
				return err
			}
			return nil
		}); err != nil {
			return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
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
			if err = is.interactionRepo.SubFollowCount(ctx, conn, uint(req.MyUserId)); err != nil {
				zap.L().Error("is.interactionRepo.SubFollowCount err:", zap.Error(err))
				return err
			}
			if err = is.interactionRepo.SubFollowerCount(ctx, conn, uint(req.ToUserId)); err != nil {
				zap.L().Error("is.interactionRepo.SubFollowerCount err:", zap.Error(err))
				return err
			}
			return nil
		}); err != nil {
			return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
		}

	}
	// 模型转换，返回参数
	return is.respRepo.FollowSBResponse(errcode.StatusOK, model.MsgNil, &inter.FollowActionResponse{}), nil
}

func (is *InteractionServiceImpl) FollowList(ctx context.Context, req *inter.FollowListRequest) (resp *inter.FollowListResponse, err error) {
	// 没有什么逻辑需要校验
	// 直接进行业务
	// 在 relation表中查询
	List, err := is.interactionRepo.GetFollowList(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
		return is.respRepo.FollowListResponse(errcode.ErrDB, err.Error(), &inter.FollowListResponse{}), nil
	}
	if List == nil {
		return is.respRepo.FollowListResponse(errcode.StatusOK, model.MsgNil, &inter.FollowListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FollowListResponse)
	resp.UserList = make([]*inter.User, 0)
	for _, v := range List {
		resp.UserList = append(resp.UserList, &inter.User{
			Id:              int64(v.ID),
			Name:            v.Username,
			FollowCount:     &v.FollowCount,
			FollowerCount:   &v.FollowerCount,
			IsFollow:        v.IsFollow,
			Avatar:          &v.Avatar,
			BackgroundImage: &v.BackgroundImage,
			Signature:       &v.Signature,
			TotalFavorited:  &v.TotalFavorited,
			WorkCount:       &v.WorkCount,
			FavoriteCount:   &v.TotalFavorited,
		})
	}
	return is.respRepo.FollowListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
func (is *InteractionServiceImpl) FansList(ctx context.Context, req *inter.FansListRequest) (resp *inter.FansListResponse, err error) {
	List, err := is.interactionRepo.GetFansList(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
		return is.respRepo.FansListResponse(errcode.ErrDB, err.Error(), &inter.FansListResponse{}), nil
	}
	if len(List) == 0 {
		return is.respRepo.FansListResponse(errcode.StatusOK, model.MsgNil, &inter.FansListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FansListResponse)
	resp.UserList = make([]*inter.User, 0)
	for _, v := range List {
		resp.UserList = append(resp.UserList, &inter.User{
			Id:              int64(v.ID),
			Name:            v.Username,
			FollowCount:     &v.FollowCount,
			FollowerCount:   &v.FollowerCount,
			IsFollow:        v.IsFollow,
			Avatar:          &v.Avatar,
			BackgroundImage: &v.BackgroundImage,
			Signature:       &v.Signature,
			TotalFavorited:  &v.TotalFavorited,
			WorkCount:       &v.WorkCount,
			FavoriteCount:   &v.TotalFavorited,
		})
	}
	return is.respRepo.FansListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) FriendList(ctx context.Context, req *inter.FriendListRequest) (resp *inter.FriendListResponse, err error) {
	// 没什么逻辑
	List, err := is.interactionRepo.GetFriendList(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
		return is.respRepo.FriendListResponse(errcode.ErrDB, err.Error(), &inter.FriendListResponse{}), nil
	}
	if len(List) == 0 {
		return is.respRepo.FriendListResponse(errcode.StatusOK, model.MsgNil, &inter.FriendListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FriendListResponse)
	resp.UserList = make([]*inter.User, 0)
	for _, v := range List {
		resp.UserList = append(resp.UserList, &inter.User{
			Id:              int64(v.ID),
			Name:            v.Username,
			FollowCount:     &v.FollowCount,
			FollowerCount:   &v.FollowerCount,
			IsFollow:        v.IsFollow,
			Avatar:          &v.Avatar,
			BackgroundImage: &v.BackgroundImage,
			Signature:       &v.Signature,
			TotalFavorited:  &v.TotalFavorited,
			WorkCount:       &v.WorkCount,
			FavoriteCount:   &v.TotalFavorited,
		})
	}
	return is.respRepo.FriendListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
