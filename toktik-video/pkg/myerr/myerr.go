package myerr

import "github.com/Happy-Why/toktik-common/errcode"

var (
	UploadVideoErr = errcode.NewErr(300001, "上传视频失败")
	GetCoverErr    = errcode.NewErr(300002, "获取视频封面失败")
	UploadCoverErr = errcode.NewErr(300003, "上传视频封面失败")
	SubmitVideoErr = errcode.NewErr(300004, "投稿视频失败")
)
