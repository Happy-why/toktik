package global

import (
	"github.com/Happy-Why/toktik-common/logger"
	"github.com/Happy-Why/toktik-common/oss"
	"github.com/Happy-Why/toktik-common/utils"
	"github.com/Happy-Why/toktik-video/internal/model/config"
)

var (
	PbSettings config.Public  // Public配置
	PvSettings config.Private // Private配置
	Logger     *logger.Log    // 日志
	OSS        oss.OSS
	SnowFlake  *utils.SnowFlake
)
