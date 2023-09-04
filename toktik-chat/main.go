package main

import (
	"fmt"
	"toktik-chat/internal/global"
	"toktik-chat/pkg/rpc"
	"toktik-chat/pkg/rpc/client"
	"toktik-chat/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)

	// RPC 注册
	client.NewRpcClient()
	fmt.Println("-----Chat Service Start ! ! !-----")
	rpc.RegisterRPC()
}
