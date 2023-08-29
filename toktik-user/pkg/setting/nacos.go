package setting

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	log2 "log"
	"toktik-user/internal/global"
)

type NacosClient struct {
	confClient config_client.IConfigClient
}

func InitNacos() *NacosClient {
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.Settings.Nacos.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.Settings.Nacos.Addr,
			ContextPath: global.Settings.Nacos.ContextPath,
			Port:        uint64(global.Settings.Nacos.Port),
			Scheme:      global.Settings.Nacos.Scheme,
		},
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log2.Fatalln(err)
	}
	nc := &NacosClient{
		confClient: configClient,
	}
	return nc
}
