package global

import (
	"toktik-common/logger"
	"toktik-common/oss"
	"toktik-common/utils"
	"toktik-video/internal/model/config"
)

var (
	Settings  config.Config // Public配置
	Logger    *logger.Log   // 日志
	OSS       oss.OSS
	SnowFlake *utils.SnowFlake
)
