package api

import (
	//_ "github.com/Happy-Why/toktik-api/internal/api/user"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/interaction/interactionservice"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user/userservice"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/video/videoservice"
)

var (
	UserClient        userservice.Client
	InteractionClient interactionservice.Client
	VideoClient       videoservice.Client
)
