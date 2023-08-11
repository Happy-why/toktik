package main

import (
	"fmt"
	srv "github.com/Happy-Why/toktik-common/serveHTTP"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/Happy-Why/toktik-user/pkg/router"
	"github.com/Happy-Why/toktik-user/pkg/setting"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化
	setting.InitAllSetting()
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
	srv.Run(route, global.PbSettings.Server.Name, global.PbSettings.Server.Addr, stop)
}
