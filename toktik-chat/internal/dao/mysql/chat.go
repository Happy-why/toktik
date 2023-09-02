package mysql

import (
	"context"
	"time"
	"toktik-chat/internal/model/auto"
)

type ChatDao struct {
	conn *GormConn
}

func NewChatDao() *ChatDao {
	return &ChatDao{
		conn: NewGormConn(),
	}
}

// 按升序取消息

func (cd *ChatDao) GetMessageList(c context.Context, userId, targetId int64, preMsgTime int64) ([]*auto.Message, error) {
	messageList := make([]*auto.Message, 0)
	preMsgTime = preMsgTime / 100
	err := cd.conn.Session(c).Model(&auto.Message{}).
		Where("user_id = ? AND to_user_id = ?", userId, targetId).
		Where("created_at > ?", time.Unix(preMsgTime, 0)).
		Order("created_at asc").Find(&messageList).Error
	return nil, err
}
