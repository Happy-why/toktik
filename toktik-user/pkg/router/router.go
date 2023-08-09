package router

import (
	"context"
	user "github.com/Happy-Why/toktik-rpc/kitex_gen/user/userservice"
	"github.com/Happy-Why/toktik-user/internal/service"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/gin-gonic/gin"
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
	//r, err := etcd.NewEtcdRegistry([]string{"654"})
	//if err != nil {
	//	panic(err)
	//}
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8881")
	if err != nil {
		panic(err)
	}

	svr := user.NewServer(
		new(service.UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "user"}),
		server.WithMiddleware(CommonMiddleware), // middleware
		server.WithMiddleware(ServerMiddleware),
	)
	go func() {
		err = svr.Run()
		if err != nil {
			panic(err)
		}
	}()
	return svr
}

// CommonMiddleware common middleware print some rpc info、real request and real response
func CommonMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)
		// get real request
		klog.Infof("real request: %+v\n", req)
		// get remote service information
		klog.Infof("remote service name: %s, remote method: %s\n", ri.To().ServiceName(), ri.To().Method())
		if err = next(ctx, req, resp); err != nil {
			return err
		}
		// get real response
		klog.Infof("real response: %+v\n", resp)
		return nil
	}
}

// ServerMiddleware server middleware print client address
func ServerMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)
		// get client information
		klog.Infof("client address: %v\n", ri.From().Address())
		if err = next(ctx, req, resp); err != nil {
			return err
		}
		return nil
	}
}
