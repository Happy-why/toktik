package repo

import (
	"toktik-common/errcode"
	"toktik-interaction/internal/model"
	inter "toktik-rpc/kitex_gen/interaction"
)

type HandlerResp interface {
	FollowSBResponse(err errcode.Err, msg string, resp *inter.FollowActionResponse) *inter.FollowActionResponse
	FollowListResponse(err errcode.Err, msg string, resp *inter.FollowListResponse) *inter.FollowListResponse
	FansListResponse(err errcode.Err, msg string, resp *inter.FansListResponse) *inter.FansListResponse
	FriendListResponse(err errcode.Err, msg string, resp *inter.FriendListResponse) *inter.FriendListResponse
	IsFollowTargetResponse(err errcode.Err, msg string, resp *inter.IsFollowTargetResponse) *inter.IsFollowTargetResponse
	IsFollowManyTargetsResponse(err errcode.Err, msg string, resp *inter.IsFollowManyTargetsResponse) *inter.IsFollowManyTargetsResponse
	IsFriendResponse(err errcode.Err, msg string, resp *inter.IsFriendResponse) *inter.IsFriendResponse
}

type HandlerResps struct {
	model.FollowSBHandler
	model.FollowListHandler
	model.FansListHandler
	model.FriendListHandler
	model.IsFollowTargetHandler
	model.IsFollowManyTargetsHandler
	model.IsFriendHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
