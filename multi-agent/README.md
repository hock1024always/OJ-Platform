# 多智能体协作平台 (Multi-Agent Collaboration Platform)

一个基于 Kubernetes 的多智能体协作系统，模拟公司组织架构，通过总控 Orchestrator 协调多个 AI Agent 协作完成任务。

## 核心特性

- 🏢 **公司隐喻**: 研发、测试、架构、运维等角色分工协作
- 🔒 **K8s 隔离**: 每个 Agent 运行在独立的 Namespace 中，资源受限
- 🎯 **总控调度**: Orchestrator 负责任务分解、分配和协调
- 🛠️ **MCP Skills**: 模块化技能系统，支持代码生成、审查、测试等
- 💬 **自然语言**: 用户通过自然语言描述需求，系统自动分配任务
- 📊 **可视化**: Web 界面实时查看 Agent 状态和协作流程

## 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                     可视化工作区 (Web UI)                     │
├─────────────────────────────────────────────────────────────┤
│              Orchestrator 总控 (Go + K8s Client)              │
│  - 任务解析 (NLP)    - Agent 调度    - 消息总线               │
├─────────────────────────────────────────────────────────────┤
│              Kubernetes 集群                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Namespace   │  │ Namespace   │  │ Namespace   │         │
│  │ agent-dev   │  │ agent-test  │  │ agent-ops   │         │
│  │  ┌───────┐  │  │  ┌───────┐  │  │  ┌───────┐  │         │
│  │  │ Agent │  │  │  │ Agent │  │  │  │ Agent │  │         │
│  │  │ [Pod] │  │  │  │ [Pod] │  │  │  │ [Pod] │  │         │
│  │  └───────┘  │  │  └───────┘  │  │  └───────┘  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## 快速开始

### 1. 部署到 Kubernetes

```bash
# 创建 Namespace 和 RBAC
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/orchestrator.yaml

# 部署 Orchestrator
kubectl apply -f deploy/k8s/orchestrator.yaml
```

### 2. 启动 Orchestrator (本地开发)

```bash
cd multi-agent
go run cmd/orchestrator/main.go
```

### 3. 启动 Agent (本地开发)

```bash
# 终端 1 - 研发 Agent
AGENT_ID=dev-1 AGENT_NAME=研发-1 AGENT_TYPE=developer go run cmd/agent-runtime/main.go

# 终端 2 - 测试 Agent
AGENT_ID=test-1 AGENT_NAME=测试-1 AGENT_TYPE=tester go run cmd/agent-runtime/main.go
```

### 4. 启动前端

```bash
cd web
npm install
npm start
```

访问 http://localhost:3000

## API 接口

### Agent 管理

```bash
# 列出所有 Agent
GET /api/v1/agents

# 创建 Agent
POST /api/v1/agents
{
  "name": "研发-1",
  "type": "developer",
  "resources": {
    "cpu": "1",
    "memory": "2Gi"
  }
}

# 删除 Agent
DELETE /api/v1/agents/:id
```

### 任务管理

```bash
# 列出所有任务
GET /api/v1/tasks

# 创建任务
POST /api/v1/tasks
{
  "title": "开发登录系统",
  "description": "实现用户登录功能",
  "created_by": "user"
}

# NLP 任务创建
POST /api/v1/nlp/task
{
  "input": "帮我开发一个电商订单系统"
}
```

### WebSocket

```bash
ws://localhost:8080/ws
```

## MCP Skills

内置 Skills:

| Skill | 描述 | 适用 Agent |
|-------|------|-----------|
| `code_generation` | 代码生成 | developer |
| `code_review` | 代码审查 | developer, tester |
| `debug` | 调试 | developer |
| `test_generation` | 测试生成 | tester |
| `system_design` | 系统设计 | architect |
| `deploy` | 部署 | devops |
| `api_call` | API 调用 | devops |
| `file_operation` | 文件操作 | developer |

## 配置

### Agent 类型配置 `configs/agent-types.yaml`

```yaml
agentTypes:
  - name: developer
    displayName: "研发工程师"
    resources:
      requests:
        cpu: "1"
        memory: "2Gi"
    skills:
      - code_generation
      - code_review
```

## 开发计划

- [x] 基础架构设计
- [x] Orchestrator 总控实现
- [x] K8s Namespace 隔离
- [x] Web 可视化界面
- [x] MCP Skill 系统
- [ ] LLM 集成 (DeepSeek/OpenAI)
- [ ] 消息队列 (NATS/Redis)
- [ ] 持久化存储
- [ ] 监控告警

## 与 OJ 平台的关系

这个多智能体平台可以作为 OJ 平台的扩展:

1. **题目生成**: AI Agent 协作生成新的算法题目
2. **题解生成**: 自动生成多种解法和讲解
3. **测试用例生成**: 自动生成边界测试用例
4. **代码审查**: 对用户提交的代码进行智能审查

## License

MIT
