package myerr

import "toktik-common/errcode"

var (
	IsNotFriend = errcode.NewErr(400001, "不是朋友关系！")
)
