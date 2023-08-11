package config

import "time"

type Public struct {
	Server Server
	Logger Logger
	Rpc    Rpc
}

type Private struct {
	Token Token
	Etcd  Etcd
}

type Server struct {
	Name                  string
	RunMode               string
	Addr                  string
	DefaultContextTimeout time.Duration
}

type Rpc struct {
	Name        string
	Addr        string
	ServerAddrs map[string]string
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

type Etcd struct {
	Addr []string
}

type Token struct {
	Key              string
	UserTokenExp     time.Duration
	AuthorizationKey string
}
