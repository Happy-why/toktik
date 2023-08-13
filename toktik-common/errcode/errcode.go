package errcode

import (
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
)

// 编写常用的一些错误处理公共方法，标准化我们的错误输出

type Err interface {
	Error() string
	ECode() int64
	WithDetails(details ...string) Err
}

var globalMap map[int64]Err
var once sync.Once

func NewErr(code int64, msg string) Err {
	once.Do(func() {
		globalMap = make(map[int64]Err)
	})
	if _, ok := globalMap[code]; ok {
		panic("错误码已存在")
	}
	err := &myErr{Code: code, Msg: msg}
	globalMap[code] = err
	return err
}

type myErr struct {
	Code    int64    `json:"status_code"` // 状态码，0-成功，其他值-失败
	Msg     string   `json:"status_msg"`  // 返回状态描述
	Details []string `json:"-"`           // 详细信息
}

func (m *myErr) ECode() int64 {
	return m.Code
}

func (m *myErr) Error() string {
	return fmt.Sprintf("%v", m.Msg)
}

func (m *myErr) WithDetails(details ...string) Err {
	var newErr = &myErr{}
	_ = copier.Copy(newErr, m)
	m.Msg = m.Msg + "," + details[0]
	newErr.Details = append(newErr.Details, details...)
	return newErr
}
