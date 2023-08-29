package auto

import (
	"strconv"
	"time"
)

type Message struct {
	BaseModel
	UserId   uint   `json:"from_user_id"`
	ToUserId uint   `json:"to_user_id"`
	Content  string `json:"content"`
}

func CreateChatHistoryKey(userId, targetId int64) string {
	if userId < targetId {
		return "chat::history::" + strconv.Itoa(int(userId)) + "+" + strconv.Itoa(int(targetId))
	}
	return "chat::history::" + strconv.Itoa(int(targetId)) + "+" + strconv.Itoa(int(userId))
}

func CreateChatMessageKey(userId, targetId int64) string {
	if userId < targetId {
		return "chat::message::" + strconv.Itoa(int(userId)) + "+" + strconv.Itoa(int(targetId))
	}
	return "chat::message::" + strconv.Itoa(int(targetId)) + "+" + strconv.Itoa(int(userId))
}

func CreateMessageContent(content string) string {
	timeNow := time.Now().Unix()
	return strconv.FormatInt(timeNow, 10) + "$" + content
}
