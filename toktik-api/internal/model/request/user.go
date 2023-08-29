package request

import (
	"toktik-api/pkg/myerr"
	"toktik-common/errcode"
)

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (r RegisterRequest) Verify() errcode.Err {
	if r.Username == "" || len(r.Username) > 32 || r.Password == "" || len(r.Password) > 32 {
		return myerr.ErrUserNameORPassWord
	}
	return nil
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (r LoginRequest) Verify() errcode.Err {
	if r.Username == "" || len(r.Username) > 32 || r.Password == "" || len(r.Password) > 32 {
		return myerr.ErrUserNameORPassWord
	}
	return nil
}

type UserIndexRequest struct {
	UserId   int64  `json:"user_id" form:"user_id"`
	Token    string `json:"token" form:"token"`
	MyUserId int64  `json:"my_user_id" form:"my_user_id"`
}
