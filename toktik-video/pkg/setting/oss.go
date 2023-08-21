package setting

import (
	"github.com/Happy-Why/toktik-common/oss/aliyun"
	"github.com/Happy-Why/toktik-video/internal/global"
)

func init() {
	o := oss{}
	Settings = append(Settings, o)
}

type oss struct {
}

func (oss) InitSetting() {
	global.OSS = aliyun.Init(aliyun.Config{
		Endpoint:        global.PvSettings.AliyunOSS.Endpoint,
		AccessKeyId:     global.PvSettings.AliyunOSS.AccessKeyId,
		AccessKeySecret: global.PvSettings.AliyunOSS.AccessKeySecret,
		BucketName:      global.PvSettings.AliyunOSS.BucketName,
	})
}
