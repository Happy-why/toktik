package cache

import (
	"context"
	"fmt"
	"time"
	"toktik-user/internal/model/auto"
)

type UserCache struct {
	rCache *RdbCache
}

func NewUserCache() *UserCache {
	return &UserCache{
		rCache: GetRdbCache(),
	}
}

func (u *UserCache) SetToken(c context.Context, key, value string, expire time.Duration) error {
	return u.rCache.Set(c, key, value, expire)
}

func (u *UserCache) GetToken(c context.Context, key string) (string, error) {
	return u.rCache.Get(c, key)
}

func (u *UserCache) HSetUserInfo(c context.Context, key string, value map[string]interface{}) error {
	err := u.rCache.HSet(c, key, value)
	if err != nil {
		return err
	}
	_, err = u.rCache.Expire(c, key, time.Hour*168)
	return err
}

func (u *UserCache) HSetUserCountInfo(c context.Context, key string, value map[string]interface{}) error {
	return u.rCache.HSet(c, key, value) // 需要频繁更改的信息，不设置过期时间
}

func (u *UserCache) HGetUserInfo(c context.Context, key string) (*auto.User, error) {
	userMap, err := u.rCache.HGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return auto.CreateUserInfo(userMap)
}

func (u *UserCache) HGetUserCountInfo(c context.Context, key string) (*auto.UserCount, error) {
	userMap, err := u.rCache.HGetAll(c, key)
	fmt.Println("userMap:", userMap)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return auto.CreateUserCountInfo(userMap)
}

func (u *UserCache) AddFollowCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.FollowCount, 1)
	return err
}

func (u *UserCache) AddFollowerCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.FollowerCount, 1)
	return err
}

func (u *UserCache) SubFollowCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.FollowCount, -1)
	return err
}

func (u *UserCache) SubFollowerCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.FollowerCount, -1)
	return err
}

func (u *UserCache) AddWorkCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.WorkCount, 1)
	return err
}

func (u *UserCache) AddFavoriteCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.FavoriteCount, 1)
	return err
}

func (u *UserCache) AddTotalFavoriteCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.TotalFavorited, 1)
	return err
}

func (u *UserCache) SubFavoriteCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.FavoriteCount, -1)
	return err
}

func (u *UserCache) SubTotalFavoriteCount(c context.Context, key string) error {
	_, err := u.rCache.IncrHMCount(c, key, auto.TotalFavorited, -1)
	return err
}
