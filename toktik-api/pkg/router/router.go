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
