package setting

import (
	"bytes"
	"flag"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	log2 "log"
	"strings"
	"toktik-api/internal/global"
	"toktik-common/setting"
)

var (
	bootstrap   string
	configPaths string // 配置文件路径
	configName  string // 配置文件名
	configType  string // 配置文件类型
)

func setupFlag() {
	// 命令行参数绑定
	flag.StringVar(&bootstrap, "bootstrap", "bootstrap", "配置文件名")
	flag.StringVar(&configName, "name", "config", "配置文件名")
	flag.StringVar(&configType, "type", "yaml", "配置文件类型")
	flag.StringVar(&configPaths, "path", global.RootDir+"/config", "指定要使用的配置文件路径,多个路径用逗号隔开")
	flag.Parse()
}

// InitBootStrap 读取配置文件
func InitBootStrap() {
	setupFlag()
	// 配置 nacos
	NacosConfig, err := setting.NewSetting(bootstrap, configType, strings.Split(configPaths, ",")...) // 引入配置文件路径
	if err != nil {
		log2.Fatalln("初始化配置文件有误:" + err.Error())
	}
	if err = NacosConfig.BindAll(&global.Settings); err != nil {
		panic("初始化配置文件有误:" + err.Error())
	}
	// 初始化nacos 客户端和服务端
	nacosClient := InitNacos()
	// 获取 远程的配置文件
	configYaml, err := nacosClient.confClient.GetConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  global.Settings.Nacos.Group,
	})
	if err != nil {
		log2.Fatalln(err)
	}
	// 监听远程配置文件的更改
	if err = nacosClient.confClient.ListenConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  global.Settings.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			//
			Setting, err := setting.ReadConfigFromBuf(configType, bytes.NewBuffer([]byte(configYaml)))
			if err != nil {
				log2.Printf("load nacos config changed err : %s \n", err.Error())
			}
			//所有的配置应该重新读取
			if err = Setting.BindAll(&global.Settings); err != nil {
				panic("初始化配置文件有误:" + err.Error())
			}
			//log2.Printf("load nacos config changed %s \n", data)
			log2.Printf("配置更改")
		},
	}); err != nil {
		log2.Fatalln(err)
	}
	if configYaml != "" {
		Setting, err := setting.ReadConfigFromBuf(configType, bytes.NewBuffer([]byte(configYaml)))
		if err != nil {
			log2.Fatalln(err)
		}
		if err = Setting.BindAll(&global.Settings); err != nil {
			panic("初始化配置文件有误:" + err.Error())
		}
	} else {
		// 在调用其他组件的Init时，这个init会首先执行并且把配置文件绑定到全局的结构体上
		Setting, err := setting.NewSetting(configName, configType, strings.Split(configPaths, ",")...) // 引入配置文件路径
		if err != nil {
			panic("初始化配置文件有误:" + err.Error())
		}
		if err = Setting.BindAll(&global.Settings); err != nil {
			panic("初始化配置文件有误:" + err.Error())
		}
	}
}
