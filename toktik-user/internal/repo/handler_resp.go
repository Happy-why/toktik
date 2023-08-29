package repo

import (
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/user"
	"toktik-user/internal/model"
)

type HandlerResp interface {
	RegisterResponse(err errcode.Err, msg string, resp *user.RegisterResponse) *user.RegisterResponse
	LoginResponse(err errcode.Err, msg string, resp *user.LoginResponse) *user.LoginResponse
	UserIndexResponse(err errcode.Err, msg string, resp *user.UserIndexResponse) *user.UserIndexResponse
	AddFollowCountResponse(err errcode.Err, msg string, resp *user.AddFollowCountResponse) *user.AddFollowCountResponse
	SubFollowCountResponse(err errcode.Err, msg string, resp *user.SubFollowCountResponse) *user.SubFollowCountResponse
	GetUserListResponse(err errcode.Err, msg string, resp *user.GetUserListResponse) *user.GetUserListResponse
	AddUserWorkCountResponse(err errcode.Err, msg string, resp *user.AddUserWorkCountResponse) *user.AddUserWorkCountResponse
	UpdateUserFavoriteCountResponse(err errcode.Err, msg string, resp *user.UpdateUserFavoriteCountResponse) *user.UpdateUserFavoriteCountResponse
}

type HandlerResps struct {
	model.RegisterHandler
	model.LoginHandler
	model.UserIndexHandler
	model.AddFollowCountHandler
	model.SubFollowCountHandler
	model.GutUserListHandler
	model.AddUserWorkCountHandler
	model.UpdateUserFavoriteCountHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
