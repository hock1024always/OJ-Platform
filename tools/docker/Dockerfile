# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# 安装运行时依赖：
#   gcc/g++    - C/C++ 判题编译
#   openjdk17  - Java 判题编译和运行
#   go         - Go 判题编译
#   time       - /usr/bin/time 内存监控
#   sqlite     - 数据库
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    gcc g++ musl-dev \
    openjdk17-jdk \
    go \
    time \
    bash

COPY --from=builder /build/server .
COPY --from=builder /build/web ./web
COPY config.example.yaml ./config.yaml
COPY tools/testdata ./tools/testdata

VOLUME ["/app/data"]

EXPOSE 8080

CMD ["./server"]
