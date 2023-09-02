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
	GetVideoInfoResponse(err errcode.Err, msg string, resp *video.GetVideoInfoResponse) *video.GetVideoInfoResponse
	GetManyVideoInfosResponse(err errcode.Err, msg string, resp *video.GetManyVideoInfosResponse) *video.GetManyVideoInfosResponse
	AddVideoFavoriteCountResponse(err errcode.Err, msg string, resp *video.AddVideoFavoriteCountResponse) *video.AddVideoFavoriteCountResponse
	SubVideoFavoriteCountResponse(err errcode.Err, msg string, resp *video.SubVideoFavoriteCountResponse) *video.SubVideoFavoriteCountResponse
	AddVideoCommentCountResponse(err errcode.Err, msg string, resp *video.AddVideoCommentCountResponse) *video.AddVideoCommentCountResponse
	SubVideoCommentCountResponse(err errcode.Err, msg string, resp *video.SubVideoCommentCountResponse) *video.SubVideoCommentCountResponse
}

type HandlerResps struct {
	model.VideoFeedHandler
	model.VideoPublishHandler
	model.PublishListHandler
	model.GetVideoInfoHandler
	model.GetManyVideoInfosHandler
	model.AddVideoFavoriteCountHandler
	model.SubVideoFavoriteCountHandler
	model.AddVideoCommentCountHandler
	model.SubVideoCommentCountHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
