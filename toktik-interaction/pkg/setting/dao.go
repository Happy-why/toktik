package setting

import (
	"toktik-interaction/internal/cache"
	"toktik-interaction/internal/dao"
	"toktik-interaction/internal/dao/mysql"
)

func init() {
	d := Dao{}
	Settings = append(Settings, d)
}

type Dao struct {
}

func (Dao) InitSetting() {
	mysql.InitMysql()
	dao.Group.Rdb = cache.InitRedis()
}
