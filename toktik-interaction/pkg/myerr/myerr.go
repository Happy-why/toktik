package myerr

import "github.com/Happy-Why/toktik-common/errcode"

var (
	UserAlreadyFollowed = errcode.NewErr(200001, "已经关注该用户")
	UserNotFollowed     = errcode.NewErr(200002, "该用户未被关注")
	NotCanFollowSelf    = errcode.NewErr(200003, "不可以对自己操作哦~")
	OperationErr        = errcode.NewErr(200004, "错误的操作")
)
