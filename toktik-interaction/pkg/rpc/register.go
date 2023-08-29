package rpc

import (
	"toktik-interaction/internal/global"
	"toktik-interaction/internal/service"
	"toktik-rpc/kitex_gen/interaction/interactionservice"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"net"
)

func RegisterRPC() server.Server {
	r, err := etcd.NewEtcdRegistry(global.Settings.Etcd.Addr)
	if err != nil {
		zap.L().Error("etcd.NewEtcdRegistry err:", zap.Error(err))
		return nil
	}
	addr, err := net.ResolveTCPAddr("tcp", global.Settings.Rpc.Addr) //:8882
	if err != nil {
		zap.L().Error("net.ResolveTCPAddr err:", zap.Error(err))
		return nil
	}
	svr := interactionservice.NewServer(
		service.NewInteractionService(),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.Settings.Rpc.Name}),
		//server.WithMiddleware(rpcmiddleware.CommonMiddleware), // middleware
		//server.WithMiddleware(rpcmiddleware.ServerMiddleware),
		server.WithRegistry(r),
	)
	go func() {
		err := svr.Run()
		if err != nil {
			zap.L().Error("RPC svr.Run() err:", zap.Error(err))
			panic(err)
		}
	}()
	return svr
}
