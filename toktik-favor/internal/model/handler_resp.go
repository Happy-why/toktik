package model

import (
	"toktik-common/errcode"

	"toktik-rpc/kitex_gen/favor"
)

type FavoriteListHandler struct {
}

func (FavoriteListHandler) FavoriteListResponse(err errcode.Err, msg string, resp *favor.FavoriteListResponse) *favor.FavoriteListResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type FavoriteActionHandler struct {
}

func (FavoriteActionHandler) FavoriteActionResponse(err errcode.Err, msg string, resp *favor.FavoriteActionResponse) *favor.FavoriteActionResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}

type IsFavoriteVideoHandler struct {
}

func (IsFavoriteVideoHandler) IsFavoriteVideoResponse(err errcode.Err, msg string, resp *favor.IsFavoriteVideosResponse) *favor.IsFavoriteVideosResponse {
	resp.StatusCode = err.ECode()
	resp.StatusMsg = err.Error() + ":" + msg
	return resp
}
