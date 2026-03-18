# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git gcc musl-dev

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译（启用CGO以支持SQLite）
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /build/main .
COPY --from=builder /build/config.yaml .
COPY --from=builder /build/web ./web

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./main"]
