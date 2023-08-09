package router

import (
	"github.com/gin-gonic/gin"
)

var routers []Router

type Router interface {
	Route(r *gin.Engine)
}

func Register(root ...Router) {
	routers = append(routers, root...)
}

func InitRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})
	for _, root := range routers {
		root.Route(r)
	}
}

//func RegisterRPC() server.Server {
//	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8881")
//	if err != nil {
//		panic(err)
//	}
//
//	//svr := user.NewServer(
//	//	new(service.UserServiceImpl),
//	//	server.WithServiceAddr(addr),
//	//)
//	go func() {
//		err = svr.Run()
//		if err != nil {
//			panic(err)
//		}
//	}()
//	return svr
//}
