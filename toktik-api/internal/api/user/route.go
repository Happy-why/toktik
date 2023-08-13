package user

import (
	"github.com/Happy-Why/toktik-api/pkg/middleware"
	"github.com/Happy-Why/toktik-api/pkg/router"
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
	InitRpcUserClient()
	//初始化grpc的客户端连接
	h := NewHandlerUser()
	r.POST("/douyin/user/register/", h.Register)
	r.POST("/douyin/user/login/", h.Login)

	g := r.Group("", middleware.MustUser())
	{
		g.GET("/douyin/user/", h.UserIndex)
	}

}
