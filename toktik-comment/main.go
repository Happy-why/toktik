package main

import (
	"fmt"
	"toktik-comment/internal/global"
	"toktik-comment/pkg/rpc"
	"toktik-comment/pkg/rpc/client"
	"toktik-comment/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)

	// RPC 注册
	client.NewRpcClient()
	fmt.Println("-----Comment Service Start ! ! !-----")
	rpc.RegisterRPC()
}
