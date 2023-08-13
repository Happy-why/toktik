package errcode

var (
	StatusOK           = NewErr(0, "成功")
	ErrParamsNotValid  = NewErr(0001, "参数有误")
	ErrNotFound        = NewErr(0002, "未找到资源")
	ErrServer          = NewErr(0003, "系统错误")
	ErrTooManyRequests = NewErr(0004, "请求过多")
	ErrTimeOut         = NewErr(0005, "请求超时")
	ErrDB              = NewErr(0006, "db错误")
	ErrRedis           = NewErr(0007, "redis错误")
	ErrAuth            = NewErr(0010, "身份鉴权失败")
	ErrAuthExp         = NewErr(0011, "身份过期")
)
