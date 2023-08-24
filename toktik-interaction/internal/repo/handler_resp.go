package repo

import (
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-interaction/internal/model"
	inter "github.com/Happy-Why/toktik-rpc/kitex_gen/interaction"
)

type HandlerResp interface {
	FollowSBResponse(err errcode.Err, msg string, resp *inter.FollowActionResponse) *inter.FollowActionResponse
	FollowListResponse(err errcode.Err, msg string, resp *inter.FollowListResponse) *inter.FollowListResponse
	FansListResponse(err errcode.Err, msg string, resp *inter.FansListResponse) *inter.FansListResponse
	FriendListResponse(err errcode.Err, msg string, resp *inter.FriendListResponse) *inter.FriendListResponse
	IsFollowTargetResponse(err errcode.Err, msg string, resp *inter.IsFollowTargetResponse) *inter.IsFollowTargetResponse
}

type HandlerResps struct {
	model.FollowSBHandler
	model.FollowListHandler
	model.FansListHandler
	model.FriendListHandler
	model.IsFollowTargetHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
