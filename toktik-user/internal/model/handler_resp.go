package model

import (
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/user"
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

type AddFollowCountHandler struct {
}

func (AddFollowCountHandler) AddFollowCountResponse(err errcode.Err, msg string, resp *user.AddFollowCountResponse) *user.AddFollowCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type SubFollowCountHandler struct {
}

func (SubFollowCountHandler) SubFollowCountResponse(err errcode.Err, msg string, resp *user.SubFollowCountResponse) *user.SubFollowCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type GutUserListHandler struct {
}

func (GutUserListHandler) GetUserListResponse(err errcode.Err, msg string, resp *user.GetUserListResponse) *user.GetUserListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type AddUserWorkCountHandler struct {
}

func (AddUserWorkCountHandler) AddUserWorkCountResponse(err errcode.Err, msg string, resp *user.AddUserWorkCountResponse) *user.AddUserWorkCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type UpdateUserFavoriteCountHandler struct {
}

func (UpdateUserFavoriteCountHandler) UpdateUserFavoriteCountResponse(err errcode.Err, msg string, resp *user.UpdateUserFavoriteCountResponse) *user.UpdateUserFavoriteCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
