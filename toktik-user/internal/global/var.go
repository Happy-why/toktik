package global

import (
	"toktik-common/logger"
	"toktik-common/token"
	"toktik-common/utils"
	"toktik-user/internal/model/config"
)

var (
	Settings  config.Config // Public配置
	Logger    *logger.Log   // 日志
	Maker     token.Maker
	SnowFlake *utils.SnowFlake
	//RdbClient  *redis.RdbCache
)
