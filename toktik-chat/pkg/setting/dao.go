package setting

import (
	"toktik-chat/internal/cache"
	"toktik-chat/internal/dao"
	"toktik-chat/internal/dao/cron"
	"toktik-chat/internal/dao/mysql"
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
	go cron.TimingJob()
}
