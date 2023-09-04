package main

import (
	"fmt"
	"toktik-favor/internal/global"
	"toktik-favor/pkg/rpc"
	"toktik-favor/pkg/rpc/client"
	"toktik-favor/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)

	// RPC 注册
	client.NewRpcClient()
	fmt.Println("-----Favor Service Start ! ! !-----")
	rpc.RegisterRPC()

}
