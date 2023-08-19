package model

import (
	"github.com/Happy-Why/toktik-common/errcode"

	"github.com/Happy-Why/toktik-rpc/kitex_gen/video"
)

type VideoFeedHandler struct {
}

func (VideoFeedHandler) VideoFeedResponse(err errcode.Err, msg string, resp *video.VideoFeedResponse) *video.VideoFeedResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type VideoPublishHandler struct {
}

func (VideoPublishHandler) VideoPublishResponse(err errcode.Err, msg string, resp *video.VideoPublishResponse) *video.VideoPublishResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type PublishListHandler struct {
}

func (PublishListHandler) PublishListResponse(err errcode.Err, msg string, resp *video.PublishListResponse) *video.PublishListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type FavoriteListHandler struct {
}

func (FavoriteListHandler) FavoriteListResponse(err errcode.Err, msg string, resp *video.FavoriteListResponse) *video.FavoriteListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type FavoriteActionHandler struct {
}

func (FavoriteActionHandler) FavoriteActionResponse(err errcode.Err, msg string, resp *video.FavoriteActionResponse) *video.FavoriteActionResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type CommentActionHandler struct {
}

func (CommentActionHandler) CommentActionResponse(err errcode.Err, msg string, resp *video.CommentActionResponse) *video.CommentActionResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type CommentListHandler struct {
}

func (CommentListHandler) CommentListResponse(err errcode.Err, msg string, resp *video.CommentListResponse) *video.CommentListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
