.PHONY: run build test clean docker-up docker-down

# 运行应用
run:
	go run cmd/server/main.go

# 编译
build:
	go build -o bin/server cmd/server/main.go

# 测试
test:
	go test -v ./...

# 清理
clean:
	rm -rf bin/

# Docker启动
docker-up:
	docker-compose up -d

# Docker停止
docker-down:
	docker-compose down

# 数据库迁移
migrate:
	go run cmd/server/main.go migrate

# 查看日志
logs:
	docker-compose logs -f app
