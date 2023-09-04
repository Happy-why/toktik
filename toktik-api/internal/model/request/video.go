package request

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

type VideoFeedRequest struct {
	LatestTime int64  `json:"latest_time" form:"latest_time" query:"latest_time"`
	Token      string `json:"token" form:"token" query:"token"`
	UserId     int64  `json:"user_id" form:"user_id" query:"user_id"`
}

func (v *VideoFeedRequest) VerifyFeed() {
	if v.LatestTime == 0 {
		v.LatestTime = time.Now().Unix()
	}
}

type VideoPublishRequest struct {
	Data   *multipart.FileHeader `json:"data" form:"data"`
	Token  string                `json:"token" form:"token"`
	Title  string                `json:"title" form:"title"`
	UserId int64                 `json:"user_id" form:"user_id"`
}

func (v *VideoPublishRequest) VerifyFeed() (*bytes.Buffer, error) {
	if v.Data == nil {
		return nil, fmt.Errorf("找不到上传的文件呀~")
	}
	src, _ := v.Data.Open()
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, src)
	return buf, err
}

type PublishListRequest struct {
	Token    string `json:"token" form:"token" query:"token"`
	UserId   int64  `json:"user_id" form:"user_id" query:"user_id"`
	MyUserId int64  `json:"my_user_id" form:"my_user_id" query:"my_user_id"`
}
