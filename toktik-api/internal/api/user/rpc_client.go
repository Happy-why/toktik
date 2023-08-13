package user

import (
	"github.com/Happy-Why/toktik-api/internal/api"
	"github.com/Happy-Why/toktik-api/internal/global"
	"github.com/Happy-Why/toktik-api/internal/model"
	rpcmiddleware "github.com/Happy-Why/toktik-common/rpc-middleware"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
)

func InitRpcUserClient() {
	r, err := etcd.NewEtcdResolver(global.PvSettings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := userservice.NewClient(
		model.RpcUser,
		client.WithHostPorts(global.PbSettings.Rpc.ServerAddrs[model.RpcUser]),
		client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		client.WithInstanceMW(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.UserClient = c
}
