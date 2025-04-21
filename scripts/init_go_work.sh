#!/bin/bash

# 初始化 go.work 文件
go work init

# 添加服务模块
go work use ./toktik-api
go work use ./toktik-chat
go work use ./toktik-comment
go work use ./toktik-common
go work use ./toktik-favor
go work use ./toktik-interaction
go work use ./toktik-rpc
go work use ./toktik-user
go work use ./toktik-video
