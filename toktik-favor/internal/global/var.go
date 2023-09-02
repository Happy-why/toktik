package global

import (
	"toktik-common/logger"
	"toktik-favor/internal/model/config"
)

var (
	Settings config.Config // Public配置
	Logger   *logger.Log   // 日志
)
