package repo

import (
	"context"
	"time"
	"toktik-video/internal/model/auto"
)

type Cache interface {
	Put(c context.Context, key, value string, expire time.Duration) error
	Get(c context.Context, key string) (string, error)
}

type RClientRepo interface {
	GetKeys(c context.Context, keyPattern string) ([]string, error)
	KeyExist(c context.Context, key string) (bool, error)
	HSetVideoInfo(c context.Context, key string, value map[string]interface{}) error
	AddVideoFavoriteCount(c context.Context, key string) error
	SubVideoFavoriteCount(c context.Context, key string) error
	CreateFavorite(c context.Context, key string, video int64) error
	DelFavorite(c context.Context, key string, videoId int64) error
	IsFavRecordExist(c context.Context, key string, videoId int64) (bool, error)
	HGetVideoInfo(c context.Context, key string) (*auto.Video, error)
	GetFavoriteVideoIds(c context.Context, key string) ([]int64, error)
	PublishVideo(c context.Context, key string, time float64, videoId int64) error
	ZGetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error)
	AddVideoCommentCount(c context.Context, key string) error
	SubVideoCommentCount(c context.Context, key string) error
	SDelUserVideo(c context.Context, key string) error
	SGetUserVideoIds(c context.Context, key string) ([]int64, error)
	SAddUserVideoIds(c context.Context, key string, videoId []int64) error
	ZSAddCommentInfo(c context.Context, key string, time float64, content string) error
	DelComment(c context.Context, key string) error
	ZGetCommentList(c context.Context, key string) ([]int64, []*auto.Comment, error)
}
