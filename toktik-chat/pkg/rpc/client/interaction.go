package client

import (
	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-chat/internal/global"
	"toktik-chat/internal/model"
	"toktik-rpc/kitex_gen/interaction/interactionservice"
)

var InteractionClient interactionservice.Client

func InitRpcInteractionClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.Settings.Jaeger.ServerName[model.TokTikInteraction]),
		provider.WithExportEndpoint(global.Settings.Jaeger.RPCExportEndpoint),
		provider.WithInsecure(),
	)

	c, err := interactionservice.NewClient(
		model.RpcInteraction,
		client.WithSuite(tracing.NewClientSuite()),
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
