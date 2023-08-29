package client

import (
	"toktik-user/internal/global"
	"toktik-user/internal/model"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-rpc/kitex_gen/interaction/interactionservice"
)

var InteractionClient interactionservice.Client

func InitRpcInteractionClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := interactionservice.NewClient(
		model.RpcInteraction,
		//client.WithHostPorts(global.PbSettings.Rpc.ServerAddrs[model.RpcInteraction]),
		//client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		//client.WithMiddleware(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		zap.L().Error("InitRpcInteractionClient err:", zap.Error(err))
		panic(err)
	}
	InteractionClient = c
}
