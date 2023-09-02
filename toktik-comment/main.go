package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"toktik-comment/internal/global"
	"toktik-comment/pkg/rpc"
	"toktik-comment/pkg/rpc/client"
	"toktik-comment/pkg/setting"
	srv "toktik-common/serveHTTP"
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
