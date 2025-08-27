#!/bin/bash

# Claude Code Companion Docker 快速启动脚本

set -e

echo "=== Claude Code Companion Docker 快速启动 ==="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "错误: Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 创建数据目录
echo "创建数据目录..."
mkdir -p ./data/config ./data/logs

# 构建镜像
echo "构建 Docker 镜像..."
docker build -f docker/Dockerfile -t claude-code-companion:latest .

# 启动服务
echo "启动服务..."
cd docker
if command -v docker-compose &> /dev/null; then
    docker-compose up -d
else
    docker compose up -d
fi
cd ..

# 等待服务启动
echo "等待服务启动..."
sleep 5

# 检查服务状态
cd docker
if command -v docker-compose &> /dev/null; then
    docker-compose ps
else
    docker compose ps
fi
cd ..

echo ""
echo "=== 启动完成 ==="
echo "代理服务: http://localhost:8080"
echo "管理界面: http://localhost:8080/admin/"
echo ""
echo "查看日志: cd docker && docker-compose logs -f"
echo "停止服务: cd docker && docker-compose down"
echo ""
echo "配置文件位置: ./data/config/config.yaml"
echo "日志文件位置: ./data/logs/"
