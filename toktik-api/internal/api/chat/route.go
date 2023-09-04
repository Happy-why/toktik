package chat

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"toktik-api/pkg/router"
)

type RouterChat struct {
}

func init() {
	log.Println("init Interaction router success")
	rc := &RouterChat{}
	router.Register(rc)
}

func (*RouterChat) Route(r *server.Hertz) {
	InitRpcChatClient()
	//初始化grpc的客户端连接
	h := NewHandlerInteraction()
	g := r.Group("/douyin/message" /*middleware.MustUser()*/)
	{
		g.GET("/chat/", h.MessageList)
		g.POST("/action/", h.ChatAction)

	}
}
