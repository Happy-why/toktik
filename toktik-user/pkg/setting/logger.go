package setting

import (
	"github.com/Happy-Why/toktik-common/logger"
	"github.com/Happy-Why/toktik-user/internal/global"
	"go.uber.org/zap"
)

func init() {
	l := log{}
	Settings = append(Settings, l)
}

type log struct {
}

// InitSetting 日志初始化
func (log) InitSetting() {
	global.Logger = logger.NewLogger(&logger.InitStruct{
		LogSavePath:   global.PbSettings.Logger.LogSavePath,
		LogFileExt:    global.PbSettings.Logger.LogFileExt,
		MaxSize:       global.PbSettings.Logger.MaxSize,
		MaxBackups:    global.PbSettings.Logger.MaxBackups,
		MaxAge:        global.PbSettings.Logger.MaxAge,
		Compress:      global.PbSettings.Logger.Compress,
		LowLevelFile:  global.PbSettings.Logger.LowLevelFile,
		HighLevelFile: global.PbSettings.Logger.HighLevelFile,
	}, global.PbSettings.Logger.Level)
	zap.ReplaceGlobals(global.Logger.Logger)
}
