package errcode

import (
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
)

// 编写常用的一些错误处理公共方法，标准化我们的错误输出

type Err interface {
	Error() string
	ECode() int
	WithDetails(details ...string) Err
}

var globalMap map[int]Err
var once sync.Once

func NewErr(code int, msg string) Err {
	once.Do(func() {
		globalMap = make(map[int]Err)
	})
	if _, ok := globalMap[code]; ok {
		panic("错误码已存在")
	}
	err := &myErr{Code: code, Msg: msg}
	globalMap[code] = err
	return err
}

type myErr struct {
	Code    int      `json:"status_code"` // 状态码，0-成功，其他值-失败
	Msg     string   `json:"status_msg"`  // 返回状态描述
	Details []string `json:"-"`           // 详细信息
}

func (m *myErr) ECode() int {
	return m.Code
}

func (m *myErr) Error() string {
	return fmt.Sprintf("错误码:%v,错误信息:%v,详细信息:%v", m.Code, m.Msg, m.Details)
}

func (m *myErr) WithDetails(details ...string) Err {
	var newErr = &myErr{}
	_ = copier.Copy(newErr, m)
	newErr.Details = append(newErr.Details, details...)
	return newErr
}
