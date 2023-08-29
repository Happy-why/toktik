package api

import (
	"toktik-rpc/kitex_gen/chat/chatservice"
	//_ "toktik-api/internal/api/user"
	"toktik-rpc/kitex_gen/interaction/interactionservice"
	"toktik-rpc/kitex_gen/user/userservice"
	"toktik-rpc/kitex_gen/video/videoservice"
)

var (
	UserClient        userservice.Client
	InteractionClient interactionservice.Client
	VideoClient       videoservice.Client
	ChatClient        chatservice.Client
)
