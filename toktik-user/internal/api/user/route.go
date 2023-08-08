package user

import (
	"github.com/Happy-Why/toktik-user/pkg/router"
	"github.com/gin-gonic/gin"
	"log"
)

type RouterUser struct {
}

func init() {
	log.Println("init User router success")
	ru := &RouterUser{}
	router.Register(ru)
}

func (*RouterUser) Route(r *gin.Engine) {
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)

}
