#!/bin/bash
# OJ Platform 一键部署脚本
# 用法：
#   ./deploy.sh          本地部署（自动检测依赖）
#   ./deploy.sh docker   Docker 部署
#   ./deploy.sh reset    清空数据库并重新导入题库
#   ./deploy.sh stop     停止服务

set -e

# ===== 颜色 =====
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC}  $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
error()   { echo -e "${RED}[ERR]${NC}  $1"; exit 1; }

# ===== 常量 =====
PORT=8080
DB_FILE="oj_platform.db"
BIN="./bin/server"
LOG="server.log"
PID_FILE="server.pid"

# ===== 停止服务 =====
stop_service() {
    info "停止已有服务..."
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            kill "$PID"
            rm -f "$PID_FILE"
            success "服务已停止 (PID: $PID)"
        else
            rm -f "$PID_FILE"
        fi
    elif lsof -ti:"$PORT" &>/dev/null 2>&1; then
        lsof -ti:"$PORT" | xargs kill -9 2>/dev/null || true
        success "端口 $PORT 已释放"
    else
        info "无运行中的服务"
    fi
}

if [ "$1" = "stop" ]; then
    stop_service
    exit 0
fi

# ===== 检查依赖 =====
check_deps() {
    info "检查运行依赖..."

    command -v go &>/dev/null || error "未找到 Go，请安装 Go 1.21+"
    GO_VER=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' | head -1)
    success "Go $GO_VER"

    command -v gcc &>/dev/null  && success "GCC $(gcc --version | head -1 | grep -oP '\d+\.\d+\.\d+' | head -1)" \
                                || warn "GCC 未安装，C 题目将不可用"
    command -v g++ &>/dev/null  && success "G++ $(g++ --version | head -1 | grep -oP '\d+\.\d+\.\d+' | head -1)" \
                                || warn "G++ 未安装，C++ 题目将不可用"
    command -v java &>/dev/null && success "Java $(java -version 2>&1 | head -1)" \
                                || warn "Java 未安装，Java 题目将不可用"

    if ! command -v /usr/bin/time &>/dev/null; then
        warn "/usr/bin/time 未找到，内存统计可能不准确"
    else
        success "/usr/bin/time 可用"
    fi
}

# ===== 编译 =====
build() {
    info "编译项目..."
    mkdir -p bin
    CGO_ENABLED=1 go build -o "$BIN" ./cmd/server/ || error "编译失败"
    success "编译完成 → $BIN"
}

# ===== 导入题库 =====
import_data() {
    info "导入 LeetCode Hot 100 题库..."
    go run scripts/import_leetcode.go 2>&1 | tail -3
    success "题库导入完成"

    info "生成 50 组测试用例..."
    go run scripts/gen_testcases.go 2>&1 | tail -3
    success "测试用例生成完成"
}

# ===== 本地启动 =====
start_local() {
    stop_service

    info "启动服务..."
    nohup "$BIN" > "$LOG" 2>&1 &
    echo $! > "$PID_FILE"
    sleep 2

    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        success "服务已启动 (PID: $PID)"
    else
        error "服务启动失败，查看 $LOG"
    fi
}

# ===== Docker 部署 =====
deploy_docker() {
    command -v docker &>/dev/null        || error "Docker 未安装"
    command -v docker-compose &>/dev/null || command -v docker &>/dev/null || error "docker-compose 未找到"

    info "构建 Docker 镜像..."
    docker build -t oj-platform:latest .

    info "启动容器..."
    docker-compose up -d

    success "Docker 容器已启动"
}

# ===== 健康检查 =====
health_check() {
    info "健康检查..."
    for i in $(seq 1 10); do
        if curl -sf "http://localhost:$PORT/health" | grep -q "ok" 2>/dev/null; then
            success "服务健康"
            return 0
        fi
        sleep 1
    done
    warn "健康检查超时，请查看 $LOG"
}

# ===== 显示部署信息 =====
show_info() {
    echo ""
    echo "========================================"
    echo "  OJ Platform 部署成功"
    echo "========================================"
    echo ""
    echo "  访问地址  :  http://localhost:$PORT"
    echo "  健康检查  :  http://localhost:$PORT/health"
    echo "  API 文档  :  http://localhost:$PORT/api/v1/"
    echo "  日志文件  :  $LOG"
    echo ""
    echo "  注册账号（首次使用）："
    echo "  curl -X POST http://localhost:$PORT/api/v1/register \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"admin\",\"email\":\"admin@example.com\",\"password\":\"123456\"}'"
    echo ""
    echo "  停止服务  :  ./deploy.sh stop"
    echo "  重置数据  :  ./deploy.sh reset"
    echo "========================================"
}

# ===== 主流程 =====
echo ""
echo "  OJ Platform 一键部署"
echo "  ========================"
echo ""

case "$1" in
    docker)
        check_deps
        deploy_docker
        health_check
        show_info
        ;;
    reset)
        warn "即将清空数据库并重新导入题库..."
        rm -f "$DB_FILE"
        import_data
        success "数据已重置，请重启服务：./deploy.sh"
        ;;
    *)
        check_deps
        build

        # 首次运行或 reset 后导入题库
        if [ ! -f "$DB_FILE" ]; then
            import_data
        else
            info "已有数据库 ($DB_FILE)，跳过导入。如需重置请运行 ./deploy.sh reset"
        fi

        start_local
        health_check
        show_info
        ;;
esac
