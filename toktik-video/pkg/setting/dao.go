package setting

import (
	"github.com/Happy-Why/toktik-video/internal/dao"
	"github.com/Happy-Why/toktik-video/internal/dao/mysql"
	"github.com/Happy-Why/toktik-video/internal/dao/redis"
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
