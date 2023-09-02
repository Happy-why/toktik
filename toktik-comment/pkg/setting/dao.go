package setting

import (
	"toktik-comment/internal/cache"
	"toktik-comment/internal/dao"
	"toktik-comment/internal/dao/mysql"
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
	//go cron.TimingJob()
}
