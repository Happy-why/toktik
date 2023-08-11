package model

import (
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
)

type RegisterHandler struct {
}

func (RegisterHandler) RegisterResponse(err errcode.Err, msg string) *user.RegisterResponse {
	u := &user.RegisterResponse{
		StatusCode: err.ECode(),
		StatusMsg:  err.Error() + ":" + msg,
	}
	return u
}
