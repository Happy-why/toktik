package setting

import (
	"toktik-video/internal/cache"
	"toktik-video/internal/dao"
	"toktik-video/internal/dao/cron"
	"toktik-video/internal/dao/mysql"
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
