package client

import (
	"github.com/Happy-Why/toktik-rpc/kitex_gen/interaction/interactionservice"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/video/videoservice"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/Happy-Why/toktik-user/internal/model"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

var VideoClient videoservice.Client

func InitRpcVideoClient() {
	r, err := etcd.NewEtcdResolver(global.PvSettings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := interactionservice.NewClient(
		model.RpcInteraction,
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
