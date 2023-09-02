package setting

import (
	"toktik-comment/internal/global"
	"toktik-common/utils"
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
