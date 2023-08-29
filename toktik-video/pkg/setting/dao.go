package setting

import (
	"toktik-video/internal/dao"
	"toktik-video/internal/dao/mysql"
	"toktik-video/internal/dao/redis"
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
