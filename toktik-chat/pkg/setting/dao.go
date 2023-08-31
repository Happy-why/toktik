package setting

import (
	"toktik-chat/internal/dao"
	"toktik-chat/internal/dao/cron"
	"toktik-chat/internal/dao/mysql"
	"toktik-chat/internal/dao/redis"
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
	go cron.TimingJob()
}
