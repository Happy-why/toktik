package redis

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-user/internal/model/auto"
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
	return u.rClient.HSet(c, key, value)
}

func (u *UserRClient) HGetUserInfo(c context.Context, key string) (*auto.User, error) {
	userMap, err := u.rClient.HGetAll(c, key)
	fmt.Println("userMap:", userMap)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, nil
	}
	return auto.CreateUserInfo(userMap)
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
