package repo

import (
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/video"
	"toktik-video/internal/model"
)

type HandlerResp interface {
	VideoFeedResponse(err errcode.Err, msg string, resp *video.VideoFeedResponse) *video.VideoFeedResponse
	VideoPublishResponse(err errcode.Err, msg string, resp *video.VideoPublishResponse) *video.VideoPublishResponse
	PublishListResponse(err errcode.Err, msg string, resp *video.PublishListResponse) *video.PublishListResponse
	FavoriteListResponse(err errcode.Err, msg string, resp *video.FavoriteListResponse) *video.FavoriteListResponse
	FavoriteActionResponse(err errcode.Err, msg string, resp *video.FavoriteActionResponse) *video.FavoriteActionResponse
	CommentActionResponse(err errcode.Err, msg string, resp *video.CommentActionResponse) *video.CommentActionResponse
	CommentListResponse(err errcode.Err, msg string, resp *video.CommentListResponse) *video.CommentListResponse
}

type HandlerResps struct {
	model.VideoFeedHandler
	model.VideoPublishHandler
	model.PublishListHandler
	model.FavoriteListHandler
	model.FavoriteActionHandler
	model.CommentActionHandler
	model.CommentListHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
