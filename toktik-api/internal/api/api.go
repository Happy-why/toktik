package api

import (
	//_ "github.com/Happy-Why/toktik-api/internal/api/user"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user/userservice"
)

var (
	UserClient userservice.Client
)
