#!/bin/bash

# 定义一个函数用于构建 Docker 镜像
build_docker_image() {
  local service_name="$1"
  local tag="$2"

  echo "构建 Docker 镜像：$service_name:$tag"
  cd "$service_name"
  docker build -f Dockerfile -t "$service_name:$tag" .
  # shellcheck disable=SC2103
  cd ..
}

# 构建各个服务的 Docker 镜像
build_docker_image "toktik-api" "0.1"
build_docker_image "toktik-user" "0.1"
build_docker_image "toktik-interaction" "0.1"
build_docker_image "toktik-video" "0.1"
build_docker_image "toktik-chat" "0.1"
build_docker_image "toktik-favor" "0.1"
build_docker_image "toktik-comment" "0.1"

# 完成提示
echo "所有 Docker 镜像构建完成。"
