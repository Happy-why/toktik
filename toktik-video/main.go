package main

import (
	"fmt"
	"toktik-video/internal/global"
	"toktik-video/pkg/rpc"
	"toktik-video/pkg/rpc/client"
	"toktik-video/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)

	// RPC 客户端注册
	client.NewRpcClient()
	fmt.Println("-----Video Service Start ! ! !-----")
	// RPC 服务注册
	rpc.RegisterRPC()
}
