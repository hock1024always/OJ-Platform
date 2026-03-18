#!/bin/bash

# OJ Platform 快速启动脚本

echo "🚀 Starting OJ Platform..."

# 检查PostgreSQL是否运行
if ! docker ps | grep -q oj_postgres; then
    echo "📦 Starting PostgreSQL..."
    docker-compose up -d postgres
    sleep 3
fi

# 检查数据库连接
echo "🔍 Checking database connection..."
until docker exec oj_postgres pg_isready -U oj_user; do
    echo "⏳ Waiting for PostgreSQL to be ready..."
    sleep 1
done

# 运行应用
echo "🎯 Running application..."
go run cmd/server/main.go
