package repo

import (
	"context"
)

type RCacheRepo interface {
	GetKeys(c context.Context, keyPattern string) ([]string, error)
	KeyExist(c context.Context, key string) (bool, error)
	CreateFavorite(c context.Context, key string, video int64) error
	DelFavorite(c context.Context, key string, videoId int64) error
	IsFavRecordExist(c context.Context, key string, videoId int64) (bool, error)
	GetFavoriteVideoIds(c context.Context, key string) ([]int64, error)
}
