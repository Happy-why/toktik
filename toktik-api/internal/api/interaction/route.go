package interaction

import (
	"github.com/cloudwego/hertz/pkg/app/server"
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

func (*RouterInteraction) Route(r *server.Hertz) {
	InitRpcInteractionClient()
	//初始化grpc的客户端连接
	h := NewHandlerInteraction()
	g := r.Group("/douyin/relation", middleware.MustUser())
	{
		g.POST("/action/", h.FollowSB)
		g.GET("/follow/list/", h.FollowList)
		g.GET("/follower/list/", h.FansList)
		g.GET("/friend/list/", h.FriendList)
	}
}
