package redis

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type RdbCache struct {
	rdb *redis.Client
}

func NewRdbCache(rdb *redis.Client) *RdbCache {
	return &RdbCache{rdb: rdb}
}

func InitRedis() *RdbCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.PvSettings.Redis.Host + ":" + global.PvSettings.Redis.Port,
		Password: global.PvSettings.Redis.Password, // 密码
		DB:       global.PvSettings.Redis.DB,       // 数据库
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		// Redis连接失败，进行相应处理
		zap.S().Errorf("redis init failed, err:%v", err)
		fmt.Println("redis初始化失败！！！！！")
		panic(err)
	}
	return NewRdbCache(rdb)
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
