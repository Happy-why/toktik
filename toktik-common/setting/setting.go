package setting

import (
	"io"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 使用viper进行配置文件的读取和热加载

type Setting struct {
	vp  *viper.Viper
	all interface{}
}

// NewSetting 初始化本项目的配置的基础属性
// 设定配置文件的名称为 config，配置类型为 yaml，并且设置其配置路径为相对路径 configs/
func NewSetting(configName, configType string, configPaths ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName(configName)
	vp.SetConfigType(configType) // 设置配置文件类型
	for _, config := range configPaths {
		if config != "" {
			vp.AddConfigPath(config) // 可以设置多个配置路径,解决路径查找问题
		}
	}
	err := vp.ReadInConfig() // 加载配置文件
	if err != nil {
		return nil, err
	}
	s := &Setting{vp: vp}
	s.vp.WatchConfig()
	s.vp.OnConfigChange(func(in fsnotify.Event) {
		log.Println("更新配置")
		err := s.vp.Unmarshal(s.all)
		if err != nil {
			log.Fatalln("更新配置失败:" + err.Error())
		}
	})
	return s, nil
}

// ReadConfigFromBuf 从文件流中读取配置
func ReadConfigFromBuf(configType string, fileBuf io.Reader) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigType(configType)  // 设置配置文件类型
	err := vp.ReadConfig(fileBuf) // 加载配置文件
	if err != nil {
		return nil, err
	}
	s := &Setting{vp: vp}
	return s, nil
}

// BindAll 绑定配置文件
func (s *Setting) BindAll(v interface{}) error {
	// 绑定
	err := s.vp.Unmarshal(v)
	if err != nil {
		return err
	}
	s.all = v
	return nil
}
