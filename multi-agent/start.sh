#!/bin/bash

# 多智能体平台启动脚本
# 使用方法: ./start.sh [orchestrator|agent|all]

set -e

# 配置
NATS_URL="nats://localhost:4222"
LLM_API_KEY="${LLM_API_KEY:-sk-7fb6784e58794327b68e4c2289d9ddf7}"
LLM_PROVIDER="${LLM_PROVIDER:-deepseek}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_deps() {
    log_info "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed"
        exit 1
    fi
    
    log_info "Go version: $(go version)"
}

# 启动 NATS (使用 Docker)
start_nats() {
    log_info "Starting NATS server..."
    
    if docker ps | grep -q multi-agent-nats; then
        log_info "NATS already running"
        return
    fi
    
    docker run -d --name multi-agent-nats \
        -p 4222:4222 \
        -p 8222:8222 \
        nats:2.10-alpine \
        -m 8222 -js
    
    sleep 2
    log_info "NATS started on port 4222"
}

# 停止 NATS
stop_nats() {
    log_info "Stopping NATS..."
    docker stop multi-agent-nats 2>/dev/null || true
    docker rm multi-agent-nats 2>/dev/null || true
}

# 启动 Orchestrator
start_orchestrator() {
    log_info "Starting Orchestrator..."
    
    cd /home/haoqian.li/compile_dockers/oj-platform/multi-agent
    
    NATS_URL=$NATS_URL \
    LLM_API_KEY=$LLM_API_KEY \
    LLM_PROVIDER=$LLM_PROVIDER \
    go run cmd/orchestrator/main.go &
    
    ORCH_PID=$!
    echo $ORCH_PID > /tmp/orchestrator.pid
    log_info "Orchestrator started (PID: $ORCH_PID)"
}

# 停止 Orchestrator
stop_orchestrator() {
    log_info "Stopping Orchestrator..."
    
    if [ -f /tmp/orchestrator.pid ]; then
        kill $(cat /tmp/orchestrator.pid) 2>/dev/null || true
        rm /tmp/orchestrator.pid
    fi
}

# 启动 Agent
start_agent() {
    local agent_id=$1
    local agent_name=$2
    local agent_type=$3
    
    log_info "Starting Agent: $agent_name ($agent_type)..."
    
    cd /home/haoqian.li/compile_dockers/oj-platform/multi-agent
    
    AGENT_ID=$agent_id \
    AGENT_NAME="$agent_name" \
    AGENT_TYPE=$agent_type \
    NATS_URL=$NATS_URL \
    LLM_API_KEY=$LLM_API_KEY \
    LLM_PROVIDER=$LLM_PROVIDER \
    go run cmd/agent-runtime/main.go &
    
    AGENT_PID=$!
    echo $AGENT_PID > /tmp/agent-$agent_id.pid
    log_info "Agent $agent_id started (PID: $AGENT_PID)"
}

# 停止所有 Agent
stop_agents() {
    log_info "Stopping all agents..."
    
    for pid_file in /tmp/agent-*.pid; do
        if [ -f "$pid_file" ]; then
            kill $(cat "$pid_file") 2>/dev/null || true
            rm "$pid_file"
        fi
    done
}

# 停止所有服务
stop_all() {
    stop_agents
    stop_orchestrator
    stop_nats
    log_info "All services stopped"
}

# 启动所有服务
start_all() {
    check_deps
    start_nats
    sleep 2
    start_orchestrator
    sleep 2
    start_agent "dev-1" "研发工程师-1" "developer"
    start_agent "dev-2" "研发工程师-2" "developer"
    start_agent "test-1" "测试工程师-1" "tester"
    start_agent "arch-1" "架构师-1" "architect"
    start_agent "ops-1" "运维工程师-1" "devops"
    
    log_info "=========================================="
    log_info "Multi-Agent Platform Started!"
    log_info "=========================================="
    log_info "Orchestrator: http://localhost:8080"
    log_info "NATS Monitor: http://localhost:8222"
    log_info ""
    log_info "Agents running:"
    log_info "  - dev-1, dev-2 (Developer)"
    log_info "  - test-1 (Tester)"
    log_info "  - arch-1 (Architect)"
    log_info "  - ops-1 (DevOps)"
    log_info ""
    log_info "Press Ctrl+C to stop all services"
    
    # 等待中断信号
    trap stop_all EXIT
    wait
}

# 测试 API
test_api() {
    log_info "Testing API..."
    
    # 健康检查
    curl -s http://localhost:8080/health && echo ""
    
    # 创建任务
    log_info "Creating a test task..."
    curl -s -X POST http://localhost:8080/api/v1/tasks \
        -H "Content-Type: application/json" \
        -d '{"title":"测试任务","description":"这是一个测试任务","created_by":"test"}' | jq .
    
    # 列出任务
    log_info "Listing tasks..."
    curl -s http://localhost:8080/api/v1/tasks | jq .
    
    # 列出 Agent
    log_info "Listing agents..."
    curl -s http://localhost:8080/api/v1/agents | jq .
}

# 主入口
case "${1:-all}" in
    start)
        start_all
        ;;
    stop)
        stop_all
        ;;
    restart)
        stop_all
        sleep 2
        start_all
        ;;
    nats)
        start_nats
        ;;
    orchestrator)
        check_deps
        start_orchestrator
        wait
        ;;
    agent)
        check_deps
        start_agent "${2:-agent-1}" "${3:-Agent-1}" "${4:-developer}"
        wait
        ;;
    test)
        test_api
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|nats|orchestrator|agent|test}"
        echo ""
        echo "Commands:"
        echo "  start       - Start all services"
        echo "  stop        - Stop all services"
        echo "  restart     - Restart all services"
        echo "  nats        - Start NATS only"
        echo "  orchestrator- Start Orchestrator only"
        echo "  agent       - Start a single agent (args: id name type)"
        echo "  test        - Test the API"
        exit 1
        ;;
esac
