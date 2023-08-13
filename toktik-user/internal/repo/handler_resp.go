package repo

import (
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/Happy-Why/toktik-user/internal/model"
)

type HandlerResp interface {
	RegisterResponse(err errcode.Err, msg string, resp *user.RegisterResponse) *user.RegisterResponse
	LoginResponse(err errcode.Err, msg string, resp *user.LoginResponse) *user.LoginResponse
	UserIndexResponse(err errcode.Err, msg string, resp *user.UserIndexResponse) *user.UserIndexResponse
}

type HandlerResps struct {
	model.RegisterHandler
	model.LoginHandler
	model.UserIndexHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
