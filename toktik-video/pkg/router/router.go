package router

import (
	rpcmiddleware "github.com/Happy-Why/toktik-common/rpc-middleware"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/video/videoservice"
	"github.com/Happy-Why/toktik-video/internal/global"
	"github.com/Happy-Why/toktik-video/internal/service"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"net"
)

func RegisterRPC() server.Server {
	r, err := etcd.NewEtcdRegistry(global.PvSettings.Etcd.Addr)
	if err != nil {
		zap.L().Error("etcd.NewEtcdRegistry err:", zap.Error(err))
		return nil
	}
	addr, err := net.ResolveTCPAddr("tcp", global.PbSettings.Rpc.Addr) //:8882
	if err != nil {
		zap.L().Error("net.ResolveTCPAddr err:", zap.Error(err))
		return nil
	}
	svr := videoservice.NewServer(
		service.NewVideoService(),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.PbSettings.Rpc.Name}),
		server.WithMiddleware(rpcmiddleware.CommonMiddleware), // middleware
		server.WithMiddleware(rpcmiddleware.ServerMiddleware),
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
