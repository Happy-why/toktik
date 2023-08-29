package interaction

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
