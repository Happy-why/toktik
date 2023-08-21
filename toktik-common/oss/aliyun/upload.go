package aliyun

import (
	"bytes"
	"errors"
	"io"
)

func (o *OSS) UploadFile(objectKey string, reader io.Reader) (string, error) {
	bucket, err := o.newBucket()
	if err != nil {
		return "", errors.New("function OSS.NewBucket() Failed, err:" + err.Error())
	}
	// 上传文件的bytes。
	err = bucket.PutObject(objectKey, reader)
	if err != nil {
		return "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}
	url := o.CreateFileURL(o.config.BucketName, o.config.Endpoint, objectKey)
	return url, nil
}

// UploadByteFile 上传文件 objectKey：目录+objectName fileBuf：文件流.buf
// 返回 访问地址，文件key，error
func (o *OSS) UploadByteFile(objectKey string, fileBuf []byte) (string, error) {
	bucket, err := o.newBucket()
	if err != nil {
		return "", errors.New("function OSS.NewBucket() Failed, err:" + err.Error())
	}
	// 上传文件的bytes。
	err = bucket.PutObject(objectKey, bytes.NewReader(fileBuf))
	if err != nil {
		return "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}
	url := o.CreateFileURL(o.config.BucketName, o.config.Endpoint, objectKey)
	return url, nil
}

// DeleteFile 删除文件
func (o *OSS) DeleteFile() error {
	return nil
}
