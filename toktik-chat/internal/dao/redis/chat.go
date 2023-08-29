package redis

import (
	"context"
	redis2 "github.com/go-redis/redis/v8"
)

type ChatRClient struct {
	rClient *RdbCache
}

func NewChatRClient() *ChatRClient {
	return &ChatRClient{
		rClient: GetRdbCache(),
	}
}

func (cr *ChatRClient) PushHistoryMessage(c context.Context, key string, time float64, content string) error {
	_, err := cr.rClient.ZSet(c, key, []float64{time}, []interface{}{content})
	return err
}

func (cr *ChatRClient) PushDBMessage(c context.Context, key string, time float64, content string) error {
	_, err := cr.rClient.ZSet(c, key, []float64{time}, []interface{}{content})
	return err
}

func (cr *ChatRClient) ZRangeMessageList(c context.Context, key string) ([]redis2.Z, error) {
	return cr.rClient.ZGetRangeWithScores(c, key, 0, -1)
}
