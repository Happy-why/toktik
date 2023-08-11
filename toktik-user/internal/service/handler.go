package service

import (
	"context"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/Happy-Why/toktik-user/internal/dao"
	"github.com/Happy-Why/toktik-user/internal/dao/mysql"
	"github.com/Happy-Why/toktik-user/internal/model"
	"github.com/Happy-Why/toktik-user/internal/repo"
	"github.com/Happy-Why/toktik-user/pkg/myerr"
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
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// 1.可以校验参数
	// 2.校验业务逻辑(邮箱、账号、手机号是否被注册)
	exist, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return s.respRepo.RegisterResponse(errcode.ErrDB, err.Error()), nil
	}
	if exist {
		return s.respRepo.RegisterResponse(myerr.UsernameExist, model.MsgNil), nil
	}
	// 3.执行业务
	// 生成token

	// 存入用户名和密码

	// 4.返回
	return resp, nil
}
