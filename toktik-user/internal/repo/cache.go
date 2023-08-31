package repo

import (
	"context"
	"time"
	"toktik-user/internal/model/auto"
)

type Cache interface {
	Put(c context.Context, key, value string, expire time.Duration) error
	Get(c context.Context, key string) (string, error)
}

type RClientRepo interface {
	HSetUserInfo(c context.Context, key string, value map[string]interface{}) error
	HSetUserCountInfo(c context.Context, key string, value map[string]interface{}) error
	HGetUserInfo(c context.Context, key string) (*auto.User, error)
	HGetUserCountInfo(c context.Context, key string) (*auto.UserCount, error)
	AddFollowCount(c context.Context, key string) error
	AddFollowerCount(c context.Context, key string) error
	SubFollowCount(c context.Context, key string) error
	SubFollowerCount(c context.Context, key string) error
	AddWorkCount(c context.Context, key string) error
	AddFavoriteCount(c context.Context, key string) error
	AddTotalFavoriteCount(c context.Context, key string) error
	SubFavoriteCount(c context.Context, key string) error
	SubTotalFavoriteCount(c context.Context, key string) error
}
