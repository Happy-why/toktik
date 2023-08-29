package model

import (
	"toktik-common/errcode"
	inter "toktik-rpc/kitex_gen/interaction"
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

type IsFollowTargetHandler struct {
}

func (IsFollowTargetHandler) IsFollowTargetResponse(err errcode.Err, msg string, resp *inter.IsFollowTargetResponse) *inter.IsFollowTargetResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type IsFollowManyTargetsHandler struct {
}

func (IsFollowManyTargetsHandler) IsFollowManyTargetsResponse(err errcode.Err, msg string, resp *inter.IsFollowManyTargetsResponse) *inter.IsFollowManyTargetsResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type IsFriendHandler struct {
}

func (IsFriendHandler) IsFriendResponse(err errcode.Err, msg string, resp *inter.IsFriendResponse) *inter.IsFriendResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
