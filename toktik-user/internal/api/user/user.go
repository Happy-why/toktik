package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (h *HandlerUser) getCaptcha(c *gin.Context) {
	c.JSON(200, "getCaptcha success")
	fmt.Println("hello")
	return
}
