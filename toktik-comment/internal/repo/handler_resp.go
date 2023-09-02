package repo

import (
	"toktik-comment/internal/model"
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/comment"
)

type HandlerResp interface {
	CommentActionResponse(err errcode.Err, msg string, resp *comment.CommentActionResponse) *comment.CommentActionResponse
	CommentListResponse(err errcode.Err, msg string, resp *comment.CommentListResponse) *comment.CommentListResponse
}

type HandlerResps struct {
	model.CommentActionHandler
	model.CommentListHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
