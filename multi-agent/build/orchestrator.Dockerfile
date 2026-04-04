# Orchestrator Docker 镜像
FROM golang:1.21-alpine AS builder

# 安装依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /build

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建 Orchestrator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orchestrator ./cmd/orchestrator

# 运行时镜像
FROM alpine:latest

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN addgroup -g 1000 orchestrator && \
    adduser -u 1000 -G orchestrator -s /bin/sh -D orchestrator

# 设置工作目录
WORKDIR /app

# 从 builder 复制二进制文件
COPY --from=builder /build/orchestrator /app/

# 复制前端构建文件
COPY --from=builder /build/web/build /app/web/build

# 切换到非 root 用户
USER orchestrator

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
ENTRYPOINT ["/app/orchestrator"]
