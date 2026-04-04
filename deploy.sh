#!/bin/bash
# OJ Platform 一键部署脚本
#
# 用法：
#   ./deploy.sh              一键部署（编译 + 导入题库 + 启动）
#   ./deploy.sh docker       Docker 部署
#   ./deploy.sh stop         停止服务
#   ./deploy.sh restart      重启服务
#   ./deploy.sh reset        清空数据库并重新导入题库
#   ./deploy.sh status       查看服务状态

set -e

# ===== 颜色 =====
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[ OK ]${NC} $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
error()   { echo -e "${RED}[ERR]${NC}  $1"; exit 1; }

# ===== 定位项目根目录 =====
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# ===== 常量 =====
DB_FILE="oj_platform.db"
BIN="./bin/server"
LOG="server.log"
PID_FILE="server.pid"
CONFIG="config.yaml"
CONFIG_EXAMPLE="config.example.yaml"

# 从 config.yaml 读取端口（默认 8080）
get_port() {
    local p=8080
    if [ -f "$CONFIG" ]; then
        p=$(grep -E '^\s+port:' "$CONFIG" | head -1 | awk '{print $2}')
    fi
    echo "${p:-8080}"
}

# ===== 停止服务 =====
stop_service() {
    info "停止已有服务..."
    local PORT=$(get_port)

    if [ -f "$PID_FILE" ]; then
        local PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            kill "$PID"
            sleep 1
            rm -f "$PID_FILE"
            success "服务已停止 (PID: $PID)"
            return 0
        else
            rm -f "$PID_FILE"
        fi
    fi

    # fallback: 按端口查找
    if command -v ss &>/dev/null; then
        local PID=$(ss -tlnp "sport = :$PORT" 2>/dev/null | grep -oP 'pid=\K\d+' | head -1)
        if [ -n "$PID" ]; then
            kill "$PID" 2>/dev/null || true
            sleep 1
            success "端口 $PORT 已释放 (PID: $PID)"
            return 0
        fi
    elif command -v lsof &>/dev/null; then
        lsof -ti:"$PORT" 2>/dev/null | xargs kill 2>/dev/null && {
            sleep 1
            success "端口 $PORT 已释放"
            return 0
        }
    fi

    info "无运行中的服务"
}

# ===== 检查依赖 =====
check_deps() {
    info "检查运行依赖..."

    command -v go &>/dev/null || error "未找到 Go，请安装 Go 1.21+"
    GO_VER=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' | head -1)
    success "Go $GO_VER"

    command -v gcc &>/dev/null  && success "GCC $(gcc -dumpversion)" \
                                || warn "GCC 未安装，C 语言提交将不可用"
    command -v g++ &>/dev/null  && success "G++ $(g++ -dumpversion)" \
                                || warn "G++ 未安装，C++ 语言提交将不可用"
    command -v javac &>/dev/null && success "Java $(javac -version 2>&1 | awk '{print $2}')" \
                                || warn "Java 未安装，Java 语言提交将不可用"

    if [ -x /usr/bin/time ]; then
        success "/usr/bin/time 可用（内存监控）"
    else
        warn "/usr/bin/time 未找到，内存限制检测将不可用"
    fi
    echo ""
}

# ===== 初始化配置 =====
init_config() {
    if [ ! -f "$CONFIG" ]; then
        if [ -f "$CONFIG_EXAMPLE" ]; then
            cp "$CONFIG_EXAMPLE" "$CONFIG"
            success "已从 $CONFIG_EXAMPLE 创建配置文件"
        else
            error "未找到 $CONFIG 或 $CONFIG_EXAMPLE，请检查项目完整性"
        fi
    fi
}

# ===== 编译 =====
build() {
    info "编译项目..."
    mkdir -p bin
    CGO_ENABLED=1 go build -o "$BIN" ./cmd/server/ || error "编译失败"
    success "编译完成 -> $BIN"
}

# ===== 导入题库 =====
import_data() {
    # 先启动一次让 GORM 自动建表
    info "初始化数据库（自动建表）..."
    nohup "$BIN" > /dev/null 2>&1 &
    local TMP_PID=$!
    sleep 3
    kill "$TMP_PID" 2>/dev/null || true
    sleep 1
    success "数据库初始化完成"

    # 导入预置测试数据
    if [ -f "tools/testdata/test_cases.sql" ]; then
        if command -v sqlite3 &>/dev/null; then
            info "导入预置题库和测试数据..."
            sqlite3 "$DB_FILE" < tools/testdata/test_cases.sql 2>/dev/null && \
                success "题库导入完成" || warn "部分数据可能已存在，已跳过"
        else
            warn "sqlite3 未安装，跳过预置数据导入。可通过 /create-problem.html 手动出题"
        fi
    fi

    # 解压附加数据包
    if [ -f "tools/testdata/oj_test_data_v1.0.tar.gz" ]; then
        info "解压附加测试数据..."
        tar -xzf tools/testdata/oj_test_data_v1.0.tar.gz -C tools/testdata/ 2>/dev/null || true
        success "附加数据解压完成"
    fi
}

# ===== 本地启动 =====
start_local() {
    stop_service

    local PORT=$(get_port)
    info "启动服务 (端口: $PORT)..."
    nohup "$BIN" > "$LOG" 2>&1 &
    echo $! > "$PID_FILE"
    sleep 2

    local PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        success "服务已启动 (PID: $PID)"
    else
        error "服务启动失败，查看日志: tail -50 $LOG"
    fi
}

# ===== 健康检查 =====
health_check() {
    local PORT=$(get_port)
    info "健康检查..."
    for i in $(seq 1 10); do
        if curl -sf "http://localhost:$PORT/health" 2>/dev/null | grep -q "ok"; then
            success "服务健康"
            return 0
        fi
        sleep 1
    done
    warn "健康检查超时，服务可能仍在启动中，请查看 $LOG"
}

# ===== Docker 部署 =====
deploy_docker() {
    command -v docker &>/dev/null || error "Docker 未安装"

    info "构建 Docker 镜像..."
    docker build -t oj-platform:latest -f tools/docker/Dockerfile . || error "Docker 构建失败"

    if command -v docker-compose &>/dev/null; then
        docker-compose -f tools/docker/docker-compose.yml up -d
    else
        docker compose -f tools/docker/docker-compose.yml up -d
    fi

    success "Docker 容器已启动"
}

# ===== 查看状态 =====
show_status() {
    local PORT=$(get_port)
    if [ -f "$PID_FILE" ]; then
        local PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            success "服务运行中 (PID: $PID, 端口: $PORT)"
            return 0
        fi
    fi
    warn "服务未运行"
    return 1
}

# ===== 显示部署信息 =====
show_info() {
    local PORT=$(get_port)
    local LOCAL_IP=$(hostname -I 2>/dev/null | awk '{print $1}')
    LOCAL_IP=${LOCAL_IP:-"localhost"}

    echo ""
    echo "========================================"
    echo "  OJ Platform 部署成功"
    echo "========================================"
    echo ""
    echo "  访问地址  :  http://${LOCAL_IP}:${PORT}"
    echo "  出题管理  :  http://${LOCAL_IP}:${PORT}/create-problem.html"
    echo "  管理面板  :  http://${LOCAL_IP}:${PORT}/admin.html"
    echo "  健康检查  :  http://localhost:${PORT}/health"
    echo "  日志文件  :  $LOG"
    echo ""
    echo "  首次使用请注册账号："
    echo "  curl -X POST http://localhost:${PORT}/api/v1/register \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"admin\",\"email\":\"admin@oj.com\",\"password\":\"admin123\"}'"
    echo ""
    echo "  管理命令："
    echo "    ./deploy.sh stop      停止服务"
    echo "    ./deploy.sh restart   重启服务"
    echo "    ./deploy.sh status    查看状态"
    echo "    ./deploy.sh reset     重置数据库"
    echo "========================================"
}

# ===== 主流程 =====
echo ""
echo "  OJ Platform 部署工具"
echo "  ========================"
echo ""

case "${1:-deploy}" in
    deploy)
        check_deps
        init_config
        build

        if [ ! -f "$DB_FILE" ]; then
            import_data
        else
            info "已有数据库 ($DB_FILE)，跳过导入。如需重置: ./deploy.sh reset"
        fi

        start_local
        health_check
        show_info
        ;;
    docker)
        init_config
        deploy_docker
        health_check
        show_info
        ;;
    stop)
        stop_service
        ;;
    restart)
        stop_service
        sleep 1
        init_config
        if [ ! -f "$BIN" ]; then
            check_deps
            build
        fi
        start_local
        health_check
        success "服务已重启"
        ;;
    reset)
        warn "即将清空数据库并重新导入..."
        stop_service
        rm -f "$DB_FILE"
        init_config
        if [ ! -f "$BIN" ]; then
            check_deps
            build
        fi
        import_data
        start_local
        health_check
        show_info
        ;;
    status)
        show_status
        ;;
    *)
        echo "用法: $0 {deploy|docker|stop|restart|reset|status}"
        echo ""
        echo "  deploy   - 编译并部署（默认）"
        echo "  docker   - Docker 容器部署"
        echo "  stop     - 停止服务"
        echo "  restart  - 重启服务"
        echo "  reset    - 清空数据库并重新导入"
        echo "  status   - 查看服务状态"
        exit 1
        ;;
esac
