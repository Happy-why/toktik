package setting

import (
	"toktik-interaction/internal/dao"
	"toktik-interaction/internal/dao/mysql"
	"toktik-interaction/internal/dao/redis"
)

func init() {
	d := Dao{}
	Settings = append(Settings, d)
}

type Dao struct {
}

func (Dao) InitSetting() {
	mysql.InitMysql()
	dao.Group.Rdb = redis.InitRedis()
}
