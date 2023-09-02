package repo

import (
	"context"
	"toktik-comment/internal/model/auto"
)

type RCacheRepo interface {
	GetKeys(c context.Context, keyPattern string) ([]string, error)
	KeyExist(c context.Context, key string) (bool, error)
	ZSAddCommentInfo(c context.Context, key string, time float64, content string) error
	DelComment(c context.Context, key string) error
	ZGetCommentList(c context.Context, key string) ([]int64, []*auto.Comment, error)
}
