#!/bin/bash

# 创建目录函数
create_dir() {
    if [ ! -d "$1" ]; then
        mkdir -p "$1"
        if [ $? -eq 0 ]; then
            echo "目录创建成功：$1"
        else
            echo "目录创建失败：$1"
            exit 1
        fi
    else
        echo "目录已经存在：$1"
    fi
}

# 在deploy下创建目录
cd deploy

# 创建deploy/data/目录
create_dir "data"

# 创建mysql目录结构
create_dir "data/mysql/data"
create_dir "data/mysql/conf"
create_dir "data/mysql/logs"

# 创建redis目录结构
create_dir "data/redis/data"
create_dir "data/redis/conf"

# 创建redis.conf文件
if [ ! -f "data/redis/conf/redis.conf" ]; then
    touch "data/redis/conf/redis.conf"
fi

# 创建etcd目录结构
create_dir "data/etcd/data"
chmod 777 data/etcd/data
# 创建nacos目录结构
create_dir "data/nacos/data"

# 创建elasticsearch目录结构
create_dir "data/es/data"
create_dir "data/es/logs"
create_dir "data/es/plugins"

# 添加 es权限
chmod 777 data/es/data
chmod 777 data/es/logs
chmod 777 data/es/plugins

# 创建logstash目录结构
create_dir "data/logstash/log"

cd ..
