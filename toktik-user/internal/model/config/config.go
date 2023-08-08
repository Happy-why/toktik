package config

import "time"

type Public struct {
	Logger Logger
}

type Private struct {
	Token Token
	Mysql Mysql
	Redis Redis
}

type Server struct {
	Name                  string
	RunMode               string
	Addr                  string
	DefaultContextTimeout time.Duration
}

type Logger struct {
	Level         string
	LogSavePath   string // 保存路径
	LogFileExt    string // 日志文件后缀
	MaxSize       int    // 备份的大小(M)
	MaxBackups    int    // 最大备份数
	MaxAge        int    // 最大备份天数
	Compress      bool   // 是否压缩过期日志
	LowLevelFile  string // 低级别文件名
	HighLevelFile string // 高级别文件名
}

type Token struct {
	Key              string
	UserTokenExp     time.Duration
	AuthorizationKey string
}

type Mysql struct {
	Username string
	Password string
	Host     string
	Port     string
	DB       string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}
