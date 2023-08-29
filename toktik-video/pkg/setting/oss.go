package setting

import (
	"toktik-common/oss/aliyun"
	"toktik-video/internal/global"
)

func init() {
	o := oss{}
	Settings = append(Settings, o)
}

type oss struct {
}

func (oss) InitSetting() {
	global.OSS = aliyun.Init(aliyun.Config{
		Endpoint:        global.Settings.AliyunOSS.Endpoint,
		AccessKeyId:     global.Settings.AliyunOSS.AccessKeyId,
		AccessKeySecret: global.Settings.AliyunOSS.AccessKeySecret,
		BucketName:      global.Settings.AliyunOSS.BucketName,
	})
}
