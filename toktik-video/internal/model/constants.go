package model

const (
	MsgNil = ""
)

type FileType string

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
	RpcUser        = "rpc-user"
	RpcVideo       = "rpc-video"
	RpcInteraction = "rpc-interaction"
	RpcChat        = "rpc-chat"
)
