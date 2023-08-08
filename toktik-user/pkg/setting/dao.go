package setting

import (
	"github.com/Happy-Why/toktik-user/internal/dao"
	"github.com/Happy-Why/toktik-user/internal/dao/mysql"
	"github.com/Happy-Why/toktik-user/internal/dao/redis"
)

type Dao struct {
}

func (Dao) Init() {
	mysql.InitMysql()
	dao.Group.Rdb = redis.InitRedis()
}
