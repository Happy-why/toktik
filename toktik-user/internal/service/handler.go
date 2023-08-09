package service

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	fmt.Printf("req:%#v\n", req)
	resp = &user.RegisterResponse{UserId: 99988, Token: "woshi"}
	return
}
