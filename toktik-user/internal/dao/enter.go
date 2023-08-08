package dao

import (
	"github.com/Happy-Why/toktik-user/internal/dao/redis"
	"gorm.io/gorm"
)

type group struct {
	Mdb *gorm.DB
	Rdb *redis.RdbCache
}

var Group = new(group)
