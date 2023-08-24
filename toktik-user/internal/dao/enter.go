package dao

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type group struct {
	Mdb *gorm.DB
	Rdb *redis.Client
}

var Group = new(group)
