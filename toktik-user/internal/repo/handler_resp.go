package repo

import (
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/Happy-Why/toktik-user/internal/model"
)

type HandlerResp interface {
	RegisterResponse(err errcode.Err, msg string) *user.RegisterResponse
}

type HandlerResps struct {
	model.RegisterHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
