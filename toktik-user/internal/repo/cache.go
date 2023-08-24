package repo

import (
	"context"
	"github.com/Happy-Why/toktik-user/internal/model/auto"
	"time"
)

type Cache interface {
	Put(c context.Context, key, value string, expire time.Duration) error
	Get(c context.Context, key string) (string, error)
}

type RClientRepo interface {
	HSetUserInfo(c context.Context, key string, value map[string]interface{}) error
	HGetUserInfo(c context.Context, key string) (*auto.User, error)
	AddFollowCount(c context.Context, key string) error
	AddFollowerCount(c context.Context, key string) error
	SubFollowCount(c context.Context, key string) error
	SubFollowerCount(c context.Context, key string) error
}
