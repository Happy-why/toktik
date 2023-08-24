package client

//func NewRpcClientRepo() *RpcClient {
//	return &RpcClient{InteractionClient,VideoClient}
//}

//type RpcClient struct {
//	 interactionservice.Client
//	 videoservice.Client
//}

func NewRpcClient() {
	InitRpcInteractionClient()
}
