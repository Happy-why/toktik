package config

import "time"

type Config struct {
	Nacos      Nacos
	Server     Server
	Logger     Logger
	Rpc        Rpc
	Token      Token
	Etcd       Etcd
	Jaeger     Jaeger
	Limit      Limit
	Prometheus Prometheus
}

type Nacos struct {
	Namespace   string `json:"namespace"`
	Group       string `json:"group"`
	Addr        string `json:"addr"`
	Port        int    `json:"port"`
	Scheme      string `json:"scheme"`
	ContextPath string `json:"context_path"`
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

type Jaeger struct {
	HTTPEndpoint      string
	RPCExportEndpoint string
	ServerName        map[string]string
}

type Prometheus struct {
	Post string
	Path string
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
type Limit struct {
	IPLimit  IPLimit  `yaml:"IPLimit"`
	APILimit APILimit `yaml:"APILimit"`
}

type IPLimit struct {
	Cap     int64 `yaml:"Cap"`
	GenNum  int64 `yaml:"GenNum"`
	GenTime int64 `yaml:"GenTime"`
	Cost    int64 `yaml:"Cost"`
}

type APILimit struct {
	Upload []Bucket `yaml:"Upload"`
	Email  []Bucket `yaml:"Email"`
}

type Bucket struct {
	Count    int           `yaml:"Count"`
	Duration time.Duration `yaml:"Duration"`
	Burst    int           `yaml:"Burst"`
}
