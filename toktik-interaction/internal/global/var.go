package global

import (
	"github.com/Happy-Why/toktik-common/logger"
	"github.com/Happy-Why/toktik-interaction/internal/model/config"
)

var (
	PbSettings config.Public  // Public配置
	PvSettings config.Private // Private配置
	Logger     *logger.Log    // 日志
)
