# 多智能体平台 - 优化规划

本文档记录后续优化方向和详细设计方案。

---

## 1. LLVM 合作者接入流程

### 1.1 接入流程

```
┌─────────────────────────────────────────────────────────────┐
│                    合作者接入流程                            │
├─────────────────────────────────────────────────────────────┤
│  1. Fork 项目 → 2. 实现接口 → 3. 提交 PR → 4. 测试验证      │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 接口定义

文件位置: `pkg/compiler/interface.go`

```go
type CompilerPlugin interface {
    // 编译器名称 (如: llvm-clang, gcc, go-compiler)
    Name() string

    // 支持的语言列表
    SupportedLanguages() []string

    // 编译入口 (源码 → 可执行文件/字节码)
    Compile(ctx context.Context, req *CompileRequest) (*CompileResult, error)

    // 可选：优化级别
    OptimizationLevels() []string // -O0, -O1, -O2, -O3
}

type CompileRequest struct {
    Language    string            // c, cpp, rust, etc.
    SourceCode  string            // 源代码
    Options     map[string]string // 编译选项
    Timeout     time.Duration     // 超时
}

type CompileResult struct {
    Success     bool
    Binary      []byte            // 编译产物
    Errors      []CompileError    // 编译错误
    Warnings    []string
    Stats       CompileStats      // 编译耗时、内存等
}
```

### 1.3 目录结构

```
pkg/compilers/
├── interface.go          # 接口定义
├── base.go               # 基础实现（沙箱调用）
├── factory.go            # 编译器工厂
├── llvm/                 # LLVM 合作者实现
│   ├── compiler.go
│   └── README.md         # 接入说明
├── gcc/
│   └── compiler.go
└── go/
    └── compiler.go
```

### 1.4 合作者接入文档模板

文件位置: `pkg/compilers/llvm/README.md`

```markdown
# LLVM 编译器插件接入指南

## 快速开始

1. 实现 `CompilerPlugin` 接口
2. 注册到编译器工厂: `compiler.Register(&LLVMCompiler{})`
3. 在 `configs/compilers.yaml` 配置

## 接口实现示例

```go
package llvm

type LLVMCompiler struct{}

func (c *LLVMCompiler) Name() string {
    return "llvm-clang"
}

func (c *LLVMCompiler) SupportedLanguages() []string {
    return []string{"c", "cpp", "rust"}
}

func (c *LLVMCompiler) Compile(ctx context.Context, req *CompileRequest) (*CompileResult, error) {
    // 1. 写入源码到沙箱
    // 2. 调用 clang 编译
    // 3. 返回编译结果
}
```

## 测试要求

- 单元测试覆盖率 > 80%
- 集成测试通过沙箱安全验证
- 性能基准测试

## 配置示例

```yaml
# configs/compilers.yaml
compilers:
  - name: llvm-clang
    enabled: true
    languages: [c, cpp]
    optimization_levels: [-O0, -O1, -O2, -O3]
    timeout: 30s
```
```

### 1.5 实现清单

- [ ] 定义 `CompilerPlugin` 接口
- [ ] 实现编译器工厂模式
- [ ] 创建基础编译器实现（调用 go-judge 沙箱）
- [ ] 编写合作者接入文档
- [ ] 添加编译器配置文件支持

---

## 2. RAG 向量库 + 题库自动化

### 2.1 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    RAG 解题助手架构                          │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐    ┌──────────┐    ┌──────────────────────┐   │
│  │ 题库 YAML │ →  │ 向量化   │ →  │ Vector Store         │   │
│  │ (标准格式)│    │ Embedding│    │ (Milvus/Chroma)      │   │
│  └──────────┘    └──────────┘    └──────────────────────┘   │
│                                          ↑                   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ MCP Tool: search_similar_problems, get_hint, explain │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 题库标准格式

文件位置: `problems/001-two-sum.yaml`

```yaml
id: "001"
title: "两数之和"
difficulty: easy
tags: [array, hash-table, two-pointers]

description: |
  给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数，并返回它们的数组下标。

examples:
  - input: "nums = [2,7,11,15], target = 9"
    output: "[0,1]"
    explanation: "因为 nums[0] + nums[1] == 9，返回 [0, 1]"

constraints:
  - "2 <= nums.length <= 10^4"
  - "-10^9 <= nums[i] <= 10^9"
  - "只会存在一个有效答案"

# RAG 向量化字段
keywords: [数组, 哈希表, 双指针, 查找, 两数之和]
related_problems: [015, 167, 170]  # 相似题目 ID
solution_patterns:
  - pattern: "哈希表一次遍历"
    hint: "使用哈希表存储已遍历元素及其索引，查找 target - nums[i]"
    complexity: "O(n) 时间, O(n) 空间"
  - pattern: "暴力枚举"
    hint: "两层循环枚举所有可能的数对"
    complexity: "O(n²) 时间, O(1) 空间"

# MCP 工具元数据
mcp_tools:
  - name: "get_hint"
    args: ["level"]  # hint 级别 1-3
  - name: "check_solution"
    args: ["code", "language"]
  - name: "explain_pattern"
    args: ["pattern_name"]
```

### 2.3 自动入库流程

文件位置: `pkg/problem/ingest.go`

```go
// ProblemIngestor 题目入库器
type ProblemIngestor struct {
    vectorStore  VectorStore
    db           *sql.DB
    embedClient  EmbeddingClient
}

// IngestProblem 将题目入库（YAML → 向量库 + 关系库）
func (p *ProblemIngestor) InestProblem(ctx context.Context, yamlPath string) error {
    // 1. 解析 YAML
    problem, err := ParseProblemYAML(yamlPath)
    if err != nil {
        return err
    }

    // 2. 验证题目格式
    if err := problem.Validate(); err != nil {
        return err
    }

    // 3. 向量化 (调用 Embedding API)
    searchableText := problem.ToSearchableText()
    embedding, err := p.embedClient.Embed(ctx, searchableText)
    if err != nil {
        return err
    }

    // 4. 存入向量库
    if err := p.vectorStore.Insert(ctx, &VectorEntry{
        ID:        problem.ID,
        Vector:    embedding,
        Metadata:  problem.ToMetadata(),
    }); err != nil {
        return err
    }

    // 5. 存入关系数据库
    if err := p.db.InsertProblem(ctx, problem); err != nil {
        return err
    }

    return nil
}

// ToSearchableText 生成用于向量检索的文本
func (p *Problem) ToSearchableText() string {
    return fmt.Sprintf("%s %s %s %s",
        p.Title,
        p.Description,
        strings.Join(p.Keywords, " "),
        strings.Join(p.Tags, " "),
    )
}
```

### 2.4 MCP 工具定义

文件位置: `pkg/mcp/rag_tools.go`

```go
// RAGTools RAG 相关 MCP 工具
var RAGTools = []mcp.Tool{
    {
        Name:        "search_similar_problems",
        Description: "搜索相似题目，用于找到与当前问题类似的已解决问题",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "query": map[string]any{
                    "type":        "string",
                    "description": "题目描述或关键词",
                },
                "top_k": map[string]any{
                    "type":        "integer",
                    "description": "返回数量，默认5",
                    "default":     5,
                },
            },
            "required": []string{"query"},
        },
        Handler: handleSearchSimilarProblems,
    },
    {
        Name:        "get_hint",
        Description: "获取解题提示，根据难度级别返回不同程度的提示",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "problem_id": map[string]any{
                    "type":        "string",
                    "description": "题目ID",
                },
                "level": map[string]any{
                    "type":        "integer",
                    "description": "提示级别 1-3，1最简单，3最详细",
                    "minimum":     1,
                    "maximum":     3,
                },
            },
            "required": []string{"problem_id", "level"},
        },
        Handler: handleGetHint,
    },
    {
        Name:        "explain_pattern",
        Description: "解释解题模式的原理和应用",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "pattern": map[string]any{
                    "type":        "string",
                    "description": "解题模式名称，如：哈希表一次遍历、双指针、滑动窗口等",
                },
            },
            "required": []string{"pattern"},
        },
        Handler: handleExplainPattern,
    },
    {
        Name:        "check_solution",
        Description: "检查代码解决方案的正确性",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "problem_id": map[string]any{
                    "type":        "string",
                    "description": "题目ID",
                },
                "code": map[string]any{
                    "type":        "string",
                    "description": "源代码",
                },
                "language": map[string]any{
                    "type":        "string",
                    "description": "编程语言",
                    "enum":        []string{"python", "go", "java", "cpp", "c"},
                },
            },
            "required": []string{"problem_id", "code", "language"},
        },
        Handler: handleCheckSolution,
    },
}
```

### 2.5 向量库选型

| 方案 | 优点 | 缺点 | 推荐场景 |
|------|------|------|---------|
| **Chroma** | 轻量、易部署、Python原生 | 生产级性能一般 | 开发测试、小规模 |
| **Milvus** | 高性能、分布式、云原生 | 部署复杂 | 生产环境、大规模 |
| **Qdrant** | Rust实现、性能好、API友好 | 社区较小 | 中等规模 |
| **pgvector** | PostgreSQL扩展、运维简单 | 性能一般 | 已有PG基础设施 |

**推荐方案**: Chroma (开发) + Milvus (生产)

### 2.6 实现清单

- [ ] 定义题目 YAML 标准格式
- [ ] 实现 YAML 解析器
- [ ] 集成 Embedding 服务 (DeepSeek/OpenAI)
- [ ] 集成向量库 (Chroma/Milvus)
- [ ] 实现 MCP RAG 工具
- [ ] 创建题目入库 CLI 工具
- [ ] 批量导入 LeetCode 题库

---

## 3. 可视化监控体系

### 3.1 监控指标体系

```
┌─────────────────────────────────────────────────────────────┐
│                      监控指标体系                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ 网络流量    │  │ Token 流速  │  │ 容器资源            │  │
│  │ - 入站/出站 │  │ - 输入/输出 │  │ - CPU/Memory        │  │
│  │ - 连接数    │  │ - 模型分布  │  │ - 磁盘 I/O          │  │
│  │ - 请求延迟  │  │ - 费用统计  │  │ - 网络 I/O          │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ 算法题执行监控                                       │    │
│  │ - 内存峰值 / 时间消耗 / 测试用例通过率               │    │
│  │ - 沙箱资源隔离状态                                   │    │
│  │ - 编译错误统计                                       │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Agent 状态监控                                       │    │
│  │ - 在线/离线状态                                      │    │
│  │ - 当前任务进度                                       │    │
│  │ - 历史执行统计                                       │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 技术栈选型

| 监控类型 | 数据源 | 存储 | 可视化 |
|---------|--------|------|--------|
| 网络流量 | cAdvisor / Node Exporter | Prometheus | Grafana |
| Token 流速 | LLM API 响应 | Prometheus | 自定义 Dashboard |
| 容器资源 | Docker Stats / cgroups | Prometheus | Grafana |
| 算法内存 | go-judge 沙箱 | ClickHouse | 自定义 |
| Agent 状态 | WebSocket 心跳 | Redis | 实时面板 |

### 3.3 指标采集实现

文件位置: `pkg/metrics/collector.go`

```go
package metrics

import (
    "context"
    "time"

    "github.com/docker/docker/client"
    "github.com/prometheus/client_golang/prometheus"
)

// MetricsCollector 指标采集器
type MetricsCollector struct {
    dockerCli  *client.Client
    registry   *prometheus.Registry

    // Prometheus 指标
    tokenCounter    *prometheus.CounterVec
    executionLatency prometheus.Histogram
    memoryPeak      prometheus.Histogram
    networkBytes    *prometheus.CounterVec
    containerCPU    *prometheus.GaugeVec
    containerMemory *prometheus.GaugeVec
}

// NewMetricsCollector 创建采集器
func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        tokenCounter: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "llm_tokens_total",
                Help: "Total number of LLM tokens processed",
            },
            []string{"model", "type"}, // type: input/output
        ),
        executionLatency: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name:    "execution_latency_seconds",
                Help:    "Code execution latency in seconds",
                Buckets: []float64{.01, .05, .1, .5, 1, 5, 10},
            },
        ),
        memoryPeak: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name:    "execution_memory_mb",
                Help:    "Peak memory usage in MB",
                Buckets: []float64{1, 5, 10, 50, 100, 500, 1000},
            },
        ),
        networkBytes: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "network_bytes_total",
                Help: "Total network bytes transferred",
            },
            []string{"direction"}, // direction: in/out
        ),
        containerCPU: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "container_cpu_percent",
                Help: "Container CPU usage percentage",
            },
            []string{"container_name"},
        ),
        containerMemory: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "container_memory_mb",
                Help: "Container memory usage in MB",
            },
            []string{"container_name"},
        ),
    }
}

// Collect 采集所有指标
func (c *MetricsCollector) Collect(ctx context.Context) {
    // 1. 容器资源指标
    c.collectContainerMetrics(ctx)

    // 2. 网络流量指标
    c.collectNetworkMetrics(ctx)
}

// collectContainerMetrics 采集容器指标
func (c *MetricsCollector) collectContainerMetrics(ctx context.Context) {
    containers, err := c.dockerCli.ContainerList(ctx, types.ContainerListOptions{})
    if err != nil {
        return
    }

    for _, ctr := range containers {
        stats, err := c.dockerCli.ContainerStats(ctx, ctr.ID, false)
        if err != nil {
            continue
        }

        // 解析 stats 并更新 Prometheus 指标
        name := ctr.Names[0]
        c.containerCPU.WithLabelValues(name).Set(stats.CPUStats.CPUUsage.Percent)
        c.containerMemory.WithLabelValues(name).Set(float64(stats.MemoryStats.Usage) / 1024 / 1024)
    }
}

// RecordTokenUsage 记录 Token 使用
func (c *MetricsCollector) RecordTokenUsage(model string, inputTokens, outputTokens int) {
    c.tokenCounter.WithLabelValues(model, "input").Add(float64(inputTokens))
    c.tokenCounter.WithLabelValues(model, "output").Add(float64(outputTokens))
}

// RecordExecution 记录执行指标
func (c *MetricsCollector) RecordExecution(duration time.Duration, peakMemoryMB float64) {
    c.executionLatency.Observe(duration.Seconds())
    c.memoryPeak.Observe(peakMemoryMB)
}
```

### 3.4 监控 API 端点

```go
// 在 Orchestrator 中添加监控端点
func (o *Orchestrator) setupMetricsRoutes() {
    // Prometheus 指标端点
    http.Handle("/metrics", promhttp.Handler())

    // 自定义监控 API
    http.HandleFunc("/api/v1/metrics/summary", o.handleMetricsSummary)
    http.HandleFunc("/api/v1/metrics/tokens", o.handleTokenMetrics)
    http.HandleFunc("/api/v1/metrics/containers", o.handleContainerMetrics)
}
```

### 3.5 前端监控面板

文件位置: `web/src/components/dashboard/MonitorPanel.tsx`

```typescript
interface MonitorData {
  network: {
    bytesIn: number;
    bytesOut: number;
    connections: number;
  };
  tokens: {
    input: number;
    output: number;
    byModel: Record<string, { input: number; output: number }>;
  };
  containers: {
    name: string;
    cpu: number;
    memory: number;
    status: string;
  }[];
  executions: {
    latency: number[];
    memoryPeak: number[];
    successRate: number;
  };
}

// WebSocket 实时推送
const useMonitorData = () => {
  const [data, setData] = useState<MonitorData | null>(null);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === 'metrics_update') {
        setData(msg.data);
      }
    };
    return () => ws.close();
  }, []);

  return data;
};
```

### 3.6 实现清单

- [ ] 集成 Prometheus 指标采集
- [ ] 实现 Docker 容器监控
- [ ] 实现 Token 使用统计
- [ ] 实现网络流量监控
- [ ] 集成 go-judge 执行指标
- [ ] 创建监控 API 端点
- [ ] 开发前端监控面板
- [ ] 配置 Grafana Dashboard (可选)

---

## 4. 像素画风 AI 公司可视化

### 4.1 设计理念

创建一个像素风格的 AI 公司场景，用户可以：
- 手动创建工位（Agent 工作位置）
- 选择模型和工具分配给 Agent
- 通过可视化界面与 Agent 交互
- 实时查看 Agent 工作状态

### 4.2 界面布局

```
┌─────────────────────────────────────────────────────────────┐
│                  🏢 像素风 AI 公司                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                  │
│   │ 👨‍💻 │ │ 👩‍💻 │ │ 🤖  │ │ 📊  │ │ 🚀  │   ← 工位格子      │
│   │dev-1│ │dev-2│ │arch │ │test │ │ops  │                  │
│   └─────┘ └─────┘ └─────┘ └─────┘ └─────┘                  │
│     🟢      🟡       🔴      🟢      🟢    ← 状态指示       │
│                                                             │
│   ┌─────────────────────────────────────────────────────┐  │
│   │                    任务看板                          │  │
│   │  [待分配] [进行中] [审核中] [已完成]                  │  │
│   └─────────────────────────────────────────────────────┘  │
│                                                             │
│   ┌─────────────────────────────────────────────────────┐  │
│   │  💬 工作区对话                                        │  │
│   │  用户: 帮我实现一个登录功能                           │  │
│   │  dev-1: 好的，我来编写代码...                         │  │
│   └─────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 前端技术栈

```yaml
框架: React 18 + TypeScript
样式:
  - 自定义 CSS 像素边框
  - 8-bit 字体: Press Start 2P
  - 动画: CSS Keyframes
状态管理: Zustand
实时通信: WebSocket
图标: 像素风 Sprite 图集
```

### 4.4 核心组件设计

#### 4.4.1 工位组件

文件位置: `web/src/components/pixel/Workstation.tsx`

```typescript
interface Workstation {
  id: string;
  agent?: Agent;
  position: { x: number; y: number };
  status: 'empty' | 'idle' | 'working' | 'error';
}

interface Agent {
  id: string;
  name: string;
  type: 'developer' | 'tester' | 'architect' | 'devops';
  model: 'deepseek' | 'gpt-4' | 'claude-3';
  tools: MCPTool[];
  currentTask?: Task;
  avatar: string; // 像素头像 URL
}

const Workstation: React.FC<{ workstation: Workstation }> = ({ workstation }) => {
  const statusColor = {
    empty: 'gray',
    idle: 'green',
    working: 'yellow',
    error: 'red',
  };

  return (
    <div className="workstation pixel-border">
      <div className="avatar">
        {workstation.agent ? (
          <PixelAvatar agent={workstation.agent} />
        ) : (
          <EmptySlot onClick={() => openCreateAgentModal()} />
        )}
      </div>
      <div className="status-indicator" style={{ color: statusColor[workstation.status] }}>
        ●
      </div>
      <div className="agent-name">
        {workstation.agent?.name || '空工位'}
      </div>
      {workstation.agent?.currentTask && (
        <div className="task-progress">
          <ProgressBar progress={workstation.agent.currentTask.progress} />
        </div>
      )}
    </div>
  );
};
```

#### 4.4.2 任务看板

文件位置: `web/src/components/pixel/TaskBoard.tsx`

```typescript
interface Task {
  id: string;
  title: string;
  status: 'pending' | 'in_progress' | 'reviewing' | 'completed';
  assignee?: string;
  priority: 'low' | 'medium' | 'high';
  createdAt: Date;
}

const TaskBoard: React.FC = () => {
  const { tasks } = useTaskStore();

  const columns = [
    { id: 'pending', title: '待分配', icon: '📋' },
    { id: 'in_progress', title: '进行中', icon: '🔄' },
    { id: 'reviewing', title: '审核中', icon: '🔍' },
    { id: 'completed', title: '已完成', icon: '✅' },
  ];

  return (
    <div className="task-board pixel-border">
      {columns.map(col => (
        <div key={col.id} className="column">
          <div className="column-header">
            {col.icon} {col.title}
          </div>
          <div className="task-list">
            {tasks.filter(t => t.status === col.id).map(task => (
              <TaskCard key={task.id} task={task} />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};
```

#### 4.4.3 工作区对话

文件位置: `web/src/components/pixel/ChatPanel.tsx`

```typescript
const ChatPanel: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');

  const sendMessage = async () => {
    // 1. 添加用户消息
    setMessages(prev => [...prev, { role: 'user', content: input }]);

    // 2. 发送到后端，由 Orchestrator 分配给合适的 Agent
    const response = await fetch('/api/v1/nlp/parse', {
      method: 'POST',
      body: JSON.stringify({ query: input }),
    });

    // 3. WebSocket 会推送 Agent 的响应
    setInput('');
  };

  return (
    <div className="chat-panel pixel-border">
      <div className="messages">
        {messages.map((msg, i) => (
          <ChatBubble key={i} message={msg} />
        ))}
      </div>
      <div className="input-area">
        <input
          value={input}
          onChange={e => setInput(e.target.value)}
          className="pixel-input"
          placeholder="输入任务或问题..."
        />
        <button onClick={sendMessage} className="pixel-btn">
          发送
        </button>
      </div>
    </div>
  );
};
```

#### 4.4.4 创建 Agent 弹窗

```typescript
const CreateAgentModal: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  const [config, setConfig] = useState({
    name: '',
    type: 'developer',
    model: 'deepseek',
    tools: [] as string[],
  });

  const agentTypes = [
    { id: 'developer', label: '研发工程师', icon: '👨‍💻' },
    { id: 'tester', label: '测试工程师', icon: '📊' },
    { id: 'architect', label: '架构师', icon: '🤖' },
    { id: 'devops', label: '运维工程师', icon: '🚀' },
  ];

  const models = [
    { id: 'deepseek', label: 'DeepSeek', icon: '🧠' },
    { id: 'gpt-4', label: 'GPT-4', icon: '🤖' },
    { id: 'claude-3', label: 'Claude 3', icon: '🎭' },
  ];

  return (
    <div className="modal-overlay">
      <div className="modal pixel-border">
        <h2>创建新 Agent</h2>

        <div className="form-group">
          <label>名称</label>
          <input
            value={config.name}
            onChange={e => setConfig({ ...config, name: e.target.value })}
            className="pixel-input"
          />
        </div>

        <div className="form-group">
          <label>类型</label>
          <div className="option-grid">
            {agentTypes.map(type => (
              <div
                key={type.id}
                className={`option-card ${config.type === type.id ? 'selected' : ''}`}
                onClick={() => setConfig({ ...config, type: type.id })}
              >
                <span className="icon">{type.icon}</span>
                <span className="label">{type.label}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="form-group">
          <label>模型</label>
          <div className="option-grid">
            {models.map(model => (
              <div
                key={model.id}
                className={`option-card ${config.model === model.id ? 'selected' : ''}`}
                onClick={() => setConfig({ ...config, model: model.id })}
              >
                <span className="icon">{model.icon}</span>
                <span className="label">{model.label}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="form-group">
          <label>MCP 工具</label>
          <ToolSelector
            selected={config.tools}
            onChange={tools => setConfig({ ...config, tools })}
          />
        </div>

        <div className="actions">
          <button onClick={onClose} className="pixel-btn secondary">取消</button>
          <button onClick={() => createAgent(config)} className="pixel-btn">创建</button>
        </div>
      </div>
    </div>
  );
};
```

### 4.5 像素风样式

文件位置: `web/src/styles/pixel.css`

```css
/* 像素边框效果 */
.pixel-border {
  border: 4px solid #2d2d2d;
  box-shadow:
    inset -4px -4px 0px 0px #1a1a1a,
    inset 4px 4px 0px 0px #4a4a4a,
    8px 8px 0px 0px rgba(0,0,0,0.2);
  background: #f0f0f0;
}

/* 像素按钮 */
.pixel-btn {
  font-family: 'Press Start 2P', cursive;
  font-size: 12px;
  padding: 12px 24px;
  border: none;
  background: #4CAF50;
  color: white;
  cursor: pointer;
  box-shadow:
    inset -4px -4px 0px 0px #388E3C,
    inset 4px 4px 0px 0px #81C784;
  transition: transform 0.1s;
}

.pixel-btn:hover {
  transform: scale(1.05);
}

.pixel-btn:active {
  transform: scale(0.95);
  box-shadow:
    inset 4px 4px 0px 0px #388E3C,
    inset -4px -4px 0px 0px #81C784;
}

/* 像素输入框 */
.pixel-input {
  font-family: 'Press Start 2P', cursive;
  font-size: 10px;
  padding: 12px;
  border: 4px solid #2d2d2d;
  background: #ffffff;
  box-shadow: inset 4px 4px 0px 0px #e0e0e0;
}

/* 工位格子 */
.workstation {
  width: 120px;
  height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 8px;
  transition: transform 0.2s;
}

.workstation:hover {
  transform: translateY(-4px);
}

/* 状态指示灯动画 */
.status-indicator {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Agent 工作动画 */
.working-animation {
  animation: typing 0.5s infinite;
}

@keyframes typing {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}
```

### 4.6 实时状态同步

```typescript
// WebSocket 消息类型
type WSMessage =
  | { type: 'agent_status'; agentId: string; status: AgentStatus }
  | { type: 'task_update'; task: Task }
  | { type: 'chat_message'; message: Message }
  | { type: 'metrics_update'; data: MonitorData };

// 前端状态管理
const useOfficeStore = create<OfficeState>((set) => ({
  workstations: [],
  tasks: [],
  messages: [],

  updateAgentStatus: (agentId, status) =>
    set(state => ({
      workstations: state.workstations.map(w =>
        w.agent?.id === agentId ? { ...w, status } : w
      ),
    })),

  addMessage: (message) =>
    set(state => ({
      messages: [...state.messages, message],
    })),
}));
```

### 4.7 实现清单

- [ ] 创建 React 项目结构
- [ ] 实现像素风 CSS 样式库
- [ ] 开发 Workstation 组件
- [ ] 开发 TaskBoard 组件
- [ ] 开发 ChatPanel 组件
- [ ] 开发 CreateAgentModal 组件
- [ ] 实现 WebSocket 实时同步
- [ ] 添加 Agent 动画效果
- [ ] 集成后端 API

---

## 5. 优先级与里程碑

### 5.1 优先级排序

| 优先级 | 任务 | 理由 | 预计工作量 |
|-------|------|------|-----------|
| **P0** | 2. RAG 向量库 | 解题助手核心能力，立即可用 | 中 |
| **P1** | 1. LLVM 接口 | 为合作者提供规范，扩展性 | 小 |
| **P1** | 4. 像素风可视化 | 用户体验，产品差异化 | 大 |
| **P2** | 3. 监控体系 | 生产环境必需，可渐进实现 | 中 |

### 5.2 里程碑规划

```
Phase 1 (当前)
├── RAG 向量库集成
├── 题库标准格式定义
└── MCP RAG 工具

Phase 2
├── 编译器插件接口
├── 合作者接入文档
└── 像素风 UI 框架

Phase 3
├── 完整像素风界面
├── Agent 创建/管理
└── 实时交互功能

Phase 4
├── 监控体系完善
├── Grafana 集成
└── 生产环境部署
```

---

## 6. 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) - 系统架构设计
- [ROADMAP.md](./ROADMAP.md) - 功能开发路线图
- [README.md](./README.md) - 项目说明
