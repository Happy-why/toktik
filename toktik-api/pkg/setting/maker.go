package setting

import (
	"toktik-api/internal/global"
	"toktik-common/token"
)

func init() {
	m := maker{}
	Settings = append(Settings, m)
}

type maker struct {
}

// InitSetting 初始化
func (maker) InitSetting() {
	var err error
	global.Maker, err = token.NewPasetoMaker(global.Settings.Token.Key)
	if err != nil {
		panic(err)
	}
}
