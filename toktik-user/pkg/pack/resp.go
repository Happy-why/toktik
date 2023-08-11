package pack

// BuildBaseResp build baseResp from error
//func BuildBaseResp(err error, any2 any) *any {
//	if err == nil {
//		return baseResp(errcode.StatusOK)
//	}
//	e := errcode.NewMyErr()
//	if errors.As(err, &e) {
//		return baseResp(e)
//	}
//
//	s := errcode.ErrServer.WithError(err.Error())
//	return baseResp(s)
//}
//
//func baseResp(err errcode.Err) *any {
//	return &user.BaseResp{
//		StatusCode:    int64(err.ECode()),
//		StatusMessage: err.Error(),
//		ServiceTime:   time.Now().Unix(),
//	}
//}
