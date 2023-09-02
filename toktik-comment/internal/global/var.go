package global

import (
	"toktik-comment/internal/model/config"
	"toktik-common/logger"
	"toktik-common/utils"
)

var (
	Settings  config.Config // Public配置
	Logger    *logger.Log   // 日志
	SnowFlake *utils.SnowFlake
)
