package myerr

import "toktik-common/errcode"

var (
	UsernameExist          = errcode.NewErr(100001, "用户名已存在")
	CreateTokenErr         = errcode.NewErr(100002, "生成Token失败")
	PasswordErr            = errcode.NewErr(100003, "密码错误")
	UserNotFound           = errcode.NewErr(100004, "用户名不存在！")
	IsFollowManyTargetsErr = errcode.NewErr(100005, "查询是否关注很多用户，失败！")
)
