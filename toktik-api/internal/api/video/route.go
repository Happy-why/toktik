package video

import (
	"github.com/Happy-Why/toktik-api/pkg/middleware"
	"github.com/Happy-Why/toktik-api/pkg/router"
	"github.com/gin-gonic/gin"
	"log"
)

type RouterVideo struct {
}

func init() {
	log.Println("init Interaction router success")
	rv := &RouterVideo{}
	router.Register(rv)
}

func (*RouterVideo) Route(r *gin.Engine) {
	InitRpcInteractionClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	r.GET("/douyin/feed/")
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/publish/action/", h.VideoFeed)
		g.GET("/publish/list/")
		g.POST("/favorite/action/")
		g.GET("/favorite/list/")
		g.POST("/comment/action/")
		g.GET("/comment/list/")
	}
}
