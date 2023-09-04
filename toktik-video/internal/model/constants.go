package model

const (
	MsgNil = ""
)

type FileType string

const (
	RpcSuccess = 0
)

const (
	MP4 = ".mp4"
	JPG = ".jpg"
)

const (
	VIDEO = "video"
	COVER = "cover"
)

const (
	FAVORITE = iota + 1
	CANCELFAVORITE
)

const (
	COMMENT = iota + 1
	CANCELCOMMENT
)

const (
	RpcUser        = "rpc-user"
	RpcVideo       = "rpc-video"
	RpcInteraction = "rpc-interaction"
	RpcChat        = "rpc-chat"
	RpcFavor       = "rpc-favor"
	RpcComment     = "rpc-comment"
)

const (
	TokTikApi         = "toktik-api"
	TokTikUser        = "toktik-user"
	TokTikInteraction = "toktik-interaction"
	TokTikVideo       = "toktik-video"
	TokTikChat        = "toktik-chat"
	TokTikFavor       = "toktik-favor"
	TokTikComment     = "toktik-comment"
)
