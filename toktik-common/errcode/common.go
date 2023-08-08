package errcode

var (
	StatusOK           = NewErr(0, "成功")
	ErrParamsNotValid  = NewErr(1001, "参数有误")
	ErrNotFound        = NewErr(1002, "未找到资源")
	ErrServer          = NewErr(1003, "系统错误")
	ErrTooManyRequests = NewErr(1004, "请求过多")
	ErrTimeOut         = NewErr(1005, "请求超时")
)
