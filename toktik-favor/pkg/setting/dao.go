package setting

import (
	"toktik-favor/internal/cache"
	"toktik-favor/internal/dao"
	"toktik-favor/internal/dao/cron"
	"toktik-favor/internal/dao/mysql"
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
	go cron.TimingJob()
}
