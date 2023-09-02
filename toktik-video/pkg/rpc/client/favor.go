package client

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-rpc/kitex_gen/favor/favorservice"
	"toktik-video/internal/global"
	"toktik-video/internal/model"
)

var FavorClient favorservice.Client

func InitRpcFavorCache() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := favorservice.NewClient(
		model.RpcFavor,
		//client.WithHostPorts(global.Settings.Rpc.ServerAddrs[model.RpcUser]),
		//client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		//client.WithInstanceMW(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		zap.L().Error("InitRpcFavorCache err:", zap.Error(err))
		panic(err)
	}
	FavorClient = c
}
