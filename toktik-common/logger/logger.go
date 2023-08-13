package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Log struct {
	*zap.Logger
}

var initLog = new(InitStruct)

type InitStruct struct {
	LogSavePath   string // 保存路径
	LogFileExt    string // 日志文件后缀
	MaxSize       int    // 备份的大小(M)
	MaxBackups    int    // 最大备份数
	MaxAge        int    // 最大备份天数
	Compress      bool   // 是否压缩过期日志
	LowLevelFile  string // 低级别文件名
	HighLevelFile string // 高级别文件名
}

func NewLogger(x *InitStruct, level string) *Log {
	initLog = x
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})
	// 多个日志文件
	var cores []zapcore.Core
	lowFileWriteSyncer := getLogWriter(initLog.LogSavePath + initLog.LowLevelFile + initLog.LogFileExt)
	highFileWriteSyncer := getLogWriter(initLog.LogSavePath + initLog.HighLevelFile + initLog.LogFileExt)
	encoder := getEncoder()
	lowFileCore := zapcore.NewCore(encoder, lowFileWriteSyncer, lowPriority)
	highFileCore := zapcore.NewCore(encoder, highFileWriteSyncer, highPriority)
	cores = append(cores, lowFileCore, highFileCore)
	if level == "debug" { // 如果是debug级别好需要输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}
	core := zapcore.NewTee(cores...)
	return &Log{zap.New(core, zap.AddCaller())} // 增加函数调用信息
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,   // 结尾字符
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 将Level序列化为全大写字符串。例如, InfoLevel被序列化为INFO
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // 格式化时间戳
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 普通的Log Encoder
}

func getLogWriter(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{ // 日志切割
		Filename:   filename,
		MaxSize:    initLog.MaxSize,    // M
		MaxBackups: initLog.MaxBackups, // 备份数量
		MaxAge:     initLog.MaxAge,     // 最大备份天数
		Compress:   initLog.Compress,   // 压缩过期日志
	}
	return zapcore.AddSync(lumberJackLogger)
}
