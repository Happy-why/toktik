package redis

import (
	"context"
	"strconv"
	"time"
)

type InteractionRClient struct {
	rClient *RdbCache
}

func NewInteractionRClient() *InteractionRClient {
	return &InteractionRClient{
		rClient: GetRdbCache(),
	}
}

func (ir *InteractionRClient) KeyExist(c context.Context, key string) (bool, error) {
	exist, err := ir.rClient.KeyExist(c, key)
	if err != nil {
		return false, err
	}
	if exist == 1 {
		return true, nil
	}
	return false, nil
}

func (ir *InteractionRClient) SAddFollow(c context.Context, key string, targetId int64) error {
	_, err := ir.rClient.SAdd(c, key, targetId)
	if err != nil {
		return err
	}
	_, err = ir.rClient.Expire(c, key, time.Hour*48)
	return err
}

func (ir *InteractionRClient) SAddManyIds(c context.Context, key string, Ids []int64) error {
	Str := make([]string, len(Ids))
	for i, v := range Ids {
		Str[i] = strconv.FormatInt(v, 10)
	}
	_, err := ir.rClient.SAdd(c, key, Str)
	if err != nil {
		return err
	}
	_, err = ir.rClient.Expire(c, key, time.Hour*48)
	return err
}

func (ir *InteractionRClient) SAddFriend(c context.Context, key string, targetId int64) error {
	_, err := ir.rClient.SAdd(c, key, targetId)
	if err != nil {
		return err
	}
	_, err = ir.rClient.Expire(c, key, time.Hour*48)
	return err
}

func (ir *InteractionRClient) DelFollow(c context.Context, key string) error {
	return ir.rClient.Del(c, key)
}

func (ir *InteractionRClient) DelFan(c context.Context, key string) error {
	return ir.rClient.Del(c, key)
}

func (ir *InteractionRClient) DelFriend(c context.Context, key string) error {
	return ir.rClient.Del(c, key)
}

func (ir *InteractionRClient) IsFollow(c context.Context, key string, targetId int64) (bool, error) {
	return ir.rClient.SIsExist(c, key, targetId)
}

func (ir *InteractionRClient) IsFriend(c context.Context, key string, targetId int64) (bool, error) {
	return ir.rClient.SIsExist(c, key, targetId)
}

func (ir *InteractionRClient) SGetAllIds(c context.Context, key string) ([]int64, error) {
	Str, err := ir.rClient.SGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(Str) == 0 {
		return nil, nil
	}
	Ids := make([]int64, len(Str))
	for i, v := range Str {
		Ids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	return Ids, err
}
