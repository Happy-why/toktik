package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	srv "toktik-common/serveHTTP"
	"toktik-video/internal/global"
	"toktik-video/pkg/rpc"
	"toktik-video/pkg/rpc/client"
	"toktik-video/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)
	// 初始化 gin
	route := gin.Default()

	// RPC 注册
	kr := rpc.RegisterRPC()
	client.NewRpcClient()
	//服务端配置
	stop := func() { kr.Stop() }
	fmt.Println("------------------------------------------------")
	srv.Run(route, global.Settings.Server.Name, global.Settings.Server.Addr, stop)
}
