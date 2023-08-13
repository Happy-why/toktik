package global

import (
	"github.com/Happy-Why/toktik-common/logger"
	"github.com/Happy-Why/toktik-common/token"
	"github.com/Happy-Why/toktik-common/utils"
	"github.com/Happy-Why/toktik-user/internal/model/config"
)

var (
	PbSettings config.Public  // Public配置
	PvSettings config.Private // Private配置
	Logger     *logger.Log    // 日志
	Maker      token.Maker
	SnowFlake  *utils.SnowFlake
)
