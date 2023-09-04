package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
	"toktik-common/errcode"
	"toktik-interaction/internal/cache"
	"toktik-interaction/internal/dao/mysql"
	"toktik-interaction/internal/model"
	"toktik-interaction/internal/model/auto"
	"toktik-interaction/internal/repo"
	"toktik-interaction/pkg/myerr"
	"toktik-interaction/pkg/rpc/client"
	inter "toktik-rpc/kitex_gen/interaction"
	"toktik-rpc/kitex_gen/user"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
	interactionRepo repo.InteractionRepo
	respRepo        repo.HandlerResp
	transaction     mysql.Transaction
	rCache          repo.RCacheRepo
}

func NewInteractionService() *InteractionServiceImpl {
	return &InteractionServiceImpl{
		interactionRepo: mysql.NewInteractionDao(),
		respRepo:        repo.NewHandlerResps(),
		transaction:     mysql.NewTransaction(),
		rCache:          cache.NewInteractionrCache(),
	}
}

func (is *InteractionServiceImpl) FollowSB(ctx context.Context, req *inter.FollowActionRequest) (resp *inter.FollowActionResponse, err error) {
	if req.MyUserId == req.ToUserId {
		return is.respRepo.FollowSBResponse(myerr.NotCanFollowSelf, model.MsgNil, &inter.FollowActionResponse{}), nil
	}
	// ① 逻辑校验，在 relation 表中查询是否存在 关注关系
	relationID, exist, err := is.interactionRepo.IsRelationExist(ctx, req.MyUserId, req.ToUserId)
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("is.interactionRepo.IsRelationExist err:", zap.Error(err))
		return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
	}
	relationInfo := &auto.Relation{
		BaseModel: auto.BaseModel{},
		UserId:    uint(req.MyUserId),
		TargetId:  uint(req.ToUserId),
	}
	// 进行业务

	switch req.ActionType {
	case model.Follow:
		if exist {
			return is.respRepo.FollowSBResponse(myerr.UserAlreadyFollowed, model.MsgNil, &inter.FollowActionResponse{}), nil
		}
		// goroutine 写法 先处理 follow
		//follow := make(chan errcode.RespErr)
		//go func() {
		//	// 先删除缓存
		//	err = is.rCache.SDelFollow(ctx, followKey, req.ToUserId)
		//	if err != nil {
		//		zap.L().Error("is.rCache.SDelFollow err:", zap.Error(err))
		//		follow <- errcode.RespErr{Err:errcode.ErrRedis, ErrStr: err.Error()}
		//	}
		//	// 进行事务，关注后，将当前用户的follow_count 和 对方用户的follower_count 增1
		//	if err = is.transaction.Action(func(conn mysql.DbConn) error {
		//		if err = is.interactionRepo.FollowUserAction(ctx, conn, relationInfo); err != nil {
		//			zap.L().Error("is.interactionRepo.FollowUserAction err:", zap.Error(err))
		//			return err
		//		}
		//		return nil
		//	}); err != nil {
		//		follow <- errcode.RespErr{Err:errcode.ErrDB, ErrStr: err.Error()}
		//	}
		//	// 延迟删除缓存，立即删除缓存，请求2 先读取数据库中的旧值，会向缓存中写入旧值，所以需要等待一次redis写入时间再删除
		//	time.Sleep(300 * time.Millisecond)
		//	err = is.rCache.SDelFollow(ctx, followKey, req.ToUserId)
		//	if err != nil {
		//		zap.L().Error("is.rCache.SDelFollow err:", zap.Error(err))
		//		follow <- errcode.RespErr{Err:errcode.ErrRedis, ErrStr: err.Error()}
		//	}
		//	follow <- errcode.RespErr{Err:errcode.StatusOK, ErrStr: model.MsgNil}
		//}()
		//respErr := <- follow
		//if respErr.Err.ECode() != 0{
		//	return is.respRepo.FollowSBResponse(respErr.Err, respErr.ErrStr, &inter.FollowActionResponse{}), nil
		//}
		// ② 先删除缓存 ，删除自己的关注缓存
		followKey := auto.CreateFollowKey(req.MyUserId)
		err = is.rCache.DelFollow(ctx, followKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFollow err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// 删除 对方的粉丝缓存
		fanKey := auto.CreateFanKey(req.ToUserId)
		err = is.rCache.DelFan(ctx, fanKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFan err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// ③ 进行事务，关注后，将当前用户的follow_count 和 对方用户的follower_count 增1
		if err = is.transaction.Action(func(conn mysql.DbConn) error {
			if err = is.interactionRepo.FollowUserAction(ctx, conn, relationInfo); err != nil {
				zap.L().Error("is.interactionRepo.FollowUserAction err:", zap.Error(err))
				return err
			}
			return nil
		}); err != nil {
			return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// ④ 延迟删除缓存，立即删除缓存的话,请求2 先读取数据库中的旧值，会向缓存中写入旧值，所以需要等待一次redis写入时间再删除
		time.Sleep(50 * time.Millisecond)
		err = is.rCache.DelFollow(ctx, followKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFollow err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		err = is.rCache.DelFan(ctx, fanKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFan err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// ⑤ 判断对方是否关注自己，是则删除好友缓存
		isFollowMe, err := is.IsFollowTarget2(ctx, req.ToUserId, req.MyUserId)
		if err != nil {
			zap.L().Error("is.IsFollowTarget2 err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		if isFollowMe {
			friendKey1 := auto.CreateFriendKey(req.MyUserId)
			if err = is.rCache.DelFriend(ctx, friendKey1); err != nil {
				zap.L().Error("is.rCache.DelFriend err:", zap.Error(err))
				return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
			}
			friendKey2 := auto.CreateFriendKey(req.ToUserId)
			if err = is.rCache.DelFriend(ctx, friendKey2); err != nil {
				zap.L().Error("is.rCache.DelFriend err:", zap.Error(err))
				return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
			}
		}
		// ⑥ 计数
		addFollowResp, _ := client.UserCache.AddFollowCount(ctx, &user.AddFollowCountRequest{UserId: req.MyUserId, TargetId: req.ToUserId})
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
		// 删除自己的关注缓存
		followKey := auto.CreateFollowKey(req.MyUserId)
		err = is.rCache.DelFollow(ctx, followKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFollow err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// 删除对方的粉丝缓存
		fanKey := auto.CreateFanKey(req.ToUserId)
		err = is.rCache.DelFan(ctx, fanKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFan err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// 如果是好友，则删除双方好友缓存
		isFollowMe, err := is.IsFollowTarget2(ctx, req.ToUserId, req.MyUserId)
		if err != nil {
			zap.L().Error("is.IsFollowTarget2 err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		if isFollowMe {
			// 删除好友缓存,双向删除
			friendKey1 := auto.CreateFriendKey(req.MyUserId)
			err = is.rCache.DelFriend(ctx, friendKey1)
			if err != nil {
				zap.L().Error("is.rCache.DelFriend err:", zap.Error(err))
				return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
			}
			friendKey2 := auto.CreateFriendKey(req.ToUserId)
			err = is.rCache.DelFriend(ctx, friendKey2)
			if err != nil {
				zap.L().Error("is.rCache.DelFriend err:", zap.Error(err))
				return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
			}
		}
		// 进行事务，关注后，将当前用户的follow_count 和 对方用户的follower_count -1
		if err = is.transaction.Action(func(conn mysql.DbConn) error {
			if err = is.interactionRepo.CancelFollowUser(ctx, conn, relationInfo); err != nil {
				zap.L().Error("is.interactionRepo.CancelFollowUser err:", zap.Error(err))
				return nil
			}
			return nil
		}); err != nil {
			return is.respRepo.FollowSBResponse(errcode.ErrDB, err.Error(), &inter.FollowActionResponse{}), nil
		}
		time.Sleep(50 * time.Millisecond)
		err = is.rCache.DelFollow(ctx, followKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFollow err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		err = is.rCache.DelFan(ctx, fanKey)
		if err != nil {
			zap.L().Error("is.rCache.SDelFan err:", zap.Error(err))
			return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
		}
		// 如果是好友，则删除双方好友缓存
		if isFollowMe {
			// 删除好友缓存,双向删除
			friendKey1 := auto.CreateFriendKey(req.MyUserId)
			err = is.rCache.DelFriend(ctx, friendKey1)
			if err != nil {
				zap.L().Error("is.rCache.DelFriend err:", zap.Error(err))
				return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
			}
			friendKey2 := auto.CreateFriendKey(req.ToUserId)
			err = is.rCache.DelFriend(ctx, friendKey2)
			if err != nil {
				zap.L().Error("is.rCache.DelFriend err:", zap.Error(err))
				return is.respRepo.FollowSBResponse(errcode.ErrRedis, err.Error(), &inter.FollowActionResponse{}), nil
			}
		}
		subFollowResp, _ := client.UserCache.SubFollowCount(ctx, &user.SubFollowCountRequest{UserId: req.MyUserId, TargetId: req.ToUserId})
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
	// 先查缓存
	followKey := auto.CreateFollowKey(req.UserId)
	followIds, err := is.rCache.SGetAllIds(ctx, followKey)
	if err != nil {
		zap.L().Error("is.rCache.SGetAllIds err:", zap.Error(err))
		return is.respRepo.FollowListResponse(errcode.ErrRedis, err.Error(), &inter.FollowListResponse{}), nil
	}
	if followIds == nil {
		// 缓存没查到，查数据库
		followIds, err = is.interactionRepo.GetFollowIDs(ctx, uint(req.UserId))
		if err != nil {
			zap.L().Error("is.interactionRepo.GetFollowList err:", zap.Error(err))
			return is.respRepo.FollowListResponse(errcode.ErrDB, err.Error(), &inter.FollowListResponse{}), nil
		}
		// 没有关注的人，直接返回
		if len(followIds) == 0 {
			return is.respRepo.FollowListResponse(errcode.StatusOK, model.MsgNil, &inter.FollowListResponse{}), nil
		}
		// 添加缓存
		if err = is.rCache.SAddManyIds(ctx, followKey, followIds); err != nil {
			zap.L().Error("is.rCache.SAddManyIds err:", zap.Error(err))
			return is.respRepo.FollowListResponse(errcode.ErrRedis, err.Error(), &inter.FollowListResponse{}), nil
		}
	}

	userListResp, err := client.UserCache.GetUserList(ctx, &user.GetUserListRequest{UserId: req.MyUserId, TargetId: followIds})
	if userListResp == nil {
		zap.L().Error("client.UserCache.GetUserList 返回空指针")
		return is.respRepo.FollowListResponse(errcode.ErrServer, model.MsgNil, &inter.FollowListResponse{}), nil
	}
	if userListResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.UserCache.GetUserList err:", zap.Error(err))
		return is.respRepo.FollowListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &inter.FollowListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FollowListResponse)
	resp.UserList = make([]*user.User, len(userListResp.UserList))
	for i, v := range userListResp.UserList {
		resp.UserList[i] = v
	}
	return is.respRepo.FollowListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) FansList(ctx context.Context, req *inter.FansListRequest) (resp *inter.FansListResponse, err error) {
	// 先查缓存
	fanKey := auto.CreateFanKey(req.UserId)
	fanIds, err := is.rCache.SGetAllIds(ctx, fanKey)
	if err != nil {
		zap.L().Error("is.rCache.SGetAllIds err:", zap.Error(err))
		return is.respRepo.FansListResponse(errcode.ErrRedis, err.Error(), &inter.FansListResponse{}), nil
	}
	if fanIds == nil {
		// 缓存没查到，查数据库
		fanIds, err = is.interactionRepo.GetFansIDs(ctx, uint(req.UserId))
		if err != nil {
			zap.L().Error("is.interactionRepo.GetFansIDs err:", zap.Error(err))
			return is.respRepo.FansListResponse(errcode.ErrDB, err.Error(), &inter.FansListResponse{}), nil
		}
		// 没有粉丝，直接返回
		if len(fanIds) == 0 {
			return is.respRepo.FansListResponse(errcode.StatusOK, model.MsgNil, &inter.FansListResponse{}), nil
		}
		// 添加缓存
		if err = is.rCache.SAddManyIds(ctx, fanKey, fanIds); err != nil {
			zap.L().Error("is.rCache.SAddManyIds err:", zap.Error(err))
			return is.respRepo.FansListResponse(errcode.ErrRedis, err.Error(), &inter.FansListResponse{}), nil
		}
	}
	userListResp, err := client.UserCache.GetUserList(ctx, &user.GetUserListRequest{UserId: req.MyUserId, TargetId: fanIds})
	if userListResp == nil {
		zap.L().Error("client.UserCache.GetUserList 返回空指针")
		return is.respRepo.FansListResponse(errcode.ErrServer, model.MsgNil, &inter.FansListResponse{}), nil
	}
	if userListResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.UserCache.GetUserList err:", zap.Error(err))
		return is.respRepo.FansListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &inter.FansListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FansListResponse)
	resp.UserList = make([]*user.User, len(fanIds))
	for i, v := range userListResp.UserList {
		resp.UserList[i] = v
	}
	return is.respRepo.FansListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) FriendList(ctx context.Context, req *inter.FriendListRequest) (resp *inter.FriendListResponse, err error) {
	// 没什么逻辑
	// 先查缓存
	friendKey := auto.CreateFriendKey(req.UserId)
	friendIds, err := is.rCache.SGetAllIds(ctx, friendKey)
	if err != nil {
		zap.L().Error("is.rCache.SGetAllIds err:", zap.Error(err))
		return is.respRepo.FriendListResponse(errcode.ErrRedis, err.Error(), &inter.FriendListResponse{}), nil
	}
	if friendIds == nil {
		// 缓存没查到，查数据库
		friendIds, err = is.interactionRepo.GetFriendIDs(ctx, uint(req.UserId))
		if err != nil {
			zap.L().Error("is.interactionRepo.GetFriendIDs err:", zap.Error(err))
			return is.respRepo.FriendListResponse(errcode.ErrDB, err.Error(), &inter.FriendListResponse{}), nil
		}
		// 没有好友，直接返回
		if len(friendIds) == 0 {
			return is.respRepo.FriendListResponse(errcode.StatusOK, model.MsgNil, &inter.FriendListResponse{}), nil
		}
		// 添加缓存
		if err = is.rCache.SAddManyIds(ctx, friendKey, friendIds); err != nil {
			zap.L().Error("is.rCache.SAddManyIds err:", zap.Error(err))
			return is.respRepo.FriendListResponse(errcode.ErrRedis, err.Error(), &inter.FriendListResponse{}), nil
		}
	}

	userListResp, err := client.UserCache.GetUserList(ctx, &user.GetUserListRequest{UserId: req.MyUserId, TargetId: friendIds})
	if userListResp == nil {
		zap.L().Error("client.UserCache.GetUserList 返回空指针")
		return is.respRepo.FriendListResponse(errcode.ErrServer, model.MsgNil, &inter.FriendListResponse{}), nil
	}
	if userListResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.UserCache.GetUserList err:", zap.Error(err))
		return is.respRepo.FriendListResponse(errcode.CreateErr(userListResp.StatusCode, model.MsgNil), userListResp.StatusMsg, &inter.FriendListResponse{}), nil
	}
	// 模型转换，返回
	resp = new(inter.FriendListResponse)
	resp.UserList = make([]*user.User, len(userListResp.UserList))
	for i, v := range userListResp.UserList {
		resp.UserList[i] = v
	}
	return is.respRepo.FriendListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) IsFollowTarget(ctx context.Context, req *inter.IsFollowTargetRequest) (resp *inter.IsFollowTargetResponse, err error) {
	resp = new(inter.IsFollowTargetResponse)
	if req.TargetId == req.UserId {
		resp.Exist = false
		return is.respRepo.IsFollowTargetResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	// 处理业务
	// 先查缓存，查不到 查数据库
	followKey := auto.CreateFollowKey(req.UserId)
	exist, err := is.rCache.IsFollow(ctx, followKey, req.TargetId)
	if err != nil {
		zap.L().Error("is.rCache.IsFollow err:", zap.Error(err))
		return is.respRepo.IsFollowTargetResponse(errcode.ErrRedis, err.Error(), &inter.IsFollowTargetResponse{}), nil
	}
	// 模型转换
	resp.Exist = exist
	if exist {
		return is.respRepo.IsFollowTargetResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	keyExist, err := is.rCache.KeyExist(ctx, followKey)
	if err != nil {
		zap.L().Error("is.rCache.KeyExist err:", zap.Error(err))
		return is.respRepo.IsFollowTargetResponse(errcode.ErrRedis, err.Error(), &inter.IsFollowTargetResponse{}), nil
	}
	if keyExist { // 缓存未过期
		return is.respRepo.IsFollowTargetResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	// 因为缓存过期，到relation表中 查询 userId 是否关注 targetId
	exist, err = is.interactionRepo.IsFollowUser(ctx, req.UserId, req.TargetId)
	if err != nil {
		zap.L().Error("is.interactionRepo.IsFollowUser err:", zap.Error(err))
		return is.respRepo.IsFollowTargetResponse(errcode.ErrDB, err.Error(), &inter.IsFollowTargetResponse{}), nil
	}
	//TODO 设置缓存，可以用消息队列异步进行 , 去数据库中查询该用户关注的所有人，然后设置为缓存
	targetIds, err := is.interactionRepo.GetFollowIDs(ctx, uint(req.UserId))
	if err != nil {
		zap.L().Error("is.interactionRepo.GetFollowIDs err:", zap.Error(err))
		return is.respRepo.IsFollowTargetResponse(errcode.ErrDB, err.Error(), &inter.IsFollowTargetResponse{}), nil
	}
	if targetIds == nil { // 没有id可以缓存
		return is.respRepo.IsFollowTargetResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	err = is.rCache.SAddManyIds(ctx, followKey, targetIds)
	if err != nil {
		zap.L().Error("is.rCache.SAddAllFollow err:", zap.Error(err))
		return is.respRepo.IsFollowTargetResponse(errcode.ErrRedis, err.Error(), &inter.IsFollowTargetResponse{}), nil
	}
	// 模型转换
	resp = &inter.IsFollowTargetResponse{
		Exist: exist,
	}
	return is.respRepo.IsFollowTargetResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) IsFollowManyTargets(ctx context.Context, req *inter.IsFollowManyTargetsRequest) (resp *inter.IsFollowManyTargetsResponse, err error) {
	resp = new(inter.IsFollowManyTargetsResponse)
	resp.ManyExist = make([]bool, len(req.TargetIds))
	followKey := auto.CreateFollowKey(req.UserId)
	for i, v := range req.TargetIds {
		exist, err := is.rCache.IsFollow(ctx, followKey, v)
		if err != nil {
			zap.L().Error("is.rCache.IsFollow err:", zap.Error(err))
			return is.respRepo.IsFollowManyTargetsResponse(errcode.ErrRedis, err.Error(), &inter.IsFollowManyTargetsResponse{}), nil
		}
		resp.ManyExist[i] = exist
	}
	return is.respRepo.IsFollowManyTargetsResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) IsFriend(ctx context.Context, req *inter.IsFriendRequest) (resp *inter.IsFriendResponse, err error) {
	resp = new(inter.IsFriendResponse)
	// 查缓存
	friendKey := auto.CreateFriendKey(req.UserId)
	ok, err := is.rCache.IsFriend(ctx, friendKey, req.TargetId)
	if err != nil {
		zap.L().Error(" is.rCache.IsFriend err:", zap.Error(err))
		return is.respRepo.IsFriendResponse(errcode.ErrDB, err.Error(), &inter.IsFriendResponse{}), nil
	}
	// 模型转换
	resp.IsFriend = ok
	if ok {
		return is.respRepo.IsFriendResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	// ok == false 缓存过期 或者 不是好友
	keyExist, err := is.rCache.KeyExist(ctx, friendKey)
	if err != nil {
		zap.L().Error("is.rCache.KeyExist err:", zap.Error(err))
		return is.respRepo.IsFriendResponse(errcode.ErrRedis, err.Error(), &inter.IsFriendResponse{}), nil
	}
	if keyExist { // 缓存未过期
		return is.respRepo.IsFriendResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}
	// 缓存未命中，查数据库
	isFriend, err := is.interactionRepo.IsFriend(ctx, req.UserId, req.TargetId)
	if err != nil {
		zap.L().Error("is.interactionRepo.IsFriend err:", zap.Error(err))
		return is.respRepo.IsFriendResponse(errcode.ErrDB, err.Error(), &inter.IsFriendResponse{}), nil
	}
	resp.IsFriend = isFriend
	// 添加缓存
	err = is.rCache.SAddFriend(ctx, friendKey, req.TargetId)
	if err != nil {
		zap.L().Error("is.rCache.SAddFriend err:", zap.Error(err))
		return is.respRepo.IsFriendResponse(errcode.ErrRedis, err.Error(), &inter.IsFriendResponse{}), nil
	}
	return is.respRepo.IsFriendResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (is *InteractionServiceImpl) IsFollowTarget2(ctx context.Context, userId, targetId int64) (bool, error) {
	// 处理业务
	// 先查缓存，查不到 查数据库
	followKey := auto.CreateFollowKey(userId)
	exist, err := is.rCache.IsFollow(ctx, followKey, targetId)
	if err != nil {
		zap.L().Error("is.rCache.IsFollow err:", zap.Error(err))
		return false, err
	}
	if exist {
		return exist, nil
	}
	keyExist, err := is.rCache.KeyExist(ctx, followKey)
	if err != nil {
		zap.L().Error("is.rCache.KeyExist err:", zap.Error(err))
		return false, err
	}
	if keyExist { // 缓存未过期
		return exist, nil
	}
	// 因为缓存过期，到relation表中 查询 userId 是否关注 targetId
	exist, err = is.interactionRepo.IsFollowUser(ctx, userId, targetId)
	if err != nil {
		zap.L().Error("is.interactionRepo.IsFollowUser err:", zap.Error(err))
		return false, err
	}
	// 如果查到有关注，则添加缓存
	if exist {
		if err = is.rCache.SAddFollow(ctx, followKey, targetId); err != nil {
			zap.L().Error("is.rCache.SAddFollow err:", zap.Error(err))
			return false, err
		}
	}
	return exist, nil
}
