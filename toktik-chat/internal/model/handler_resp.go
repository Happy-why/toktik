package model

import (
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/chat"
)

type ChatActionHandler struct {
}

func (ChatActionHandler) ChatActionResponse(err errcode.Err, msg string, resp *chat.ChatActionResponse) *chat.ChatActionResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type MessageListHandler struct {
}

func (MessageListHandler) MessageListResponse(err errcode.Err, msg string, resp *chat.MessageListResponse) *chat.MessageListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type GetFriendLatestMessageHandler struct {
}

func (MessageListHandler) GetFriendLatestMessageResponse(err errcode.Err, msg string, resp *chat.GetFriendLatestMessageResponse) *chat.GetFriendLatestMessageResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
