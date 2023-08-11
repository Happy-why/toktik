package myerr

import "github.com/Happy-Why/toktik-common/errcode"

var (
	RegisterFailed        = errcode.NewErr(1001, "用户注册失败")
	ErrUserNameORPassWord = errcode.NewErr(1002, "账号或密码格式错误")
)
