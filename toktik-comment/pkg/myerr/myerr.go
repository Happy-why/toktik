package myerr

import "toktik-common/errcode"

var (
	OperationErr = errcode.NewErr(600001, "没有这个功能哦~")
)
