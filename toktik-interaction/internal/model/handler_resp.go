package model

import (
	"github.com/Happy-Why/toktik-common/errcode"
	inter "github.com/Happy-Why/toktik-rpc/kitex_gen/interaction"
)

type FollowSBHandler struct {
}

func (FollowSBHandler) FollowSBResponse(err errcode.Err, msg string, resp *inter.FollowActionResponse) *inter.FollowActionResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type FollowListHandler struct {
}

func (FollowListHandler) FollowListResponse(err errcode.Err, msg string, resp *inter.FollowListResponse) *inter.FollowListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type FansListHandler struct {
}

func (FansListHandler) FansListResponse(err errcode.Err, msg string, resp *inter.FansListResponse) *inter.FansListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type FriendListHandler struct {
}

func (FriendListHandler) FriendListResponse(err errcode.Err, msg string, resp *inter.FriendListResponse) *inter.FriendListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
