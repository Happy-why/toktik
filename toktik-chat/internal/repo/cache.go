package repo

import (
	"context"
	"time"
	"toktik-chat/internal/model/auto"
)

type Cache interface {
	Put(c context.Context, key, value string, expire time.Duration) error
	Get(c context.Context, key string) (string, error)
}

type RClientRepo interface {
	PushHistoryMessage(c context.Context, key string, time float64, content string) error
	PushDBMessage(c context.Context, key string, time float64, content string) error
	ZRangeMessageList(c context.Context, key string) ([]*auto.Message, error)
	PushManyHistoryMessage(c context.Context, key string, messageInfos []*auto.Message) error
}
