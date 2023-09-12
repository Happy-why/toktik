package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/monitor-prometheus"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
	log2 "log"
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
	"toktik-api/pkg/sentinel"
	"toktik-api/pkg/setting"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	log2.Printf("config:%#v\n", global.Settings)

	// 初始化sentinel
	sentinel.InitSentinel()

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
		server.WithTracer(prometheus.NewServerTracer(global.Settings.Prometheus.Post, global.Settings.Prometheus.Path)),
		//server.WithMaxRequestBodySize(),
		tracer,
	)
	// 注册 中间件
	hz.Use(gzip.Gzip(gzip.DefaultCompression))
	hz.Use(middleware.Auth())
	hz.Use(tracing.ServerMiddleware(cfg))
	// hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
	hz.Use(hertzSentinel.SentinelServerMiddleware(
		hertzSentinel.WithServerResourceExtractor(func(c context.Context, ctx *app.RequestContext) string {
			return model.SentinelApi
		}),
		hertzSentinel.WithServerBlockFallback(func(c context.Context, ctx *app.RequestContext) {
			ctx.AbortWithStatusJSON(400, utils.H{
				"err":  "too many request; the quota used up",
				"code": 10222,
			})
		})))
	hz.Use(middleware.LimitIP())
	// 路由注册
	router.InitRouter(hz)

	log2.Println("-----API Server Start ! ! !-----")
	hz.Spin()
}
