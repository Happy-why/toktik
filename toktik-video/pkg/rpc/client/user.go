package client

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-rpc/kitex_gen/user/userservice"
	"toktik-video/internal/global"
	"toktik-video/internal/model"
)

var UserClient userservice.Client

func InitRpcUserClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := userservice.NewClient(
		model.RpcUser,
		//client.WithHostPorts(global.Settings.Rpc.ServerAddrs[model.RpcUser]),
		//client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		//client.WithInstanceMW(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		zap.L().Error("InitRpcUserClient err:", zap.Error(err))
		panic(err)
	}
	UserClient = c
}
