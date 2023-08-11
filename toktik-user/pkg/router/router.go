package router

import (
	rpcmiddleware "github.com/Happy-Why/toktik-common/rpc-middleware"
	user "github.com/Happy-Why/toktik-rpc/kitex_gen/user/userservice"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/Happy-Why/toktik-user/internal/service"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/gin-gonic/gin"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"net"
)

var routers []Router

type Router interface {
	Route(r *gin.Engine)
}

func Register(root ...Router) {
	routers = append(routers, root...)
}

func InitRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})
	for _, root := range routers {
		root.Route(r)
	}
}

func RegisterRPC() server.Server {
	r, err := etcd.NewEtcdRegistry(global.PvSettings.Etcd.Addr)
	if err != nil {
		zap.L().Error("etcd.NewEtcdRegistry err:", zap.Error(err))
		return nil
	}
	addr, err := net.ResolveTCPAddr("tcp", global.PbSettings.Rpc.Addr) //:8881
	if err != nil {
		zap.L().Error("net.ResolveTCPAddr err:", zap.Error(err))
		return nil
	}
	svr := user.NewServer(
		service.NewUserService(),
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
