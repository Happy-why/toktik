package myerr

import "toktik-common/errcode"

var (
	IsNotFriend        = errcode.NewErr(400001, "不是朋友关系！")
	MessageCanNotEmpty = errcode.NewErr(400002, "不能发空消息哦~！")
)
