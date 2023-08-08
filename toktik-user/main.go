package main

import (
	"fmt"
	srv "github.com/Happy-Why/toktik-common/ServeHTTP"
	_ "github.com/Happy-Why/toktik-user/internal/api/user"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/Happy-Why/toktik-user/pkg/router"
	"github.com/Happy-Why/toktik-user/pkg/setting"
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
	kr := router.RegisterRPC()
	//服务端配置
	stop := func() { kr.Stop() }
	fmt.Println("------------------------------------------------")
	srv.Run(route, "user", "127.0.0.1:8081", stop)
}
