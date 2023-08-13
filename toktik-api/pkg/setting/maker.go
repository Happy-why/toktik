package setting

import (
	"github.com/Happy-Why/toktik-api/internal/global"
	"github.com/Happy-Why/toktik-common/token"
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
	global.Maker, err = token.NewPasetoMaker(global.PvSettings.Token.Key)
	if err != nil {
		panic(err)
	}
}
