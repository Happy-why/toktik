chcp 65001
@echo off
:loop
@echo off&amp;color 0A
cls
echo,
echo 请选择要编译的系统环境：
echo,
echo 1. Windows_amd64
echo 2. linux_amd64

set/p action=请选择:
if %action% == 1 goto build_Windows_amd64
if %action% == 2 goto build_linux_amd64

:build_Windows_amd64
echo 编译Windows版本64位
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o toktik-api/bin/toktik-api toktik-api/main.go
go build -o toktik-user/bin/toktik-user toktik-user/main.go
go build -o toktik-interaction/bin/toktik-interaction toktik-interaction/main.go
go build -o toktik-video/bin/toktik-video toktik-video/main.go
go build -o toktik-chat/bin/toktik-chat toktik-chat/main.go
go build -o toktik-favor/bin/toktik-favor toktik-favor/main.go
go build -o toktik-comment/bin/toktik-comment toktik-comment/main.go
:build_linux_amd64
echo 编译Linux版本64位
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o toktik-api/bin/toktik-api toktik-api/main.go
go build -o toktik-user/bin/toktik-user toktik-user/main.go
go build -o toktik-interaction/bin/toktik-interaction toktik-interaction/main.go
go build -o toktik-video/bin/toktik-video toktik-video/main.go
go build -o toktik-chat/bin/toktik-chat toktik-chat/main.go
go build -o toktik-favor/bin/toktik-favor toktik-favor/main.go
go build -o toktik-comment/bin/toktik-comment toktik-comment/main.go