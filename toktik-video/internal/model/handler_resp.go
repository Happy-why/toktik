package model

import (
	"toktik-common/errcode"

	"toktik-rpc/kitex_gen/video"
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

type GetVideoInfoHandler struct {
}

func (GetVideoInfoHandler) GetVideoInfoResponse(err errcode.Err, msg string, resp *video.GetVideoInfoResponse) *video.GetVideoInfoResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type GetManyVideoInfosHandler struct {
}

func (GetManyVideoInfosHandler) GetManyVideoInfosResponse(err errcode.Err, msg string, resp *video.GetManyVideoInfosResponse) *video.GetManyVideoInfosResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type AddVideoFavoriteCountHandler struct {
}

func (AddVideoFavoriteCountHandler) AddVideoFavoriteCountResponse(err errcode.Err, msg string, resp *video.AddVideoFavoriteCountResponse) *video.AddVideoFavoriteCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type SubVideoFavoriteCountHandler struct {
}

func (SubVideoFavoriteCountHandler) SubVideoFavoriteCountResponse(err errcode.Err, msg string, resp *video.SubVideoFavoriteCountResponse) *video.SubVideoFavoriteCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type AddVideoCommentCountHandler struct {
}

func (AddVideoCommentCountHandler) AddVideoCommentCountResponse(err errcode.Err, msg string, resp *video.AddVideoCommentCountResponse) *video.AddVideoCommentCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type SubVideoCommentCountHandler struct {
}

func (SubVideoCommentCountHandler) SubVideoCommentCountResponse(err errcode.Err, msg string, resp *video.SubVideoCommentCountResponse) *video.SubVideoCommentCountResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
