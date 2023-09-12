#!/bin/bash

# 1. 创建容器卷目录
chmod 777 scripts/build-volumes.sh
./scripts/build-volumes.sh

# 2. 创建每个服务的可执行文件
chmod 777 scripts/build-bin.sh
./scripts/build-bin.sh

# 3. 构建每个服务的Docker镜像
chmod 777 scripts/build-dockerfile.sh
./scripts/build-dockerfile.sh

