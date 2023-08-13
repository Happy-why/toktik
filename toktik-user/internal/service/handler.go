package service

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-common/utils"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/Happy-Why/toktik-user/internal/dao"
	"github.com/Happy-Why/toktik-user/internal/dao/mysql"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/Happy-Why/toktik-user/internal/model"
	"github.com/Happy-Why/toktik-user/internal/model/auto"
	"github.com/Happy-Why/toktik-user/internal/repo"
	"github.com/Happy-Why/toktik-user/pkg/myerr"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	respRepo repo.HandlerResp
	cache    repo.Cache
	userRepo repo.UserRepo
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{
		cache:    dao.Group.Rdb,
		respRepo: repo.NewHandlerResps(),
		userRepo: mysql.NewUserDao(),
	}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(c context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1.可以校验参数
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
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
		Model:          gorm.Model{ID: uint(userID)},
		Username:       req.Username,
		Password:       hashPassword,
		Avatar:         global.PbSettings.Rules.DefaultUserAvatar,
		IsFollow:       false,
		FollowCount:    0,
		FollowerCount:  0,
		Signature:      global.PbSettings.Rules.DefaultUserSignature,
		TotalFavorited: 0,
		WorkCount:      0,
	}
	err = s.userRepo.UserRegister(ctx, userInfo)
	if err != nil {
		zap.L().Error("s.userRepo.UserRegister err:", zap.Error(err))
		return s.respRepo.RegisterResponse(errcode.ErrDB, err.Error(), &user.RegisterResponse{}), nil
	}
	// TODO 注册之后要不要自己关注自己呢？
	// 生成token TODO 后续看情况考虑 可以将 token 过期时间加到 Payload.content 中。
	token, content, err := CreateToken(userID)
	if err != nil {
		zap.L().Error("CreateToken err:", zap.Error(err))
		return s.respRepo.RegisterResponse(myerr.CreateTokenErr, err.Error(), &user.RegisterResponse{}), nil
	}
	// 将 token 中的 content 加入缓存，以便身份校验时的判断身份
	sss := make(chan struct{}, 0)
	go func() {
		if err := s.cache.Put(c, model.TypeUser+strconv.Itoa(int(userInfo.ID)), content, global.PvSettings.Token.UserTokenExp); err != nil {
			zap.L().Error("s.cache.Put err:", zap.Error(err))
			return
		}
		sss <- struct{}{}
		// TODO 使用 go程 的错误处理以及 response问题
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

	// 生成token TODO 后续看情况考虑 可以将 token 过期时间加到 Payload.content 中。
	token, content, err := CreateToken(int64(userInfo.ID))
	if err != nil {
		zap.L().Error("CreateToken err:", zap.Error(err))
		return s.respRepo.LoginResponse(myerr.CreateTokenErr, err.Error(), &user.LoginResponse{}), nil
	}

	// TODO go程测试 将 token 中的 content 加入缓存，以便身份校验时的判断身份
	sss := make(chan struct{}, 0)
	go func() {
		if err := s.cache.Put(c, model.TypeUser+strconv.Itoa(int(userInfo.ID)), content, global.PvSettings.Token.UserTokenExp); err != nil {
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
	c, cancel := context.WithTimeout(ctx, 200*time.Second)
	defer cancel()
	// 从缓存中查询 token 中的 content，校验身份
	content, err := s.cache.Get(c, model.TypeUser+strconv.Itoa(int(req.UserId)))
	fmt.Println("content:", content)
	if err != nil {
		zap.L().Error("TokenVerify cache get user error", zap.Error(err))
		return nil, err
	}
	return nil, nil
}

func (s *UserServiceImpl) UserIndex(ctx context.Context, req *user.UserIndexRequest) (resp *user.UserIndexResponse, err error) {
	// 校验逻辑
	// 进行业务
	// 根据 user_id 查询 userInfo
	userInfo, err := s.userRepo.GetUserInfoByUserID(ctx, req.UserId)
	if err == gorm.ErrRecordNotFound {
		return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, &user.UserIndexResponse{}), nil
	}
	if err != nil {
		zap.L().Error("s.userRepo.GetUserInfoByUserID err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrDB, err.Error(), &user.UserIndexResponse{}), nil
	}
	// 通过 查询 relation表 判断是否关注该用户 来决定 is_follow字段 , 获得的 userInfo.ID 是 target_id
	ok, err := s.userRepo.IsFollowUser(ctx, req.MyUserId, req.UserId)
	if err != nil {
		zap.L().Error("s.userRepo.IsFollowUser err:", zap.Error(err))
		return s.respRepo.UserIndexResponse(errcode.ErrDB, err.Error(), &user.UserIndexResponse{}), nil
	}
	if ok {
		userInfo.IsFollow = true
	}
	resp = &user.UserIndexResponse{
		User: &user.User{
			Id:              int64(userInfo.ID),
			Name:            userInfo.Username,
			FollowCount:     &userInfo.FollowCount,
			FollowerCount:   &userInfo.FollowerCount,
			IsFollow:        userInfo.IsFollow,
			Avatar:          &userInfo.Avatar,
			BackgroundImage: &userInfo.BackgroundImage,
			Signature:       &userInfo.Signature,
			TotalFavorited:  &userInfo.TotalFavorited,
			WorkCount:       &userInfo.WorkCount,
			FavoriteCount:   &userInfo.TotalFavorited,
		},
	}
	return s.respRepo.UserIndexResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
