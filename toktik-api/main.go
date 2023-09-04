package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	_ "toktik-api/internal/api/chat"
	_ "toktik-api/internal/api/comment"
	_ "toktik-api/internal/api/favor"
	_ "toktik-api/internal/api/interaction"
	_ "toktik-api/internal/api/user"
	_ "toktik-api/internal/api/video"
	"toktik-api/internal/global"
	"toktik-api/internal/model"
	"toktik-api/pkg/middleware"
	"toktik-api/pkg/router"
	"toktik-api/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)

	// 注册链路追踪
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.Settings.Jaeger.ServerName[model.TokTikApi]),
		provider.WithExportEndpoint(global.Settings.Jaeger.RPCExportEndpoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	tracer, cfg := tracing.NewServerTracer()
	// 初始化 hz
	hz := server.Default(
		server.WithHostPorts(global.Settings.Server.Addr),
		//server.WithMaxRequestBodySize(),
		tracer,
	)
	// 注册 中间件
	hz.Use(gzip.Gzip(gzip.DefaultCompression))
	hz.Use(middleware.Auth())
	hz.Use(tracing.ServerMiddleware(cfg))
	// 路由注册
	router.InitRouter(hz)

	fmt.Println("-----API Server Start ! ! !-----")
	hz.Spin()
}
