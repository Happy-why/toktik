package user

import (
	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-api/internal/api"
	"toktik-api/internal/global"
	"toktik-api/internal/model"
	"toktik-rpc/kitex_gen/user/userservice"
)

func InitRpcUserClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.Settings.Jaeger.ServerName[model.TokTikUser]),
		provider.WithExportEndpoint(global.Settings.Jaeger.RPCExportEndpoint),
		provider.WithInsecure(),
	)

	c, err := userservice.NewClient(
		model.RpcUser,
		client.WithSuite(tracing.NewClientSuite()),
		//client.WithHostPorts(global.Settings.Rpc.ServerAddrs[model.RpcUser]),
		//client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		//client.WithInstanceMW(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		zap.L().Error("apiServer InitRpc User Client err:", zap.Error(err))
		panic(err)
	}
	api.UserClient = c
}
