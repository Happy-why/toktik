package repo

import (
	"context"
	"toktik-chat/internal/model/auto"
)

type RCacheRepo interface {
	KeyExist(c context.Context, key string) (bool, error)
	PushHistoryMessage(c context.Context, key string, time float64, content string) error
	PushDBMessage(c context.Context, key string, time float64, content string) error
	ZRangeMessageList(c context.Context, key string, preMsgTime int64) ([]*auto.Message, error)
	PushManyHistoryMessage(c context.Context, key string, messageInfos []*auto.Message) error
	ZGetFriendLatestMessage(c context.Context, key string, myUserId int64) (string, int32, error)
}
