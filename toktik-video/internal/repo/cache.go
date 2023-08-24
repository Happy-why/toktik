package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(c context.Context, key, value string, expire time.Duration) error
	Get(c context.Context, key string) (string, error)
}

type RClientRepo interface {
	HSetVideoInfo(c context.Context, key string, value map[string]interface{}) error
}
