package myerr

import "toktik-common/errcode"

var (
	UsernameExist  = errcode.NewErr(100001, "用户名已存在")
	CreateTokenErr = errcode.NewErr(100002, "生成Token失败")
	PasswordErr    = errcode.NewErr(100003, "密码错误")
	UserNotFound   = errcode.NewErr(100004, "人没了！！！")
)
