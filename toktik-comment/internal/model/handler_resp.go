package model

import (
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/comment"
)

type CommentActionHandler struct {
}

func (CommentActionHandler) CommentActionResponse(err errcode.Err, msg string, resp *comment.CommentActionResponse) *comment.CommentActionResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type CommentListHandler struct {
}

func (CommentListHandler) CommentListResponse(err errcode.Err, msg string, resp *comment.CommentListResponse) *comment.CommentListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
