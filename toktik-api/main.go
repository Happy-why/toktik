package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	_ "toktik-api/internal/api/chat"
	_ "toktik-api/internal/api/interaction"
	_ "toktik-api/internal/api/user"
	_ "toktik-api/internal/api/video"
	"toktik-api/internal/global"
	"toktik-api/pkg/middleware"
	"toktik-api/pkg/router"
	"toktik-api/pkg/setting"
	"toktik-api/pkg/tracing"
	srv "toktik-common/serveHTTP"
)

func main() {
	// 初始化
	setting.InitAllSetting()
	fmt.Printf("config:%#v\n", global.Settings)
	// 加载 jaeger
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// 初始化 gin
	route := gin.Default()
	route.Use(middleware.Auth(), middleware.Cors(), otelgin.Middleware("toktik-api"))
	pprof.Register(route)
	// 路由注册
	router.InitRouter(route)

	// RPC 注册
	//kr := router.RegisterRPC()
	//服务端配置
	//stop := func() { kr.Stop() }
	fmt.Println("------------------------------------------------")
	srv.Run(route, global.Settings.Server.Name, global.Settings.Server.Addr, nil)
}
