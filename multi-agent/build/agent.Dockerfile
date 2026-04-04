# Agent Runtime Docker 镜像
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

# 构建 Agent Runtime
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o agent-runtime ./cmd/agent-runtime

# 运行时镜像
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN addgroup -g 1000 agent && \
    adduser -u 1000 -G agent -s /bin/sh -D agent

# 设置工作目录
WORKDIR /app

# 从 builder 复制二进制文件
COPY --from=builder /build/agent-runtime /app/

# 创建 workspace 目录
RUN mkdir -p /workspace && chown -R agent:agent /workspace

# 切换到非 root 用户
USER agent

# 设置环境变量
ENV AGENT_MODE=runtime
ENV WORKSPACE_DIR=/workspace
ENV MCP_PORT=8081

# 暴露端口
EXPOSE 8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# 启动命令
ENTRYPOINT ["/app/agent-runtime"]
