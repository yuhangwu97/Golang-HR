#!/bin/bash

# Golang-HR 人力资源管理系统启动脚本

echo "🚀 启动 Golang-HR 人力资源管理系统"

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误: Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        echo "❌ 错误: 端口 $port 已被占用，请先释放该端口"
        exit 1
    fi
}

echo "🔍 检查端口占用情况..."
check_port 80
check_port 8080
check_port 3306
check_port 6379

# 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p data/mysql
mkdir -p data/redis

# 启动服务
echo "🐳 启动 Docker 容器..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "✅ 检查服务状态..."
docker-compose ps

# 显示访问地址
echo ""
echo "🎉 Golang-HR 系统启动成功！"
echo ""
echo "📱 前端地址: http://localhost"
echo "🔧 后端API: http://localhost:8080"
echo "💾 MySQL: localhost:3306"
echo "🔴 Redis: localhost:6379"
echo ""
echo "👤 默认账户:"
echo "   邮箱: admin@example.com"
echo "   密码: admin123"
echo ""
echo "🛠️  管理命令:"
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo ""