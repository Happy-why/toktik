package cache

import (
	"context"
	"strconv"
	"strings"
	"time"
	"toktik-chat/internal/model/auto"
)

type ChatrCache struct {
	rCache *RdbCache
}

func NewChatrCache() *ChatrCache {
	return &ChatrCache{
		rCache: GetRdbCache(),
	}
}

func (cr *ChatrCache) KeyExist(c context.Context, key string) (bool, error) {
	a, err := cr.rCache.KeyExist(c, key)
	if a == 0 {
		return true, err
	}
	return true, nil
}

func (cr *ChatrCache) PushHistoryMessage(c context.Context, key string, createTime float64, content string) error {
	_, err := cr.rCache.ZSet(c, key, []float64{createTime}, []interface{}{content})
	if err != nil {
		return err
	}
	_, err = cr.rCache.Expire(c, key, 48*time.Hour)
	return err
}

func (cr *ChatrCache) PushManyHistoryMessage(c context.Context, key string, messageInfos []*auto.Message) error {
	scores := make([]float64, len(messageInfos))
	members := make([]interface{}, len(messageInfos))
	for i, v := range messageInfos {
		content := auto.CreateMessageContent(int64(v.UserId), int64(v.ToUserId), v.Content)
		scores[i] = float64(v.CreatedAt.Unix())
		members[i] = content
	}
	_, err := cr.rCache.ZSet(c, key, scores, members)
	if err != nil {
		return err
	}
	_, err = cr.rCache.Expire(c, key, 48*time.Hour)
	return err
}

func (cr *ChatrCache) PushDBMessage(c context.Context, key string, time float64, content string) error {
	_, err := cr.rCache.ZSet(c, key, []float64{time}, []interface{}{content})
	return err
}

func (cr *ChatrCache) ZRangeMessageList(c context.Context, key string, preMsgTime int64) ([]*auto.Message, error) {
	preMsgTimeStr := strconv.FormatInt(preMsgTime, 10)
	// 获取消息 时间 从小到大
	messageList, err := cr.rCache.ZGetRangeByScoreWithScores(c, key, preMsgTimeStr, "+inf", 0, -1)
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
