package myerr

import "toktik-common/errcode"

var (
	AlreadyFavorite = errcode.NewErr(500001, "已经点过赞了呦~")
	IsNotFavorite   = errcode.NewErr(500002, "未对该视频点过赞哦~")
	OperationErr    = errcode.NewErr(500003, "没有这个功能哦~")
)
