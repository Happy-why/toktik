package request

import (
	"github.com/Happy-Why/toktik-api/pkg/myerr"
	"github.com/Happy-Why/toktik-common/errcode"
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
