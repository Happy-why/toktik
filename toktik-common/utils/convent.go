package utils

import (
	"strconv"
	"strings"
)

// 统一处理接口返回的响应处理方法，它也正与错误码标准化是相对应的

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}
func (s StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(s.String(), 10, 64)
	return v, err
}
func (s StrTo) Int32() (int32, error) {
	v, err := strconv.ParseInt(s.String(), 10, 32)
	return int32(v), err
}
func (s StrTo) MustInt64() int64 {
	v, _ := s.Int64()
	return v
}
func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

func (s StrTo) MustInt32() int32 {
	v, _ := s.Int32()
	return v
}

func IDToSting(id int64) string {
	return strconv.FormatInt(id, 10)
}

func StringToIDMust(id string) int64 {
	res, _ := strconv.ParseInt(id, 10, 64)
	return res
}

func LinkStr(a, b string) string {
	return a + ":" + b
}

func LinkID(a, b int64) string {
	return IDToSting(a) + ":" + IDToSting(b)
}

func ParseLinkID(str string) (a, b int64) {
	result := strings.Split(str, ":")
	if len(result) != 2 {
		return
	}
	return StringToIDMust(result[0]), StringToIDMust(result[1])
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func StringToBoolMust(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
