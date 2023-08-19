package main

import (
	"fmt"
	"github.com/Happy-Why/toktik-api/pkg/middleware"
	srv "github.com/Happy-Why/toktik-common/serveHTTP"
	"github.com/Happy-Why/toktik-interaction/internal/global"
	"github.com/Happy-Why/toktik-interaction/pkg/router"
	"github.com/Happy-Why/toktik-interaction/pkg/setting"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.PbSettings)
	fmt.Printf("config:%#v\n", global.PvSettings)
	// 初始化 gin
	route := gin.Default()
	route.Use(middleware.Auth())

	// RPC 注册
	kr := router.RegisterRPC()
	//服务端配置
	stop := func() { kr.Stop() }
	fmt.Println("------------------------------------------------")
	srv.Run(route, global.PbSettings.Server.Name, global.PbSettings.Server.Addr, stop)
}
