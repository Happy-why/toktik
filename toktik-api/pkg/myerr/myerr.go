package myerr

import "toktik-common/errcode"

var (
	RegisterFailed        = errcode.NewErr(100001, "用户注册失败")
	ErrUserNameORPassWord = errcode.NewErr(100002, "账号或密码格式错误")
	CanNotChatSelf        = errcode.NewErr(100003, "不能自言自语哦~")
)
