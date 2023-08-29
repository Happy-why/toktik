package client

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-rpc/kitex_gen/interaction/interactionservice"
	"toktik-video/internal/global"
	"toktik-video/internal/model"
)

var InteractionClient interactionservice.Client

func InitRpcInteractionClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := interactionservice.NewClient(
		model.RpcInteraction,
		//client.WithHostPorts(global.Settings.Rpc.ServerAddrs[model.RpcUser]),
		//client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		//client.WithInstanceMW(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		zap.L().Error("InitRpcInteractionClient err:", zap.Error(err))
		panic(err)
	}
	InteractionClient = c
}
