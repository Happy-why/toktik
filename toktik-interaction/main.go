package main

import (
	"fmt"
	srv "github.com/Happy-Why/toktik-common/serveHTTP"
	"github.com/Happy-Why/toktik-interaction/internal/global"
	"github.com/Happy-Why/toktik-interaction/pkg/rpc"
	"github.com/Happy-Why/toktik-interaction/pkg/rpc/client"
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

	// RPC 注册
	kr := rpc.RegisterRPC()
	client.NewRpcClient()
	//服务端配置
	stop := func() { kr.Stop() }
	fmt.Println("------------------------------------------------")
	srv.Run(route, global.PbSettings.Server.Name, global.PbSettings.Server.Addr, stop)
}
