package video

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"toktik-api/pkg/router"
)

type RouterComment struct {
}

func init() {
	log.Println("init Comment router success")
	rc := &RouterComment{}
	router.Register(rc)
}

func (*RouterComment) Route(r *server.Hertz) {
	InitRpcCommentClient()
	//初始化grpc的客户端连接
	h := NewHandlerVideo()
	r.GET("/douyin/comment/list/", h.CommentList)
	g := r.Group("/douyin" /* middleware.MustUser()*/)
	{
		g.POST("/comment/action/", h.CommentAction)
	}
}
