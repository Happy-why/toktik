package redis

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-user/internal/dao"
	"github.com/Happy-Why/toktik-user/internal/global"
	redis2 "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type RdbCache struct {
	rdb *redis2.Client
}

func GetRdbCache() *RdbCache {
	return &RdbCache{rdb: dao.Group.Rdb}
}

func NewRdbCache(rdb *redis2.Client) *RdbCache {
	return &RdbCache{rdb: rdb}
}

func InitRedis() *redis2.Client {
	rdb := redis2.NewClient(&redis2.Options{
		Addr:     global.PvSettings.Redis.Host + ":" + global.PvSettings.Redis.Port,
		Password: global.PvSettings.Redis.Password, // 密码
		DB:       global.PvSettings.Redis.DB,       // 数据库
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		// Redis连接失败，进行相应处理
		zap.L().Error("redis init error:", zap.Error(err))
		fmt.Println("redis初始化失败！！！！！")
		panic(err)
	}
	return rdb
}

func (rc *RdbCache) Put(c context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(c, key, value, expire).Err()
	return err
}

func (rc *RdbCache) Get(c context.Context, key string) (string, error) {
	fmt.Println(key)
	result, err := rc.rdb.Get(c, key).Result()
	return result, err
}

func (rc *RdbCache) HSet(c context.Context, key string, value interface{}) error {
	return rc.rdb.HSet(c, key, value).Err()
}

func (rc *RdbCache) HGet(c context.Context, key string, filed string) error {
	return rc.rdb.HGet(c, key, filed).Err()
}

func (rc *RdbCache) HGetAll(c context.Context, key string) (map[string]string, error) {
	return rc.rdb.HGetAll(c, key).Result()
}

func (rc *RdbCache) IncrHMCount(c context.Context, key, field string, incr int64) (int64, error) {
	return rc.rdb.HIncrBy(c, key, field, incr).Result()
}
