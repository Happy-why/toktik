package video

import (
	"github.com/gin-gonic/gin"
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

func (*RouterFavor) Route(r *gin.Engine) {
	InitRpcFavorClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/favorite/action/", h.FavoriteAction)
		g.GET("/favorite/list/", h.FavoriteList)
	}
}
