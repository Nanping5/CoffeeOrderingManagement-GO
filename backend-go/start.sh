#!/bin/bash

echo "==================================="
echo "咖啡点餐管理系统 - Go 后端启动脚本"
echo "==================================="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未检测到 Go，请先安装 Go 1.21 或更高版本"
    exit 1
fi

echo "Go 版本: $(go version)"

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "警告: 未找到 .env 文件，从 .env.example 复制..."
    cp .env.example .env
    echo "请编辑 .env 文件配置数据库密码"
    exit 1
fi

# 安装依赖
echo ""
echo "正在安装依赖..."
go mod download
go mod tidy

# 运行服务
echo ""
echo "正在启动服务..."
echo "服务地址: http://localhost:8080"
echo "按 Ctrl+C 停止服务"
echo ""

go run main.go
