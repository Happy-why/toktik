package redis

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
	"toktik-user/internal/dao"
	"toktik-user/internal/global"
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
		Addr:     global.Settings.Redis.Host + ":" + global.Settings.Redis.Port,
		Password: global.Settings.Redis.Password, // 密码
		DB:       global.Settings.Redis.DB,       // 数据库
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

func (rc *RdbCache) Expire(c context.Context, key string, expireTime time.Duration) (bool, error) {
	// 有目标key返回true，没有目标key，返回false
	return rc.rdb.Expire(c, key, expireTime).Result()
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

func (rc *RdbCache) SAdd(c context.Context, key string, value interface{}) (int64, error) {
	// 返回存入数据的数量
	return rc.rdb.SAdd(c, key, value).Result()
}

func (rc *RdbCache) SGetAll(c context.Context, key string) ([]string, error) {
	return rc.rdb.SMembers(c, key).Result()
}

func (rc *RdbCache) SDel(c context.Context, key string, value interface{}) (int64, error) {
	return rc.rdb.SRem(c, key, value).Result()
}

func (rc *RdbCache) SIsExist(c context.Context, key string, value interface{}) (bool, error) {
	return rc.rdb.SIsMember(c, key, value).Result()
}

func (rc *RdbCache) ZSet(c context.Context, key string, score []float64, member []interface{}) (int64, error) {
	// 返回存入数据的数量
	z := make([]*redis2.Z, 0)
	for i := 0; i < len(score); i++ {
		z = append(z, &redis2.Z{Score: score[i], Member: member[i]})
	}
	return rc.rdb.ZAdd(c, key, z...).Result()
}

func (rc *RdbCache) ZGetRange(c context.Context, key string, min, max string, offset, count int64) ([]string, error) {
	return rc.rdb.ZRevRangeByScore(c, key, &redis2.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count}).Result()
}
