package service

import (
	"context"
	"errors"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
	"toktik-common/errcode"
	"toktik-common/utils"
	"toktik-rpc/kitex_gen/interaction"
	"toktik-rpc/kitex_gen/user"
	"toktik-user/internal/cache"
	"toktik-user/internal/dao/mysql"
	"toktik-user/internal/global"
	"toktik-user/internal/model"
	"toktik-user/internal/model/auto"
	"toktik-user/internal/repo"
	"toktik-user/pkg/myerr"
	"toktik-user/pkg/rpc/client"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	respRepo    repo.HandlerResp
	UserRepo    repo.UserRepo
	rCache      repo.RCacheRepo
	transaction mysql.Transaction
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{
		respRepo:    repo.NewHandlerResps(),
		UserRepo:    mysql.NewUserDao(),
		rCache:      cache.NewUserCache(),
		transaction: mysql.NewTransaction(),
	}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1.可以校验参数
	// 2.校验业务逻辑(邮箱、账号、手机号是否被注册)
	exist, err := s.UserRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		zap.L().Error("s.UserRepo.GetUserByUsername err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrDB, err.Error(), &user.RegisterResponse{}), nil
	}
	if exist {
		return s.respRepo.RegisterResponse(myerr.UsernameExist, model.MsgNil, &user.RegisterResponse{}), nil
	}
	// 3.执行业务
	// 生成唯一id，存入用户名和密码
	userID, _ := global.SnowFlake.NextId()
	timeNow := time.Now()
	hashPassword, _ := utils.HashPassword(req.Password)
	userInfo := &auto.User{
		BaseModel:       auto.BaseModel{ID: uint(userID), CreatedAt: timeNow, UpdatedAt: timeNow},
		Username:        req.Username,
		Password:        hashPassword,
		Avatar:          global.Settings.Rules.DefaultUserAvatar,
		IsFollow:        false,
		BackgroundImage: global.Settings.Rules.DefaultUserBackGroundImage,
		Signature:       global.Settings.Rules.DefaultUserSignature,
	}
	userCountInfo := &auto.UserCount{
		BaseModel:      auto.BaseModel{CreatedAt: timeNow, UpdatedAt: timeNow},
		UserId:         uint(userID),
		FollowCount:    0,
		FollowerCount:  0,
		TotalFavorited: 0,
		WorkCount:      0,
		FavoriteCount:  0,
	}
	// 将 userInfo 存储 数据库
	if err = s.transaction.Action(func(conn mysql.DbConn) error {
		if err = s.UserRepo.SetUserInfo(ctx, conn, userInfo); err != nil {
			zap.L().Error("s.UserRepo.Useregister err:", zap.Error(err))
			return err
		}
		if err = s.UserRepo.SetUserCountInfo(ctx, conn, userCountInfo); err != nil {
			zap.L().Error("s.UserRepo.SetUserCountInfo err:", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return s.respRepo.RegisterResponse(errcode.ErrDB, err.Error(), &user.RegisterResponse{}), nil
	}
	// 将 userInfo 存储 redis
	if err = s.rCache.HSetUserInfo(ctx, auto.CreateUserKey(uint(userID)), auto.CreateMapUserInfo(userInfo)); err != nil {
		zap.L().Error("s.rCache.HSetUserInfo err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrRedis, err.Error(), &user.RegisterResponse{}), nil
	}
	// 将 userCount 存储 redis
	if err = s.rCache.HSetUserCountInfo(ctx, auto.CreateUserCountKey(uint(userID)), auto.CreateMapUserCount(userCountInfo)); err != nil {
		zap.L().Error("s.rCache.HSetUserInfo err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrRedis, err.Error(), &user.RegisterResponse{}), nil
	}
	// 生成token
	token, content, err := CreateToken(userID)
	if err != nil {
		zap.L().Error("CreateToken err:", zap.Error(err))
		return s.respRepo.RegisterResponse(myerr.CreateTokenErr, err.Error(), &user.RegisterResponse{}), nil
	}
	// 将 token 中的 content 加入缓存，以便身份校验时的判断身份
	a := make(chan struct{})
	go func() {
		if err := s.rCache.SetToken(ctx, model.TypeUserToken+strconv.Itoa(int(userInfo.ID)), content, global.Settings.Token.UserTokenExp); err != nil {
			zap.L().Error("s.cache.Put err:", zap.Error(err))
		}
		a <- struct{}{}
	}()
	<-a
	// 4.返回
	resp := &user.RegisterResponse{
		UserId: userID,
		Token:  token,
	}
	return s.respRepo.RegisterResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	c, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	// 通过 username 查询 UserInfo
	userInfo, err := s.UserRepo.GetUserInfoByUsername(c, req.Username)
	if err == gorm.ErrRecordNotFound {
		zap.L().Info("用户名不存在！")
		return s.respRepo.LoginResponse(myerr.UserNotFound, model.MsgNil, &user.LoginResponse{}), nil
	}
	if err != nil {
		zap.L().Error("s.UserRepo.GetUserByUsername err:", zap.Error(err))
		return s.respRepo.LoginResponse(errcode.ErrDB, err.Error(), &user.LoginResponse{}), nil
	}
	// 判断 password
	if err = utils.CheckPassword(req.Password, userInfo.Password); err != nil {
		return s.respRepo.LoginResponse(myerr.PasswordErr, err.Error(), &user.LoginResponse{}), nil
	}
	// 生成token
	token, content, err := CreateToken(int64(userInfo.ID))
	if err != nil {
		zap.L().Error("CreateToken err:", zap.Error(err))
		return s.respRepo.LoginResponse(myerr.CreateTokenErr, err.Error(), &user.LoginResponse{}), nil
	}

	//将 token 中的 content 加入缓存，以便身份校验时的判断身份
	a := make(chan struct{}, 0)
	go func() {
		if err := s.rCache.SetToken(c, model.TypeUserToken+strconv.Itoa(int(userInfo.ID)), content, global.Settings.Token.UserTokenExp); err != nil {
			zap.L().Error("s.cache.Put err:", zap.Error(err))
		}
		a <- struct{}{}
	}()
	<-a
	// 4.返回
	resp = &user.LoginResponse{
		UserId: int64(userInfo.ID),
		Token:  token,
	}
	return s.respRepo.LoginResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) TokenVerify(ctx context.Context, req *user.TokenVerifyRequest) (resp *user.TokenVerifyResponse, err error) {
	// 从缓存中查询 token 中的 content，校验身份
	content, err := s.rCache.GetToken(ctx, model.TypeUserToken+strconv.FormatInt(req.UserId, 10))
	fmt.Println("content:", content)
	if err == redis.Nil {
		return nil, errors.New("查无此人")
	}
	if err != nil {
		zap.L().Error("TokenVerify cache get user error", zap.Error(err))
		return nil, err
	}
	return nil, nil
}

func (s *UserServiceImpl) UserIndex(ctx context.Context, req *user.UserIndexRequest) (resp *user.UserIndexResponse, err error) {
	// 1.校验逻辑
	// 通过 查询 relation表 判断是否关注该用户 来决定 is_follow字段 , 获得的 userInfo.ID 是 target_id
	var isFollow bool
	if req.UserId == req.MyUserId {
		isFollow = false
	} else {
		isFollowResp, _ := client.InteractionClient.IsFollowTarget(ctx, &interaction.IsFollowTargetRequest{UserId: req.MyUserId, TargetId: req.UserId})
		if isFollowResp == nil {
			zap.L().Error("client.InteractionClient.IsFollowTarget 返回空指针")
			return s.respRepo.UserIndexResponse(errcode.ErrServer, model.MsgNil, &user.UserIndexResponse{}), nil
		}
		if isFollowResp.StatusCode != model.RpcSuccess {
			zap.L().Error(" client.InteractionClient.IsFollowTarget err:", zap.Error(err))
			return s.respRepo.UserIndexResponse(errcode.CreateErr(isFollowResp.StatusCode, model.MsgNil), isFollowResp.StatusMsg, &user.UserIndexResponse{}), nil
		}
		isFollow = isFollowResp.Exist
	}

	// 2.查询缓存
	userKey := auto.CreateUserKey(uint(req.UserId))
	userCountKey := auto.CreateUserCountKey(uint(req.UserId))
	userInfo, err := s.rCache.HGetUserInfo(ctx, userKey)
	if err != nil {
		zap.L().Error("s.rCache.HGetUserInfo err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrRedis, err.Error(), &user.UserIndexResponse{}), nil
	}
	userCountInfo, err := s.rCache.HGetUserCountInfo(ctx, userCountKey)
	if err != nil {
		zap.L().Error("s.rCache.HGetUserInfo err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrRedis, err.Error(), &user.UserIndexResponse{}), nil
	}
	if userInfo != nil {
		fmt.Println("走缓存！！！")
		resp = &user.UserIndexResponse{
			User: &user.User{
				Id:              int64(userInfo.ID),
				Name:            userInfo.Username,
				FollowCount:     &userCountInfo.FollowCount,
				FollowerCount:   &userCountInfo.FollowerCount,
				IsFollow:        isFollow,
				Avatar:          &userInfo.Avatar,
				BackgroundImage: &userInfo.BackgroundImage,
				Signature:       &userInfo.Signature,
				TotalFavorited:  &userCountInfo.TotalFavorited,
				WorkCount:       &userCountInfo.WorkCount,
				FavoriteCount:   &userCountInfo.FavoriteCount,
			},
		}
		return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}

	// 3.进行业务
	// 根据 user_id 查询 userInfo
	userInfo, err = s.UserRepo.GetUserInfoByUserID(ctx, req.UserId)
	if err == gorm.ErrRecordNotFound {
		zap.L().Error("s.UserRepo.GetUserInfoByUserID err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(myerr.UserNotFound, model.MsgNil, &user.UserIndexResponse{}), nil
	}
	if err != nil {
		zap.L().Error("s.UserRepo.GetUserInfoByUserID err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrDB, err.Error(), &user.UserIndexResponse{}), nil
	}
	userCountInfo, err = s.UserRepo.GetUserCountInfoByUserID(ctx, req.UserId)
	if err == gorm.ErrRecordNotFound {
		zap.L().Error("s.UserRepo.GetUserCountInfoByUserID err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(myerr.UserNotFound, model.MsgNil, &user.UserIndexResponse{}), nil
	}
	if err != nil {
		zap.L().Error("s.UserRepo.GetUserCountInfoByUserID err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrDB, err.Error(), &user.UserIndexResponse{}), nil
	}
	// 模型转换
	resp = &user.UserIndexResponse{
		User: &user.User{
			Id:              int64(userInfo.ID),
			Name:            userInfo.Username,
			FollowCount:     &userCountInfo.FollowCount,
			FollowerCount:   &userCountInfo.FollowerCount,
			IsFollow:        isFollow,
			Avatar:          &userInfo.Avatar,
			BackgroundImage: &userInfo.BackgroundImage,
			Signature:       &userInfo.Signature,
			TotalFavorited:  &userCountInfo.TotalFavorited,
			WorkCount:       &userCountInfo.WorkCount,
			FavoriteCount:   &userCountInfo.FavoriteCount,
		},
	}
	//TODO 使用消息队列 将 userInfo 存储 redis
	if err = s.rCache.HSetUserInfo(ctx, auto.CreateUserKey(uint(req.UserId)), auto.CreateMapUserInfo(userInfo)); err != nil {
		zap.L().Error("s.rCache.HSetUserInfo err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrRedis, err.Error(), &user.UserIndexResponse{}), nil
	}
	if err = s.rCache.HSetUserCountInfo(ctx, auto.CreateUserCountKey(uint(req.UserId)), auto.CreateMapUserCount(userCountInfo)); err != nil {
		zap.L().Error("s.rCache.HSetUserCountInfo err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrRedis, err.Error(), &user.UserIndexResponse{}), nil
	}
	return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) AddFollowCount(ctx context.Context, req *user.AddFollowCountRequest) (resp *user.AddFollowCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	//// 执行事务，对当前userId_follow_count+1，对targetId_follower_count+1
	//if err = s.transaction.Action(func(conn mysql.DbConn) error {
	//	if err = s.UserRepo.AddFollowCount(ctx, conn, uint(req.UserId)); err != nil {
	//		zap.L().Error("s.UserRepo.AddFollowCount err:", zap.Error(err))
	//		return err
	//	}
	//	if err = s.UserRepo.AddFollowerCount(ctx, conn, uint(req.TargetId)); err != nil {
	//		zap.L().Error("s.UserRepo.AddFollowerCount err:", zap.Error(err))
	//		return err
	//	}
	//	return nil
	//}); err != nil {
	//	return s.respRepo.AddFollowCountResponse(errcode.ErrDB, err.Error(), &user.AddFollowCountResponse{}), nil
	//}
	// 对 redis中 userInfo中的follow_count做更改
	userCountKey := auto.CreateUserCountKey(uint(req.UserId))
	targetCountKey := auto.CreateUserCountKey(uint(req.TargetId))
	err = s.rCache.AddFollowCount(ctx, userCountKey)
	if err != nil {
		zap.L().Error("s.rCache.AddFollowCount err:", zap.Error(err))
		return s.respRepo.AddFollowCountResponse(errcode.ErrRedis, err.Error(), &user.AddFollowCountResponse{}), nil
	}
	err = s.rCache.AddFollowerCount(ctx, targetCountKey)
	if err != nil {
		zap.L().Error("s.rCache.AddFollowerCount err:", zap.Error(err))
		return s.respRepo.AddFollowCountResponse(errcode.ErrRedis, err.Error(), &user.AddFollowCountResponse{}), nil
	}
	return s.respRepo.AddFollowCountResponse(errcode.StatusOK, model.MsgNil, &user.AddFollowCountResponse{}), nil
}

func (s *UserServiceImpl) SubFollowCount(ctx context.Context, req *user.SubFollowCountRequest) (resp *user.SubFollowCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 执行事务，对当前userId_follow_count-1，对targetId_follower_count-1
	//if err = s.transaction.Action(func(conn mysql.DbConn) error {
	//	if err = s.UserRepo.SubFollowCount(ctx, conn, uint(req.UserId)); err != nil {
	//		zap.L().Error("s.UserRepo.SubFollowCount err:", zap.Error(err))
	//		return err
	//	}
	//	if err = s.UserRepo.SubFollowerCount(ctx, conn, uint(req.TargetId)); err != nil {
	//		zap.L().Error("s.UserRepo.SubFollowerCount err:", zap.Error(err))
	//		return err
	//	}
	//	return nil
	//}); err != nil {
	//	return s.respRepo.SubFollowCountResponse(errcode.ErrDB, err.Error(), &user.SubFollowCountResponse{}), nil
	//}
	// 对 redis中 userInfo中的follow_count做更改
	userCountKey := auto.CreateUserCountKey(uint(req.UserId))
	targetCountKey := auto.CreateUserCountKey(uint(req.TargetId))
	err = s.rCache.SubFollowCount(ctx, userCountKey)
	if err != nil {
		zap.L().Error("s.rCache.SubFollowCount err:", zap.Error(err))
		return s.respRepo.SubFollowCountResponse(errcode.ErrRedis, err.Error(), &user.SubFollowCountResponse{}), nil
	}
	err = s.rCache.SubFollowerCount(ctx, targetCountKey)
	if err != nil {
		zap.L().Error("s.rCache.SubFollowCount err:", zap.Error(err))
		return s.respRepo.SubFollowCountResponse(errcode.ErrRedis, err.Error(), &user.SubFollowCountResponse{}), nil
	}
	return s.respRepo.SubFollowCountResponse(errcode.StatusOK, model.MsgNil, &user.SubFollowCountResponse{}), nil
}

func (s *UserServiceImpl) GetUserList(ctx context.Context, req *user.GetUserListRequest) (resp *user.GetUserListResponse, err error) {
	resp = new(user.GetUserListResponse)
	resp.UserList = make([]*user.User, len(req.TargetId))
	var userKey string
	var userCntKey string
	// 1.校验业务逻辑
	// 一次性查所有的 is_follow
	isFollowResp, _ := client.InteractionClient.IsFollowManyTargets(ctx, &interaction.IsFollowManyTargetsRequest{
		UserId:    req.UserId,
		TargetIds: req.TargetId,
	})
	if isFollowResp == nil {
		zap.L().Error("client.InteractionClient.IsFollowManyTargets 返回空指针")
		return s.respRepo.GetUserListResponse(errcode.ErrServer, myerr.IsFollowManyTargetsErr.Error(), &user.GetUserListResponse{}), nil
	}
	if isFollowResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.InteractionClient.IsFollowManyTargets err:", zap.Error(err))
		return s.respRepo.GetUserListResponse(errcode.CreateErr(isFollowResp.StatusCode, myerr.IsFollowManyTargetsErr.Error()), isFollowResp.StatusMsg, &user.GetUserListResponse{}), nil
	}
	// 2.处理业务
	// 先查缓存
	for i, targetId := range req.TargetId {
		userKey = auto.CreateUserKey(uint(targetId))
		userCntKey = auto.CreateUserCountKey(uint(targetId))
		userInfo, err := s.rCache.HGetUserInfo(ctx, userKey)
		if err != nil {
			zap.L().Error("s.rCache.HGetUserInfo err:", zap.Error(err))
			return s.respRepo.GetUserListResponse(errcode.ErrRedis, err.Error(), &user.GetUserListResponse{}), nil
		}
		// 缓存没查到，查数据库
		if userInfo == nil {
			userInfo, err = s.UserRepo.GetUserInfoByUserID(ctx, targetId)
			if err != nil {
				zap.L().Error("s.UserRepo.GetUserInfoByUserID err:", zap.Error(err))
				return s.respRepo.GetUserListResponse(errcode.ErrDB, err.Error(), &user.GetUserListResponse{}), nil
			}
			// 存缓存
			err = s.rCache.HSetUserInfo(ctx, userKey, auto.CreateMapUserInfo(userInfo))
			if err != nil {
				zap.L().Error("s.rCache.HSetUserInfo err:", zap.Error(err))
				return s.respRepo.GetUserListResponse(errcode.ErrRedis, err.Error(), &user.GetUserListResponse{}), nil
			}
		}
		userCntInfo, err := s.rCache.HGetUserCountInfo(ctx, userCntKey)
		if err != nil {
			zap.L().Error("s.rCache.HGetUserCountInfo err:", zap.Error(err))
			return s.respRepo.GetUserListResponse(errcode.ErrRedis, err.Error(), &user.GetUserListResponse{}), nil
		}
		// 缓存没查到，查数据库
		if userCntInfo == nil {
			userCntInfo, err = s.UserRepo.GetUserCountInfoByUserID(ctx, targetId)
			if err != nil {
				zap.L().Error("s.UserRepo.GetUserInfoByUserID err:", zap.Error(err))
				return s.respRepo.GetUserListResponse(errcode.ErrDB, err.Error(), &user.GetUserListResponse{}), nil
			}
			// 存缓存
			err = s.rCache.HSetUserCountInfo(ctx, userCntKey, auto.CreateMapUserCount(userCntInfo))
			if err != nil {
				zap.L().Error("s.rCache.HSetUserCountInfo err:", zap.Error(err))
				return s.respRepo.GetUserListResponse(errcode.ErrRedis, err.Error(), &user.GetUserListResponse{}), nil
			}
		}
		// 模型转换
		resp.UserList[i] = &user.User{
			Id:              targetId,
			Name:            userInfo.Username,
			FollowCount:     &userCntInfo.FollowCount,
			FollowerCount:   &userCntInfo.FollowerCount,
			IsFollow:        isFollowResp.ManyExist[i],
			Avatar:          &userInfo.Avatar,
			BackgroundImage: &userInfo.BackgroundImage,
			Signature:       &userInfo.Signature,
			TotalFavorited:  &userCntInfo.TotalFavorited,
			WorkCount:       &userCntInfo.WorkCount,
			FavoriteCount:   &userCntInfo.FavoriteCount,
		}
	}
	return s.respRepo.GetUserListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) AddUserWorkCount(ctx context.Context, req *user.AddUserWorkCountRequest) (resp *user.AddUserWorkCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 将redis中信息 更改
	userKey := auto.CreateUserCountKey(uint(req.UserId))
	err = s.rCache.AddWorkCount(ctx, userKey)
	if err != nil {
		zap.L().Error("s.rCache.AddFollowCount err:", zap.Error(err))
		return s.respRepo.AddUserWorkCountResponse(errcode.ErrRedis, err.Error(), &user.AddUserWorkCountResponse{}), nil
	}
	//TODO 数据库更改交给定时任务
	return s.respRepo.AddUserWorkCountResponse(errcode.StatusOK, model.MsgNil, &user.AddUserWorkCountResponse{}), nil
}

func (s *UserServiceImpl) UpdateUserFavoriteCount(ctx context.Context, req *user.UpdateUserFavoriteCountRequest) (resp *user.UpdateUserFavoriteCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 将redis中信息 更改
	userCntKey := auto.CreateUserCountKey(uint(req.UserId))
	authorCntKey := auto.CreateUserCountKey(uint(req.AuthorId))
	switch req.ActionType {
	case model.FAVORITE:
		err = s.rCache.AddFavoriteCount(ctx, userCntKey)
		if err != nil {
			zap.L().Error("s.rCache.AddFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
		err = s.rCache.AddTotalFavoriteCount(ctx, authorCntKey)
		if err != nil {
			zap.L().Error("s.rCache.AddTotalFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
	case model.CANCELFAVORITE:
		err = s.rCache.SubFavoriteCount(ctx, userCntKey)
		if err != nil {
			zap.L().Error("s.rCache.SubFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
		err = s.rCache.SubTotalFavoriteCount(ctx, authorCntKey)
		if err != nil {
			zap.L().Error("s.rCache.SubTotalFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
	}
	return s.respRepo.UpdateUserFavoriteCountResponse(errcode.StatusOK, model.MsgNil, &user.UpdateUserFavoriteCountResponse{}), nil
}
