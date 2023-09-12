package rpc

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"net"
	"toktik-rpc/kitex_gen/video/videoservice"
	"toktik-video/internal/global"
	"toktik-video/internal/model"
	"toktik-video/internal/service"
)

func RegisterRPC() server.Server {
	r, err := etcd.NewEtcdRegistry(global.Settings.Etcd.Addr)
	if err != nil {
		zap.L().Error("etcd.NewEtcdRegistry err:", zap.Error(err))
		return nil
	}
	addr, err := net.ResolveTCPAddr("tcp", global.Settings.Rpc.Addr) //:8883
	if err != nil {
		zap.L().Error("net.ResolveTCPAddr err:", zap.Error(err))
		return nil
	}

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.Settings.Jaeger.ServerName[model.TokTikVideo]),
		provider.WithExportEndpoint(global.Settings.Jaeger.RPCExportEndpoint),
		provider.WithInsecure(),
	)
	svr := videoservice.NewServer(
		service.NewVideoService(),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.Settings.Rpc.Name}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithTracer(prometheus.NewServerTracer(global.Settings.Prometheus.Post, global.Settings.Prometheus.Path)),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
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
