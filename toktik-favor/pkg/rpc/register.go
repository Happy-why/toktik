package rpc

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"net"
	"toktik-favor/internal/global"
	"toktik-favor/internal/model"
	"toktik-favor/internal/service"
	"toktik-rpc/kitex_gen/favor/favorservice"
)

func RegisterRPC() server.Server {
	r, err := etcd.NewEtcdRegistry(global.Settings.Etcd.Addr)
	if err != nil {
		zap.L().Error("etcd.NewEtcdRegistry err:", zap.Error(err))
		return nil
	}
	addr, err := net.ResolveTCPAddr("tcp", global.Settings.Rpc.Addr) //:8886
	if err != nil {
		zap.L().Error("net.ResolveTCPAddr err:", zap.Error(err))
		return nil
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.Settings.Jaeger.ServerName[model.TokTikFavor]),
		provider.WithExportEndpoint(global.Settings.Jaeger.RPCExportEndpoint),
		provider.WithInsecure(),
	)
	svr := favorservice.NewServer(
		service.NewFavorService(),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.Settings.Rpc.Name}),
		server.WithSuite(tracing.NewServerSuite()),
		//server.WithMiddleware(rpcmiddleware.CommonMiddleware), // middleware
		//server.WithMiddleware(rpcmiddleware.ServerMiddleware),
		server.WithRegistry(r),
	)
	err = svr.Run()
	if err != nil {
		zap.L().Error("RPC svr.Run() err:", zap.Error(err))
		panic(err)
	}
	return svr

}
