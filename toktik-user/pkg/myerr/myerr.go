package myerr

import "github.com/Happy-Why/toktik-common/errcode"

var (
	UsernameExist = errcode.NewErr(100001, "用户名已存在")
)
