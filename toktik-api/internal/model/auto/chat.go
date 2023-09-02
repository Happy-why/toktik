package auto

import (
	"strconv"
)

type Message struct {
	BaseModel
	UserId   uint   `json:"from_user_id" gorm:"index:idx_message;not null"`
	ToUserId uint   `json:"to_user_id" gorm:"index:idx_message;not null"`
	Content  string `json:"content" gorm:"not null"`
}

func CreateChatHistoryKey(userId, targetId int64) string {
	if userId < targetId {
		return "chat::history::" + strconv.FormatInt(userId, 10) + "+" + strconv.FormatInt(targetId, 10)
	}
	return "chat::history::" + strconv.FormatInt(targetId, 10) + "+" + strconv.FormatInt(userId, 10)
}

func CreateChatMessageKey(userId, targetId int64) string {
	if userId < targetId {
		return "chat::message::" + strconv.FormatInt(userId, 10) + "+" + strconv.FormatInt(targetId, 10)
	}
	return "chat::message::" + strconv.FormatInt(targetId, 10) + "+" + strconv.FormatInt(userId, 10)
}

func CreateMessageContent(userId, targetId int64, content string) string {
	return strconv.FormatInt(userId, 10) + "+" + strconv.FormatInt(targetId, 10) + "+" + content
}
