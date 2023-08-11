package setting

import (
	"flag"
	"github.com/Happy-Why/toktik-common/setting"
	"github.com/Happy-Why/toktik-user/internal/global"
	"strings"
)

func init() {
	c := config{}
	Settings = append(Settings, c)
}

var (
	configPaths       string // 配置文件路径
	privateConfigName string // private配置文件名
	publicConfigName  string // public配置文件名
	configType        string // 配置文件类型
)

type config struct {
}

func setupFlag() {
	// 命令行参数绑定
	flag.StringVar(&privateConfigName, "private_name", "private", "private配置文件名")
	flag.StringVar(&publicConfigName, "public_name", "public", "public配置文件名")
	flag.StringVar(&configType, "type", "yaml", "配置文件类型")
	flag.StringVar(&configPaths, "path", global.RootDir+"/config", "指定要使用的配置文件路径,多个路径用逗号隔开")
	flag.Parse()
}

// InitSetting  读取配置文件
func (config) InitSetting() {
	setupFlag()
	// 在调用其他组件的Init时，这个init会首先执行并且把配置文件绑定到全局的结构体上
	PbSetting, err := setting.NewSetting(publicConfigName, configType, strings.Split(configPaths, ",")...) // 引入配置文件路径
	if err != nil {
		panic("初始化配置文件有误:" + err.Error())
	}
	if err = PbSetting.BindAll(&global.PbSettings); err != nil {
		panic("初始化配置文件有误:" + err.Error())
	}

	PvSetting, err := setting.NewSetting(privateConfigName, configType, strings.Split(configPaths, ",")...) // 引入配置文件路径
	if err != nil {
		panic("初始化配置文件有误:" + err.Error())
	}
	if err = PvSetting.BindAll(&global.PvSettings); err != nil {
		panic("初始化配置文件有误:" + err.Error())
	}
}
