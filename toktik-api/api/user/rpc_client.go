package user

import (
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"log"
)

var userClient userservice.Client

func InitRpcUserClient() {
	c, err := userservice.NewClient("user", client.WithHostPorts("127.0.0.1:8881"))
	if err != nil {
		log.Fatal(err)
	}
	userClient = c
}
