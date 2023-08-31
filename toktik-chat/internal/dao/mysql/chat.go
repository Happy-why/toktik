package mysql

import (
	"context"
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

func (cd *ChatDao) GetMessageList(c context.Context, userId, targetId int64) ([]*auto.Message, error) {
	messageList := make([]*auto.Message, 0)
	err := cd.conn.Session(c).Model(&auto.Message{}).
		Where("user_id = ? AND to_user_id = ?", userId, targetId).
		Order("created_at asc").Find(&messageList).Error
	return nil, err
}
