package model

import (
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
)

type RegisterHandler struct {
}

func (RegisterHandler) RegisterResponse(err errcode.Err, msg string, resp *user.RegisterResponse) *user.RegisterResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type LoginHandler struct {
}

func (LoginHandler) LoginResponse(err errcode.Err, msg string, resp *user.LoginResponse) *user.LoginResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type UserIndexHandler struct {
}

func (UserIndexHandler) UserIndexResponse(err errcode.Err, msg string, resp *user.UserIndexResponse) *user.UserIndexResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
