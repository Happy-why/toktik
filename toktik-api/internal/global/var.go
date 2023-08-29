package global

import (
	"toktik-api/internal/model/config"
	"toktik-common/logger"
	"toktik-common/token"
)

var (
	Settings config.Config // Public配置
	Maker    token.Maker
	Logger   *logger.Log // 日志
)
