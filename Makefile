.PHONY: run build test clean deploy stop restart reset status docker-up docker-down

# 开发运行
run:
	go run cmd/server/main.go

# 编译
build:
	@mkdir -p bin
	CGO_ENABLED=1 go build -o bin/server ./cmd/server/

# 测试
test:
	go test -v ./...

# 清理编译产物
clean:
	rm -rf bin/

# 一键部署
deploy:
	./deploy.sh

# 停止服务
stop:
	./deploy.sh stop

# 重启服务
restart:
	./deploy.sh restart

# 重置数据库
reset:
	./deploy.sh reset

# 查看状态
status:
	./deploy.sh status

# Docker 构建和启动
docker-up:
	docker build -t oj-platform:latest -f tools/docker/Dockerfile .
	docker-compose -f tools/docker/docker-compose.yml up -d

# Docker 停止
docker-down:
	docker-compose -f tools/docker/docker-compose.yml down
