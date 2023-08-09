package user

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (h *HandlerUser) getCaptcha(c *gin.Context) {
	req := &user.RegisterRequest{Username: "65456456asdsa"}
	resp, err := userClient.Register(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	c.JSON(200, "getCaptcha success")
	fmt.Println("hello")
	return
}
