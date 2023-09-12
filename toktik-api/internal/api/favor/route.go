package video

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"toktik-api/pkg/middleware"
	"toktik-api/pkg/router"
)

type RouterFavor struct {
}

func init() {
	log.Println("init Favor router success")
	rf := &RouterFavor{}
	router.Register(rf)
}

func (*RouterFavor) Route(r *server.Hertz) {
	InitRpcFavorClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	r.GET("/douyin/favorite/list/", h.FavoriteList)
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/favorite/action/", h.FavoriteAction)

	}
}
