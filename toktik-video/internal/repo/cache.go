package repo

import (
	"context"
	"toktik-video/internal/model/auto"
)

type RCacheRepo interface {
	GetKeys(c context.Context, keyPattern string) ([]string, error)
	KeyExist(c context.Context, key string) (bool, error)
	HSetVideoInfo(c context.Context, key string, value map[string]interface{}) error
	AddVideoFavoriteCount(c context.Context, key string) error
	SubVideoFavoriteCount(c context.Context, key string) error
	HGetVideoInfo(c context.Context, key string) (*auto.Video, error)
	PublishVideo(c context.Context, key string, time float64, videoId int64) error
	ZGetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error)
	AddVideoCommentCount(c context.Context, key string) error
	SubVideoCommentCount(c context.Context, key string) error
	SDelUserVideo(c context.Context, key string) error
	SGetUserVideoIds(c context.Context, key string) ([]int64, error)
	SAddUserVideoIds(c context.Context, key string, videoId []int64) error
}
