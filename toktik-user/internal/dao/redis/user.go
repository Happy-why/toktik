package redis

import (
	"context"
	"fmt"
	"time"
	"toktik-user/internal/model/auto"
)

type UserRClient struct {
	rClient *RdbCache
}

func NewUserRClient() *UserRClient {
	return &UserRClient{
		rClient: GetRdbCache(),
	}
}

func (u *UserRClient) HSetUserInfo(c context.Context, key string, value map[string]interface{}) error {
	err := u.rClient.HSet(c, key, value)
	if err != nil {
		return err
	}
	_, err = u.rClient.Expire(c, key, time.Hour*168)
	return err
}

func (u *UserRClient) HSetUserCountInfo(c context.Context, key string, value map[string]interface{}) error {
	return u.rClient.HSet(c, key, value) // 需要频繁更改的信息，不设置过期时间
}

func (u *UserRClient) HGetUserInfo(c context.Context, key string) (*auto.User, error) {
	userMap, err := u.rClient.HGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return auto.CreateUserInfo(userMap)
}

func (u *UserRClient) HGetUserCountInfo(c context.Context, key string) (*auto.UserCount, error) {
	userMap, err := u.rClient.HGetAll(c, key)
	fmt.Println("userMap:", userMap)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return auto.CreateUserCountInfo(userMap)
}

func (u *UserRClient) AddFollowCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.FollowCount, 1)
	return err
}

func (u *UserRClient) AddFollowerCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.FollowerCount, 1)
	return err
}

func (u *UserRClient) SubFollowCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.FollowCount, -1)
	return err
}

func (u *UserRClient) SubFollowerCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.FollowerCount, -1)
	return err
}

func (u *UserRClient) AddWorkCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.WorkCount, 1)
	return err
}

func (u *UserRClient) AddFavoriteCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.FavoriteCount, 1)
	return err
}

func (u *UserRClient) AddTotalFavoriteCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.TotalFavorited, 1)
	return err
}

func (u *UserRClient) SubFavoriteCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.FavoriteCount, -1)
	return err
}

func (u *UserRClient) SubTotalFavoriteCount(c context.Context, key string) error {
	_, err := u.rClient.IncrHMCount(c, key, auto.TotalFavorited, -1)
	return err
}
