package setting

import (
	"toktik-user/internal/cache"
	"toktik-user/internal/dao"
	"toktik-user/internal/dao/cron"
	"toktik-user/internal/dao/mysql"
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
	dao.Group.Rdb = cache.InitRedis()
	// 开启定时任务
	go cron.TimingJob()
}
