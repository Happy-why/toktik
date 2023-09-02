package cache

import (
	"context"
	"strconv"
	"time"
)

type InteractionrCache struct {
	rCache *RdbCache
}

func NewInteractionrCache() *InteractionrCache {
	return &InteractionrCache{
		rCache: GetRdbCache(),
	}
}

func (ir *InteractionrCache) KeyExist(c context.Context, key string) (bool, error) {
	exist, err := ir.rCache.KeyExist(c, key)
	if err != nil {
		return false, err
	}
	if exist == 1 {
		return true, nil
	}
	return false, nil
}

func (ir *InteractionrCache) SAddFollow(c context.Context, key string, targetId int64) error {
	_, err := ir.rCache.SAdd(c, key, targetId)
	if err != nil {
		return err
	}
	_, err = ir.rCache.Expire(c, key, time.Hour*48)
	return err
}

func (ir *InteractionrCache) SAddManyIds(c context.Context, key string, Ids []int64) error {
	Str := make([]string, len(Ids))
	for i, v := range Ids {
		Str[i] = strconv.FormatInt(v, 10)
	}
	_, err := ir.rCache.SAdd(c, key, Str)
	if err != nil {
		return err
	}
	_, err = ir.rCache.Expire(c, key, time.Hour*48)
	return err
}

func (ir *InteractionrCache) SAddFriend(c context.Context, key string, targetId int64) error {
	_, err := ir.rCache.SAdd(c, key, targetId)
	if err != nil {
		return err
	}
	_, err = ir.rCache.Expire(c, key, time.Hour*48)
	return err
}

func (ir *InteractionrCache) DelFollow(c context.Context, key string) error {
	return ir.rCache.Del(c, key)
}

func (ir *InteractionrCache) DelFan(c context.Context, key string) error {
	return ir.rCache.Del(c, key)
}

func (ir *InteractionrCache) DelFriend(c context.Context, key string) error {
	return ir.rCache.Del(c, key)
}

func (ir *InteractionrCache) IsFollow(c context.Context, key string, targetId int64) (bool, error) {
	return ir.rCache.SIsExist(c, key, targetId)
}

func (ir *InteractionrCache) IsFriend(c context.Context, key string, targetId int64) (bool, error) {
	return ir.rCache.SIsExist(c, key, targetId)
}

func (ir *InteractionrCache) SGetAllIds(c context.Context, key string) ([]int64, error) {
	Str, err := ir.rCache.SGetAll(c, key)
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
