package main

import (
	"fmt"
	_ "github.com/Happy-Why/toktik-api/api"
	"github.com/Happy-Why/toktik-api/internal/global"
	"github.com/Happy-Why/toktik-api/pkg/router"
	"github.com/Happy-Why/toktik-api/pkg/setting"
	srv "github.com/Happy-Why/toktik-common/serveHTTP"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化
	setting.InitAll()
	fmt.Printf("config:%#v\n", global.PbSettings)
	fmt.Printf("config:%#v\n", global.PvSettings)
	// 初始化 gin
	route := gin.Default()
	// 路由注册
	router.InitRouter(route)

	// RPC 注册
	//kr := router.RegisterRPC()
	//服务端配置
	//stop := func() { kr.Stop() }
	fmt.Println("------------------------------------------------")
	srv.Run(route, "user", "127.0.0.1:8082", nil)
}
