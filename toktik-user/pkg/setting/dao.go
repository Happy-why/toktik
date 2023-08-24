package setting

import (
	"github.com/Happy-Why/toktik-user/internal/dao"
	"github.com/Happy-Why/toktik-user/internal/dao/mysql"
	"github.com/Happy-Why/toktik-user/internal/dao/redis"
)

func init() {
	d := Dao{}
	Settings = append(Settings, d)
}

type Dao struct {
}

func (Dao) InitSetting() {
	mysql.InitMysql()
	//global.RdbClient = redis.InitRedis()
	dao.Group.Rdb = redis.InitRedis()
}
