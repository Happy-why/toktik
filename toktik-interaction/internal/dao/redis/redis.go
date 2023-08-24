package redis

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-interaction/internal/dao"
	"github.com/Happy-Why/toktik-interaction/internal/global"
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
