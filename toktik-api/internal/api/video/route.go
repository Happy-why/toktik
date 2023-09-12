package video

import (
	"github.com/cloudwego/hertz/pkg/app/server"
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

func (*RouterVideo) Route(r *server.Hertz) {
	InitRpcVideoClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	r.GET("/douyin/feed/", h.VideoFeed)
	r.GET("/douyin//publish/list/", h.PublishList)
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/publish/action/", h.VideoPublish)

	}
}
