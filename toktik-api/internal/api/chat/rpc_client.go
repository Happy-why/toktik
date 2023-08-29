package chat

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"toktik-api/internal/api"
	"toktik-api/internal/global"
	"toktik-api/internal/model"
	"toktik-rpc/kitex_gen/chat/chatservice"
)

func InitRpcChatClient() {
	r, err := etcd.NewEtcdResolver(global.Settings.Etcd.Addr)
	if err != nil {
		panic(err)
	}
	c, err := chatservice.NewClient(
		model.RpcChat,
		//client.WithHostPorts(global.Settings.Rpc.ServerAddrs[model.RpcInteraction]),
		//client.WithMiddleware(rpcmiddleware.CommonMiddleware),
		//client.WithInstanceMW(rpcmiddleware.ClientMiddleware),
		client.WithResolver(r),
	)
	if err != nil {
		zap.L().Error("apiServer InitRpcChatClient err:", zap.Error(err))
		panic(err)
	}
	api.ChatClient = c
}
