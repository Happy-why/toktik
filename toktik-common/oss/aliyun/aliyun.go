package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Config struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

type OSS struct {
	config Config
}

func Init(config Config) *OSS {
	return &OSS{config: config}
}

func (o *OSS) newBucket() (*oss.Bucket, error) {
	// 创建OSSClient实例。
	client, err := oss.New(o.config.Endpoint, o.config.AccessKeyId, o.config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(o.config.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
