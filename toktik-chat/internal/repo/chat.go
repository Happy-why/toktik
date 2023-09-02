package repo

import (
	"context"
	"toktik-chat/internal/model/auto"
)

type ChatRepo interface {
	GetMessageList(c context.Context, userId, targetId int64, preMsgTime int64) ([]*auto.Message, error)
}
