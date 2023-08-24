package redis

import (
	"context"
)

type VideoRClient struct {
	rClient *RdbCache
}

func NewVideoRClient() *VideoRClient {
	return &VideoRClient{
		rClient: GetRdbCache(),
	}
}

func (vr *VideoRClient) HSetVideoInfo(c context.Context, key string, value map[string]interface{}) error {
	return vr.rClient.HSet(c, key, value)
}

func (vr *VideoRClient) IncrUserFollowCount(c context.Context, key string) error {
	vr.rClient.rdb.Pipeline()
	return nil
}
