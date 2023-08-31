package redis

import (
	"context"
	"strconv"
	"strings"
	"time"
	"toktik-chat/internal/model/auto"
)

type ChatRClient struct {
	rClient *RdbCache
}

func NewChatRClient() *ChatRClient {
	return &ChatRClient{
		rClient: GetRdbCache(),
	}
}

func (cr *ChatRClient) PushHistoryMessage(c context.Context, key string, createTime float64, content string) error {
	_, err := cr.rClient.ZSet(c, key, []float64{createTime}, []interface{}{content})
	if err != nil {
		return err
	}
	_, err = cr.rClient.Expire(c, key, 48*time.Hour)
	return err
}

func (cr *ChatRClient) PushManyHistoryMessage(c context.Context, key string, messageInfos []*auto.Message) error {
	scores := make([]float64, len(messageInfos))
	members := make([]interface{}, len(messageInfos))
	for i, v := range messageInfos {
		content := auto.CreateMessageContent(int64(v.UserId), int64(v.ToUserId), v.Content)
		scores[i] = float64(v.CreatedAt.Unix())
		members[i] = content
	}
	_, err := cr.rClient.ZSet(c, key, scores, members)
	if err != nil {
		return err
	}
	_, err = cr.rClient.Expire(c, key, 48*time.Hour)
	return err
}

func (cr *ChatRClient) PushDBMessage(c context.Context, key string, time float64, content string) error {
	_, err := cr.rClient.ZSet(c, key, []float64{time}, []interface{}{content})
	return err
}

func (cr *ChatRClient) ZRangeMessageList(c context.Context, key string) ([]*auto.Message, error) {
	messageList, err := cr.rClient.ZGetRangeWithScores(c, key, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(messageList) == 0 {
		return nil, nil
	}
	messageInfos := make([]*auto.Message, len(messageList))
	for i, message := range messageList {
		createdTime := time.Unix(int64(message.Score), 0)
		str := strings.SplitN(message.Member.(string), "+", 3)
		userId, _ := strconv.ParseInt(str[0], 10, 64)
		targetId, _ := strconv.ParseInt(str[1], 10, 64)
		content := str[2]
		messageInfos[i] = &auto.Message{
			BaseModel: auto.BaseModel{CreatedAt: createdTime},
			UserId:    uint(userId),
			ToUserId:  uint(targetId),
			Content:   content,
		}
	}
	return messageInfos, nil
}
