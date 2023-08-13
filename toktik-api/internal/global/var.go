package global

import (
	"github.com/Happy-Why/toktik-api/internal/model/config"
	"github.com/Happy-Why/toktik-common/logger"
	"github.com/Happy-Why/toktik-common/token"
)

var (
	PbSettings config.Public  // Public配置
	PvSettings config.Private // Private配置
	Maker      token.Maker
	Logger     *logger.Log // 日志
)
