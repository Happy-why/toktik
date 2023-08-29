package repo

import (
	"context"
	redis2 "github.com/go-redis/redis/v8"
	"time"
)

type Cache interface {
	Put(c context.Context, key, value string, expire time.Duration) error
	Get(c context.Context, key string) (string, error)
}

type RClientRepo interface {
	PushHistoryMessage(c context.Context, key string, time float64, content string) error
	PushDBMessage(c context.Context, key string, time float64, content string) error
	ZRangeMessageList(c context.Context, key string) ([]redis2.Z, error)
}
