package setting

import (
	"toktik-common/utils"
	"toktik-video/internal/global"
)

func init() {
	s := sf{}
	Settings = append(Settings, s)
}

type sf struct {
}

func (sf) InitSetting() {
	var err error
	global.SnowFlake, err = utils.NewSnowFlake(0, 0)
	if err != nil {
		panic(err)
	}
}
