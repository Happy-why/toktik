package service

import (
	"context"
	"errors"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
	"toktik-common/errcode"
	"toktik-common/utils"
	"toktik-rpc/kitex_gen/interaction"
	"toktik-rpc/kitex_gen/user"
	"toktik-user/internal/dao/mysql"
	"toktik-user/internal/dao/redis"
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
	cache       repo.Cache
	userRepo    repo.UserRepo
	rClient     repo.RClientRepo
	transaction mysql.Transaction
	//rpc repo.RpcRepo
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{
		cache:       redis.GetRdbCache(),
		respRepo:    repo.NewHandlerResps(),
		userRepo:    mysql.NewUserDao(),
		rClient:     redis.NewUserRClient(),
		transaction: mysql.NewTransaction(),
		//rpc: client.NewRpcClientRepo(),
	}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1.可以校验参数
	// 2.校验业务逻辑(邮箱、账号、手机号是否被注册)
	exist, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		zap.L().Error("s.userRepo.GetUserByUsername err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrDB, err.Error(), &user.RegisterResponse{}), nil
	}
	if exist {
		return s.respRepo.RegisterResponse(myerr.UsernameExist, model.MsgNil, &user.RegisterResponse{}), nil
	}
	// 3.执行业务
	// 生成唯一id，存入用户名和密码
	userID, _ := global.SnowFlake.NextId()
	hashPassword, _ := utils.HashPassword(req.Password)
	userInfo := &auto.User{
		BaseModel:       auto.BaseModel{ID: uint(userID)},
		Username:        req.Username,
		Password:        hashPassword,
		Avatar:          global.Settings.Rules.DefaultUserAvatar,
		IsFollow:        false,
		FollowCount:     0,
		FollowerCount:   0,
		BackgroundImage: global.Settings.Rules.DefaultUserBackGroundImage,
		Signature:       global.Settings.Rules.DefaultUserSignature,
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
	// 将 userInfo 存储 数据库
	if err = s.userRepo.UserRegister(ctx, userInfo); err != nil {
		zap.L().Error("s.userRepo.UserRegister err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrDB, err.Error(), &user.RegisterResponse{}), nil
	}
	// 将 userInfo 存储 redis
	if err = s.rClient.HSetUserInfo(ctx, auto.CreateUserKey(uint(userID)), auto.CreateMapUserInfo(userInfo)); err != nil {
		zap.L().Error("s.rClient.HSetUserInfo err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrRedis, err.Error(), &user.RegisterResponse{}), nil
	}

	// 生成token
	token, content, err := CreateToken(userID)
	if err != nil {
		zap.L().Error("CreateToken err:", zap.Error(err))
		return s.respRepo.RegisterResponse(myerr.CreateTokenErr, err.Error(), &user.RegisterResponse{}), nil
	}
	// 将 token 中的 content 加入缓存，以便身份校验时的判断身份
	sss := make(chan struct{}, 0)
	go func() {
		if err := s.cache.Put(ctx, model.TypeUserToken+strconv.Itoa(int(userInfo.ID)), content, global.Settings.Token.UserTokenExp); err != nil {
			zap.L().Error("s.cache.Put err:", zap.Error(err))
			return
		}
		sss <- struct{}{}
		// TODO 使用 go程 的错误处理
	}()
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
	userInfo, err := s.userRepo.GetUserInfoByUsername(c, req.Username)
	if err != nil {
		zap.L().Error("s.userRepo.GetUserByUsername err:", zap.Error(err))
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
	sss := make(chan struct{}, 0)
	go func() {
		if err := s.cache.Put(c, model.TypeUserToken+strconv.Itoa(int(userInfo.ID)), content, global.Settings.Token.UserTokenExp); err != nil {
			zap.L().Error("s.cache.Put err:", zap.Error(err))
			return
		}
		sss <- struct{}{}
		// TODO 使用 go程 的错误处理以及 response问题
	}()
	<-sss
	// 4.返回
	resp = &user.LoginResponse{
		UserId: int64(userInfo.ID),
		Token:  token,
	}
	return s.respRepo.LoginResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) TokenVerify(ctx context.Context, req *user.TokenVerifyRequest) (resp *user.TokenVerifyResponse, err error) {
	// 从缓存中查询 token 中的 content，校验身份
	content, err := s.cache.Get(ctx, model.TypeUserToken+strconv.Itoa(int(req.UserId)))
	fmt.Println("content:", content)
	if err == redis2.Nil {
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

	// 2.查询缓存
	userKey := auto.CreateUserKey(uint(req.UserId))
	userInfo, err := s.rClient.HGetUserInfo(ctx, userKey)
	fmt.Println("userInfo2:", userInfo)
	if err != nil {
		zap.L().Error("s.rClient.HGetUserInfo err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrRedis, err.Error(), &user.UserIndexResponse{}), nil
	}
	if userInfo != nil {
		fmt.Println("走缓存！！！")
		resp = &user.UserIndexResponse{
			User: &user.User{
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
		}
		return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, resp), nil
	}

	// 3.进行业务
	// 根据 user_id 查询 userInfo
	userInfo, err = s.userRepo.GetUserInfoByUserID(ctx, req.UserId)
	if err == gorm.ErrRecordNotFound {
		return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, &user.UserIndexResponse{}), nil
	}
	if err != nil {
		zap.L().Error("s.userRepo.GetUserInfoByUserID err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrDB, err.Error(), &user.UserIndexResponse{}), nil
	}
	// 模型转换
	resp = &user.UserIndexResponse{
		User: &user.User{
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
	}
	// 将 userInfo 存储 redis
	if err = s.rClient.HSetUserInfo(ctx, auto.CreateUserKey(uint(req.UserId)), auto.CreateMapUserInfo(userInfo)); err != nil {
		zap.L().Error("s.rClient.HSetUserInfo err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrRedis, err.Error(), &user.UserIndexResponse{}), nil
	}

	return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) AddFollowCount(ctx context.Context, req *user.AddFollowCountRequest) (resp *user.AddFollowCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 执行事务，对当前userId_follow_count+1，对targetId_follower_count+1
	if err = s.transaction.Action(func(conn mysql.DbConn) error {
		if err = s.userRepo.AddFollowCount(ctx, conn, uint(req.UserId)); err != nil {
			zap.L().Error("s.userRepo.AddFollowCount err:", zap.Error(err))
			return err
		}
		if err = s.userRepo.AddFollowerCount(ctx, conn, uint(req.TargetId)); err != nil {
			zap.L().Error("s.userRepo.AddFollowerCount err:", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return s.respRepo.AddFollowCountResponse(errcode.ErrDB, err.Error(), &user.AddFollowCountResponse{}), nil
	}
	// 对 redis中 userInfo中的follow_count做更改
	userKey := auto.CreateUserKey(uint(req.UserId))
	targetKey := auto.CreateUserKey(uint(req.TargetId))
	err = s.rClient.AddFollowCount(ctx, userKey)
	if err != nil {
		zap.L().Error("s.rClient.AddFollowCount err:", zap.Error(err))
		return s.respRepo.AddFollowCountResponse(errcode.ErrRedis, err.Error(), &user.AddFollowCountResponse{}), nil
	}
	err = s.rClient.AddFollowerCount(ctx, targetKey)
	if err != nil {
		zap.L().Error("s.rClient.AddFollowerCount err:", zap.Error(err))
		return s.respRepo.AddFollowCountResponse(errcode.ErrRedis, err.Error(), &user.AddFollowCountResponse{}), nil
	}
	return s.respRepo.AddFollowCountResponse(errcode.StatusOK, model.MsgNil, &user.AddFollowCountResponse{}), nil
}

func (s *UserServiceImpl) SubFollowCount(ctx context.Context, req *user.SubFollowCountRequest) (resp *user.SubFollowCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 执行事务，对当前userId_follow_count-1，对targetId_follower_count-1
	if err = s.transaction.Action(func(conn mysql.DbConn) error {
		if err = s.userRepo.SubFollowCount(ctx, conn, uint(req.UserId)); err != nil {
			zap.L().Error("s.userRepo.SubFollowCount err:", zap.Error(err))
			return err
		}
		if err = s.userRepo.SubFollowerCount(ctx, conn, uint(req.TargetId)); err != nil {
			zap.L().Error("s.userRepo.SubFollowerCount err:", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return s.respRepo.SubFollowCountResponse(errcode.ErrDB, err.Error(), &user.SubFollowCountResponse{}), nil
	}
	// 对 redis中 userInfo中的follow_count做更改
	userKey := auto.CreateUserKey(uint(req.UserId))
	targetKey := auto.CreateUserKey(uint(req.TargetId))
	err = s.rClient.SubFollowCount(ctx, userKey)
	if err != nil {
		zap.L().Error("s.rClient.SubFollowCount err:", zap.Error(err))
		return s.respRepo.SubFollowCountResponse(errcode.ErrRedis, err.Error(), &user.SubFollowCountResponse{}), nil
	}
	err = s.rClient.SubFollowerCount(ctx, targetKey)
	if err != nil {
		zap.L().Error("s.rClient.SubFollowCount err:", zap.Error(err))
		return s.respRepo.SubFollowCountResponse(errcode.ErrRedis, err.Error(), &user.SubFollowCountResponse{}), nil
	}
	return s.respRepo.SubFollowCountResponse(errcode.StatusOK, model.MsgNil, &user.SubFollowCountResponse{}), nil
}

func (s *UserServiceImpl) GetUserList(ctx context.Context, req *user.GetUserListRequest) (resp *user.GetUserListResponse, err error) {
	resp = new(user.GetUserListResponse)
	// 1.校验业务逻辑
	// 2.处理业务
	userInfos, err := s.userRepo.GetUserList(ctx, req.TargetId)
	if err != nil {
		zap.L().Error("s.userRepo.GetUserList err:", zap.Error(err))
		return s.respRepo.GetUserListResponse(errcode.ErrDB, err.Error(), &user.GetUserListResponse{}), nil
	}
	isFollowResp, _ := client.InteractionClient.IsFollowManyTargets(ctx, &interaction.IsFollowManyTargetsRequest{
		UserId:    req.UserId,
		TargetIds: req.TargetId,
	})
	if isFollowResp == nil {
		zap.L().Error("client.InteractionClient.IsFollowManyTargets 返回空指针")
		return s.respRepo.GetUserListResponse(errcode.ErrServer, model.MsgNil, &user.GetUserListResponse{}), nil
	}
	if isFollowResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.InteractionClient.IsFollowManyTargets err:", zap.Error(err))
		return s.respRepo.GetUserListResponse(errcode.CreateErr(isFollowResp.StatusCode, model.MsgNil), isFollowResp.StatusMsg, &user.GetUserListResponse{}), nil
	}
	for i, v := range userInfos {
		resp.UserList = append(resp.UserList, &user.User{
			Id:              int64(v.ID),
			Name:            v.Username,
			FollowCount:     &v.FollowCount,
			FollowerCount:   &v.FollowerCount,
			IsFollow:        isFollowResp.ManyExist[i],
			Avatar:          &v.Avatar,
			BackgroundImage: &v.BackgroundImage,
			Signature:       &v.Signature,
			TotalFavorited:  &v.TotalFavorited,
			WorkCount:       &v.WorkCount,
			FavoriteCount:   &v.FavoriteCount,
		})
	}
	return s.respRepo.GetUserListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (s *UserServiceImpl) AddUserWorkCount(ctx context.Context, req *user.AddUserWorkCountRequest) (resp *user.AddUserWorkCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 将redis中信息 更改
	userKey := auto.CreateUserKey(uint(req.UserId))
	err = s.rClient.AddWorkCount(ctx, userKey)
	if err != nil {
		zap.L().Error("s.rClient.AddFollowCount err:", zap.Error(err))
		return s.respRepo.AddUserWorkCountResponse(errcode.ErrRedis, err.Error(), &user.AddUserWorkCountResponse{}), nil
	}
	//TODO 数据库更改交给定时任务
	return s.respRepo.AddUserWorkCountResponse(errcode.StatusOK, model.MsgNil, &user.AddUserWorkCountResponse{}), nil
}

func (s *UserServiceImpl) UpdateUserFavoriteCount(ctx context.Context, req *user.UpdateUserFavoriteCountRequest) (resp *user.UpdateUserFavoriteCountResponse, err error) {
	// 1.校验业务逻辑
	// 2.处理业务
	// 将redis中信息 更改
	userKey := auto.CreateUserKey(uint(req.UserId))
	authorKey := auto.CreateUserKey(uint(req.AuthorId))
	switch req.ActionType {
	case model.FAVORITE:
		err = s.rClient.AddFavoriteCount(ctx, userKey)
		if err != nil {
			zap.L().Error("s.rClient.AddFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
		err = s.rClient.AddTotalFavoriteCount(ctx, authorKey)
		if err != nil {
			zap.L().Error("s.rClient.AddTotalFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
	case model.CANCELFAVORITE:
		err = s.rClient.SubFavoriteCount(ctx, userKey)
		if err != nil {
			zap.L().Error("s.rClient.SubFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
		err = s.rClient.SubTotalFavoriteCount(ctx, authorKey)
		if err != nil {
			zap.L().Error("s.rClient.SubTotalFavoriteCount err:", zap.Error(err))
			return s.respRepo.UpdateUserFavoriteCountResponse(errcode.ErrRedis, err.Error(), &user.UpdateUserFavoriteCountResponse{}), nil
		}
	}
	return s.respRepo.UpdateUserFavoriteCountResponse(errcode.StatusOK, model.MsgNil, &user.UpdateUserFavoriteCountResponse{}), nil

}
