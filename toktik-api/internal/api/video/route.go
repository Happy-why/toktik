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
	r.GET("/douyin/feed/", h.VideoFeed)
	g := r.Group("/douyin", middleware.MustUser())
	{
		g.POST("/publish/action/", h.VideoPublish)
		g.GET("/publish/list/", h.PublishList)
		g.POST("/favorite/action/", h.FavoriteAction)
		g.GET("/favorite/list/", h.FavoriteList)
		g.POST("/comment/action/", h.CommentAction)
		g.GET("/comment/list/", h.CommentList)
	}
}
