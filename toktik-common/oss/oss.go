package oss

import "io"

// OSS 对象存储接口
type OSS interface {
	UploadFile(objectKey string, reader io.Reader) (string, error)
	UploadByteFile(objectKey string, fileBuf []byte) (string, error)
	DeleteFile() error // 未实现
	CreateFileURL(bucketName, endpoint, objectName string) string
	CreateObjectKey(suffix string, directories ...string) string
}

// CreateObjectKey(后缀 string,...string])string
// 这里只是构建视频和封面的key，后续可以完善这个方法，参数为...string
// 每个string都是一级目录，然后用字符'/'拼接这个切片。
