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
	log.Println("init Interaction router success")
	rv := &RouterVideo{}
	router.Register(rv)
}

func (*RouterVideo) Route(r *gin.Engine) {
	InitRpcInteractionClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	r.GET("/douyin/feed/", h.VideoFeed)
	r.GET("/comment/list/", h.CommentList)
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/publish/action/", h.VideoPublish)
		g.GET("/publish/list/", h.PublishList)
		g.POST("/favorite/action/", h.FavoriteAction)
		g.GET("/favorite/list/", h.FavoriteList)
		g.POST("/comment/action/", h.CommentAction)

	}
}
