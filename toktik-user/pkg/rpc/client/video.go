package client

import (
	"toktik-rpc/kitex_gen/interaction/interactionservice"
	"toktik-rpc/kitex_gen/video/videoservice"
	"toktik-user/internal/global"
	"toktik-user/internal/model"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

var VideoClient videoservice.Client

func InitRpcVideoClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := interactionservice.NewClient(
		model.RpcVideo,
		//client.WithHostPorts(global.PbSettings.Rpc.ServerAddrs[model.RpcInteraction]),
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
