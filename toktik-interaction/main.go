package main

import (
	"fmt"
	"toktik-interaction/internal/global"
	"toktik-interaction/pkg/rpc"
	"toktik-interaction/pkg/rpc/client"
	"toktik-interaction/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)

	// RPC 注册
	client.NewRpcClient()
	fmt.Println("-----Interaction Service Start ! ! !-----")
	rpc.RegisterRPC()
}
