package chat

import (
	"github.com/gin-gonic/gin"
	"log"
	"toktik-api/pkg/middleware"
	"toktik-api/pkg/router"
)

type RouterInteraction struct {
}

func init() {
	log.Println("init Interaction router success")
	ri := &RouterInteraction{}
	router.Register(ri)
}

func (*RouterInteraction) Route(r *gin.Engine) {
	InitRpcChatClient()
	//初始化grpc的客户端连接
	h := NewHandlerInteraction()
	g := r.Group("/douyin/message", middleware.MustUser())
	{
		g.GET("/chat/", h.MessageList)
		g.POST("/action/", h.ChatAction)

	}
}
