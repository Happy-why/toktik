#!/bin/bash

echo "编译Linux版本64位"
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 编译各个服务的主程序
go build -o toktik-api/bin/toktik-api toktik-api/main.go
go build -o toktik-user/bin/toktik-user toktik-user/main.go
go build -o toktik-interaction/bin/toktik-interaction toktik-interaction/main.go
go build -o toktik-video/bin/toktik-video toktik-video/main.go
go build -o toktik-chat/bin/toktik-chat toktik-chat/main.go
go build -o toktik-favor/bin/toktik-favor toktik-favor/main.go
go build -o toktik-comment/bin/toktik-comment toktik-comment/main.go

echo "可执行文件编译完成"
