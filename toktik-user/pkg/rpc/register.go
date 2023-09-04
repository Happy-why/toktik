package rpc

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"net"
	"toktik-rpc/kitex_gen/user/userservice"
	"toktik-user/internal/global"
	"toktik-user/internal/model"
	"toktik-user/internal/service"
)

func RegisterRPC() server.Server {
	r, err := etcd.NewEtcdRegistry(global.Settings.Etcd.Addr)
	if err != nil {
		zap.L().Error("etcd.NewEtcdRegistry err:", zap.Error(err))
		return nil
	}
	addr, err := net.ResolveTCPAddr("tcp", global.Settings.Rpc.Addr) //:8881
	if err != nil {
		zap.L().Error("net.ResolveTCPAddr err:", zap.Error(err))
		return nil
	}

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.Settings.Jaeger.ServerName[model.TokTikUser]),
		provider.WithExportEndpoint(global.Settings.Jaeger.RPCExportEndpoint),
		provider.WithInsecure(),
	)

	svr := userservice.NewServer(
		service.NewUserService(),
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
