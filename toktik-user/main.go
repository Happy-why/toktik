package main

import (
	"fmt"
	"toktik-user/internal/global"
	"toktik-user/pkg/rpc"
	"toktik-user/pkg/rpc/client"
	"toktik-user/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)
	client.NewRpcClient()
	fmt.Println("-----User Service Start ! ! !-----")
	// RPC 注册
	rpc.RegisterRPC()

}
