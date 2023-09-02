package repo

import (
	"toktik-common/errcode"
	"toktik-favor/internal/model"
	"toktik-rpc/kitex_gen/favor"
)

type HandlerResp interface {
	FavoriteListResponse(err errcode.Err, msg string, resp *favor.FavoriteListResponse) *favor.FavoriteListResponse
	FavoriteActionResponse(err errcode.Err, msg string, resp *favor.FavoriteActionResponse) *favor.FavoriteActionResponse
	IsFavoriteVideoResponse(err errcode.Err, msg string, resp *favor.IsFavoriteVideosResponse) *favor.IsFavoriteVideosResponse
}

type HandlerResps struct {
	model.FavoriteListHandler
	model.FavoriteActionHandler
	model.IsFavoriteVideoHandler
}

func NewHandlerResps() *HandlerResps {
	return &HandlerResps{}
}
