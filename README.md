# OJ Platform

> 一个开源的在线判题平台，面向技术招聘场景，支持力扣风格的算法题解答，内置多语言评测引擎。

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux-lightgrey?style=flat&logo=linux)](https://linux.org)

---

## 概述

OJ Platform 是一个轻量级、可自托管的在线判题系统，专为技术招聘设计。候选人可在浏览器中直接解题，采用力扣风格的交互方式——只需编写解题函数，无需处理 I/O 模板代码。

**核心特性：**

- **力扣风格编辑器**：用户只需编写解题函数，驱动代码由平台自动拼接
- **多语言支持**：Go / C / C++ / Java，前端一键切换语言和语法高亮
- **自定义测试控制台**：分栏 Input/Output 面板，提交前可反复调试
- **性能可视化**：执行时间和内存消耗以彩色进度条实时展示
- **排行榜系统**：全局排名（按通过题数）+ 单题排名（按最快用时）
- **管理后台**：查看所有提交记录，支持筛选、分页、代码审查
- **JSON 导入题目**：管理员可通过 JSON 一键导入自定义题目及最多 100 组测试用例
- **完整 Hot 100 题库**：内置 100 道 LeetCode Hot 100，每题 50 组测试数据（见 [题单](PROBLEM_LIST.md)）
- **请求限流**：令牌桶限流中间件，全局 100 req/s，提交接口 5 req/s（per IP）
- **沙盒隔离**：`ulimit` 资源限制 + 进程组隔离，防止代码炸弹、fork 炸弹、磁盘滥写
- 并发提交处理（20 Worker 池，支持 10-100 并发用户）
- JWT 认证机制
- 后端：Go (Gin + GORM + SQLite)
- 前端：原生 JS + CodeMirror（零框架依赖）

---

## 界面预览

> 打开题目 → 选择语言 → 编写函数 → 运行测试 / 提交 → 即时反馈

```
┌─────────────────────────────────────────────────┐
│  1. 两数之和  [Easy]                              │
│  时间限制: 5000ms   内存限制: 256MB                │
├─────────────────────────────────────────────────┤
│  语言: [Go ▼] [C] [C++] [Java]                  │
│  func twoSum(nums []int, target int) []int {     │
│      // 在此编写你的代码                           │
│  }                                               │
├─────────────────────────────────────────────────┤
│  [输入] [输出]                                    │
│  2 7 11 15                                       │
│  9                                               │
├─────────────────────────────────────────────────┤
│  [运行测试]  [提交代码]                            │
└─────────────────────────────────────────────────┘
```

---

## 快速开始

### 环境要求

- Go 1.21+
- GCC / G++（C/C++ 支持）
- JDK 8+（Java 支持）
- Git

### 一键部署

```bash
git clone https://github.com/your-org/oj-platform.git
cd oj-platform

# 一键完成：环境检测 → 编译 → 导入题库 → 生成测试数据 → 启动
./deploy.sh
```

**其他命令：**

```bash
./deploy.sh docker   # Docker 容器部署
./deploy.sh reset    # 清空数据库并重新导入题库
./deploy.sh stop     # 停止服务
```

浏览器打开 `http://localhost:8080`。

首次运行请注册账号：

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","email":"admin@example.com","password":"123456"}'
```

### 手动安装

```bash
# 编译
go build -o bin/server ./cmd/server/

# 导入题库
go run scripts/import_leetcode.go

# 生成测试数据
go run scripts/gen_testcases.go

# 启动
./bin/server
```

---

## 支持语言

| 语言 | 编译器 | 编译选项 |
|------|--------|---------|
| Go | go 1.21+ | `go build` |
| C | gcc | `-O2 -lm` |
| C++ | g++ | `-O2 -std=c++17 -lm` |
| Java | javac + java | `-Xmx256m` |

编译器路径可在 `config.yaml` 中配置。

---

## 题库

内置完整 LeetCode Hot 100，共 100 道题目，覆盖 14 个专题。

完整题单见 **[PROBLEM_LIST.md](PROBLEM_LIST.md)**。

每题配备 **50 组测试用例**（含公开示例 + 自动生成的边界用例）。

---

## 请求限流

基于**令牌桶算法**，无第三方依赖，纯 Go 实现：

| 限流层 | 策略 | 适用接口 |
|--------|------|---------|
| 全局限流 | 100 req/s per IP，峰值 200 | 所有 API |
| 提交限流 | 5 req/s per IP，峰值 10 | `POST /submit`、`POST /test` |

触发限流时返回 `HTTP 429 Too Many Requests`。

实现位置：`internal/middleware/ratelimit.go`

---

## 沙盒隔离

每次代码执行均通过 `ulimit` + 进程组隔离形成轻量级沙盒：

| 限制项 | 限制值 | 说明 |
|--------|--------|------|
| 虚拟内存 (`-v`) | `memory_limit × 1024 KB` | 与配置的内存限制联动 |
| 最大文件写入 (`-f`) | 16384 blocks（~8 MB） | 防止磁盘炸弹 |
| 最大子进程数 (`-u`) | 64 | 防止 fork 炸弹 |
| 最大文件描述符 (`-n`) | 32（Java 为 64） | 限制网络连接 |
| 进程组隔离 (`Setpgid`) | 是 | 超时时一次性 kill 整棵进程树 |
| 超时强制终止 | `time_limit + 500ms` | Timer 到期后发送 SIGKILL |

实现位置：`internal/judge/judge.go` — `runBinary` / `runJavaCompiled`

---

## 项目结构

```
oj-platform/
├── cmd/
│   └── server/             # 应用入口
├── internal/
│   ├── database/           # 数据库初始化与迁移
│   ├── handlers/           # HTTP 请求处理器
│   ├── judge/              # 多语言编译执行引擎 + 沙盒
│   ├── middleware/         # JWT 认证、CORS、限流
│   ├── models/             # GORM 数据模型
│   ├── queue/              # Worker 池任务队列
│   ├── repository/         # 数据访问层
│   ├── routes/             # 路由注册
│   └── services/           # 业务逻辑
├── pkg/
│   └── config/             # 配置加载
├── scripts/
│   ├── import_leetcode.go  # 题库导入脚本（100道）
│   └── gen_testcases.go    # 50组测试用例生成器
├── web/                    # 前端静态文件 (HTML/CSS/JS)
├── docs/                   # 文档
├── PROBLEM_LIST.md         # 完整题单（100道Hot100）
├── config.yaml             # 应用配置
├── deploy.sh               # 一键部署脚本
├── Dockerfile
└── docker-compose.yml
```

---

## API 接口

详见 [docs/API.md](docs/API.md)。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/register` | 注册用户 |
| POST | `/api/v1/login` | 登录，返回 JWT |
| GET | `/api/v1/problems` | 题目列表 |
| GET | `/api/v1/problems/:id` | 题目详情 |
| POST | `/api/v1/submit` | 提交代码（限流 5/s） |
| GET | `/api/v1/submissions/:id` | 查询提交结果 |
| POST | `/api/v1/test` | 自定义测试（限流 5/s） |
| POST | `/api/v1/problems/import` | 管理员：JSON 导入题目 |
| GET | `/api/v1/leaderboard` | 全局排行榜 |
| GET | `/api/v1/problems/:id/leaderboard` | 单题排行榜 |
| GET | `/api/v1/admin/submissions` | 管理员：提交列表 |
| GET | `/api/v1/admin/submissions/:id` | 管理员：查看代码 |

---

## 管理员导入题目

```json
{
  "title": "题目名称",
  "description": "题目描述...",
  "difficulty": "Easy",
  "tags": "数组,哈希表",
  "time_limit": 5000,
  "memory_limit": 256,
  "function_template": "func solution() {\n    // Go 函数模板\n}",
  "driver_code": "package main\nimport \"fmt\"\nfunc main() { ... }",
  "test_cases": [
    {"input": "1 2 3\n4", "output": "0 1", "is_public": true}
  ]
}
```

---

## 架构

```
浏览器 → Gin HTTP 服务
           ├── RateLimit 中间件（令牌桶，100/5 req/s）
           ├── JWT Auth 中间件
           └── Handler → Service → 多语言 Judge 引擎
                                        (Go/C/C++/Java)
                          ↓
                     任务队列 (20 Workers)
                          ↓
              ulimit 沙盒 + /usr/bin/time 资源测量
                          ↓
                     SQLite (GORM)
```

---

## 配置

`config.yaml`：

```yaml
server:
  port: 8080
  mode: release

database:
  host: sqlite
  dbname: oj_platform

judge:
  go_path: /usr/local/go/bin/go
  gcc_path: /usr/bin/gcc
  gpp_path: /usr/bin/g++
  javac_path: /usr/bin/javac
  java_path: /usr/bin/java
  time_limit: 5000    # 毫秒
  memory_limit: 256   # MB（同时作为 ulimit -v 上限）
  worker_count: 20

jwt:
  secret: your-secret-key
  expire: 24          # 小时
```

---

## 贡献

欢迎参与贡献，详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

---

## 路线图

- [x] 完整 LeetCode Hot 100 题库（100 道）
- [x] 50 组测试用例（程序化生成）
- [x] 排行榜 / 提交历史
- [x] 管理后台（提交审查）
- [x] 自定义测试控制台（分栏 Input/Output）
- [x] 性能可视化（时间 & 内存进度条）
- [x] 多语言支持（C / C++ / Java）
- [x] JSON 导入自定义题目
- [x] 请求限流（令牌桶，100/5 req/s per IP）
- [x] 沙盒隔离（ulimit 资源限制 + 进程组 kill）
- [x] 一键部署脚本
- [ ] PostgreSQL 支持
- [ ] Docker 沙盒（每次提交独立容器，更强隔离）

---

## 许可证

MIT License — 详见 [LICENSE](LICENSE)。
