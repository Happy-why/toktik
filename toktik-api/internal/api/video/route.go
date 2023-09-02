package video

import (
	"github.com/gin-gonic/gin"
	"log"
	"toktik-api/pkg/middleware"
	"toktik-api/pkg/router"
)

type RouterVideo struct {
}

func init() {
	log.Println("init video router success")
	rv := &RouterVideo{}
	router.Register(rv)
}

func (*RouterVideo) Route(r *gin.Engine) {
	InitRpcVideoClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	r.GET("/douyin/feed/", h.VideoFeed)
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/publish/action/", h.VideoPublish)
		g.GET("/publish/list/", h.PublishList)
	}
}
