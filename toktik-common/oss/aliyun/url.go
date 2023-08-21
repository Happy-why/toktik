package aliyun

import (
	"net/url"
	"strings"
)

//构建 URL https://why-bucket.oss-cn-beijing.aliyuncs.com/test1/123.jpg
// URL 由 https:// + bucket_name. + endpoint + object_name

func (o *OSS) CreateFileURL(bucketName, endpoint, objectName string) string {
	// 注意: OSS自己构建的 httpURL，会自动将 objectName转码一次，所以需要在自己构建的URL中进行转码，不能提前对objectName转码
	objectName = url.QueryEscape(objectName)
	objectName = strings.Replace(objectName, "%2F", "/", -1) // "/"转码后是"%2F"，让他变回去
	str := "https://" + bucketName + "." + endpoint + "/" + objectName
	// 对url的中文 进行url编码
	return str
}

//data := url.QueryEscape("asd")
//fmt.Println(data)
//ax, _ := url.QueryUnescape("%E5%85%B0%E4%BA%AD%E5%BA%8F")
//fmt.Println(ax)
