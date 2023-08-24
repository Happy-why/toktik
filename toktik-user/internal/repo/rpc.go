package repo

import "github.com/Happy-Why/toktik-rpc/kitex_gen/interaction/interactionservice"

type RpcRepo interface {
	interactionservice.Client
}
