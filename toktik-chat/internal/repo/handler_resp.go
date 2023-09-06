package repo

import (
	"toktik-chat/internal/model"
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/chat"
)

type HandlerResp interface {
	ChatActionResponse(err errcode.Err, msg string, resp *chat.ChatActionResponse) *chat.ChatActionResponse
	MessageListResponse(err errcode.Err, msg string, resp *chat.MessageListResponse) *chat.MessageListResponse
	GetFriendLatestMessageResponse(err errcode.Err, msg string, resp *chat.GetFriendLatestMessageResponse) *chat.GetFriendLatestMessageResponse
}

type HandlerResps struct {
	model.ChatActionHandler
	model.MessageListHandler
	model.GetFriendLatestMessageHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
