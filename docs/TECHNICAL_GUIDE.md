# OJ Platform 技术实现文档

## 目录
- [系统架构](#系统架构)
- [技术栈详解](#技术栈详解)
- [核心模块实现](#核心模块实现)
- [API接口文档](#api接口文档)
- [数据库设计](#数据库设计)
- [部署指南](#部署指南)
- [性能优化](#性能优化)

---

## 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层                                │
│  HTML/CSS/JavaScript + Monaco Editor + WebSocket Client     │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTP/WebSocket
┌────────────────────▼────────────────────────────────────────┐
│                      API网关层                               │
│  Gin Router + Middleware (Auth, CORS, Rate Limit)          │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
┌───────▼──────┐ ┌──▼───────┐ ┌──▼──────────┐
│  用户服务    │ │ 题目服务  │ │  判题服务   │
│  UserService │ │ Problem   │ │ JudgeService│
└───────┬──────┘ └──┬───────┘ └──┬──────────┘
        │           │            │
        └───────────┼────────────┘
                    │
        ┌───────────▼───────────┐
        │      任务队列层        │
        │  TaskQueue + Workers  │
        └───────────┬───────────┘
                    │
        ┌───────────▼───────────┐
        │     判题引擎层         │
        │  Go Judge + Sandbox   │
        └───────────┬───────────┘
                    │
        ┌───────────▼───────────┐
        │      数据持久层        │
        │  SQLite/PostgreSQL    │
        └───────────────────────┘
```

### 请求流程

#### 1. 用户注册/登录流程
```
用户输入 → API验证 → 密码加密 → 数据库存储 → JWT生成 → 返回Token
```

#### 2. 代码提交流程
```
提交代码 → JWT验证 → 创建任务 → 任务队列 → Worker处理 →
编译执行 → 对比结果 → 更新数据库 → 返回状态
```

---

## 技术栈详解

### 后端技术栈

#### 1. Web框架 - Gin
**选择理由**：
- 高性能：基于httprouter，零内存分配路由
- 轻量级：无额外依赖，启动快
- 易扩展：中间件机制完善

**关键实现**：
```go
// 路由分组 + 中间件
v1 := r.Group("/api/v1")
v1.Use(middleware.AuthRequired())
```

#### 2. ORM - GORM
**选择理由**：
- 自动迁移：简化数据库管理
- 关联关系：支持外键、预加载
- 多数据库：PostgreSQL/SQLite/MySQL

**关键实现**：
```go
// 软删除
DeletedAt gorm.DeletedAt `gorm:"index"`

// 自动迁移
DB.AutoMigrate(&models.User{}, &models.Problem{})
```

#### 3. 配置管理 - Viper
**选择理由**：
- 多格式支持：YAML/JSON/TOML
- 环境变量：支持覆盖配置
- 热重载：配置文件监听

**配置结构**：
```yaml
server:
  port: 8080
  mode: debug

database:
  host: sqlite  # 开发环境使用SQLite

judge:
  worker_count: 20  # 并发worker数量
  time_limit: 5000  # 5秒超时
```

#### 4. 认证 - JWT
**实现细节**：
- 算法：HS256
- 过期时间：24小时
- Claims：user_id, username, exp

**Token验证流程**：
```
请求 → 提取Authorization头 → 解析Bearer Token →
验证签名 → 提取Claims → 注入Context
```

---

## 核心模块实现

### 1. 判题引擎 (Judge Engine)

#### 设计思路
```
代码提交 → 临时文件创建 → 编译 → 执行 → 结果对比 → 清理
```

#### 关键代码
```go
func (j *Judge) RunGo(code, input, expectedOutput string) *JudgeResult {
    // 1. 创建隔离的临时目录
    tmpDir, _ := os.MkdirTemp("", "oj_*")
    defer os.RemoveAll(tmpDir) // 自动清理

    // 2. 写入代码文件
    codeFile := filepath.Join(tmpDir, "main.go")
    os.WriteFile(codeFile, []byte(code), 0644)

    // 3. 编译
    binaryFile := filepath.Join(tmpDir, "main")
    compileCmd := exec.Command(j.goPath, "build", "-o", binaryFile, codeFile)

    if err := compileCmd.Run(); err != nil {
        return &JudgeResult{Status: "Compile Error"}
    }

    // 4. 执行（带超时控制）
    runCmd := exec.Command(binaryFile)
    runCmd.Stdin = strings.NewReader(input)

    start := time.Now()
    done := make(chan error, 1)
    go func() { done <- runCmd.Run() }()

    select {
    case <-time.After(time.Duration(j.timeLimit) * time.Millisecond):
        runCmd.Process.Kill()
        return &JudgeResult{Status: "Time Limit Exceeded"}
    case err := <-done:
        // 5. 结果对比
        output := strings.TrimSpace(stdout.String())
        if output == expectedOutput {
            return &JudgeResult{Status: "Accepted"}
        }
        return &JudgeResult{Status: "Wrong Answer"}
    }
}
```

#### 安全措施
1. **文件隔离**：每个任务独立临时目录
2. **资源限制**：超时自动终止进程
3. **自动清理**：defer确保资源释放
4. **无网络访问**：编译和执行不联网

### 2. 任务队列 (Task Queue)

#### 设计思路
```
生产者-消费者模式：
提交请求 → 任务入队 → Worker并发处理 → 结果回调
```

#### 关键实现
```go
type TaskQueue struct {
    tasks   chan *Task
    workers int
    wg      sync.WaitGroup
}

func (q *TaskQueue) Start(handler func(*Task) *TaskResult) {
    // 启动多个worker
    for i := 0; i < q.workers; i++ {
        go func() {
            for task := range q.tasks {
                result := handler(task)
                task.ResultChan <- result
            }
        }()
    }
}
```

#### 并发控制
- Worker数量：20（可配置）
- 队列大小：1000（防止内存溢出）
- 优雅关闭：WaitGroup确保任务完成

### 3. 用户认证系统

#### 密码安全
```go
// 注册时加密
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost, // cost=10
)

// 登录时验证
bcrypt.CompareHashAndPassword(
    []byte(user.Password),
    []byte(inputPassword),
)
```

#### JWT生成
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id":  user.ID,
    "username": user.Username,
    "exp":      time.Now().Add(24 * time.Hour).Unix(),
})

tokenString, _ := token.SignedString([]byte(secret))
```

---

## API接口文档

### 认证相关

#### 1. 用户注册
```http
POST /api/v1/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 2. 用户登录
```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### 3. 获取用户信息
```http
GET /api/v1/profile
Authorization: Bearer <token>

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

### 题目相关

#### 4. 获取题目列表
```http
GET /api/v1/problems?page=1&page_size=20

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "problems": [...],
    "page": 1,
    "pageSize": 20
  }
}
```

#### 5. 获取题目详情
```http
GET /api/v1/problems/1

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "problem": {...},
    "testCases": [...]  // 仅公开测试用例
  }
}
```

#### 6. 创建题目（需认证）
```http
POST /api/v1/problems
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "两数之和",
  "description": "...",
  "difficulty": "Easy",
  "tags": "数组,哈希表",
  "time_limit": 5000,
  "memory_limit": 256
}
```

### 判题相关

#### 7. 提交代码
```http
POST /api/v1/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "problem_id": 1,
  "code": "package main\n...",
  "language": "go"
}

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "status": "Pending",
    ...
  }
}
```

#### 8. 查询判题结果
```http
GET /api/v1/submissions/1
Authorization: Bearer <token>

Response 200:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "status": "Accepted",  // 或 "Wrong Answer"
    "time_used": 123,
    "result": ""
  }
}
```

---

## 数据库设计

### ER图
```
┌─────────────┐       ┌──────────────┐
│    users    │       │   problems   │
├─────────────┤       ├──────────────┤
│ id (PK)     │       │ id (PK)      │
│ username    │       │ title        │
│ email       │       │ description  │
│ password    │       │ difficulty   │
│ created_at  │       │ time_limit   │
│ deleted_at  │       │ memory_limit │
└──────┬──────┘       └──────┬───────┘
       │                     │
       │                     │
       │    ┌────────────────┘
       │    │
       │    │
┌──────▼────▼─────┐    ┌──────────────┐
│   submissions   │    │  test_cases  │
├─────────────────┤    ├──────────────┤
│ id (PK)         │    │ id (PK)      │
│ user_id (FK)    │    │ problem_id(FK)│
│ problem_id (FK) │    │ input        │
│ code            │    │ output       │
│ language        │    │ is_public    │
│ status          │    │ created_at   │
│ result          │    │ deleted_at   │
│ time_used       │    └──────────────┘
│ memory_used     │
│ created_at      │
│ deleted_at      │
└─────────────────┘
```

### 索引设计
```sql
-- 用户表索引
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 提交记录索引
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_problem_id ON submissions(problem_id);
CREATE INDEX idx_submissions_deleted_at ON submissions(deleted_at);

-- 测试用例索引
CREATE INDEX idx_test_cases_problem_id ON test_cases(problem_id);
```

---

## 部署指南

### 本地开发环境

#### 前置要求
- Go 1.21+
- SQLite（开发）/ PostgreSQL（生产）

#### 启动步骤
```bash
# 1. 克隆项目
git clone <repository-url>
cd oj-platform

# 2. 安装依赖
go mod download

# 3. 配置文件
cp config.yaml.example config.yaml
# 编辑config.yaml，设置数据库等

# 4. 运行
go run cmd/server/main.go
# 或
./start.sh
```

### Docker部署

#### 1. 构建镜像
```bash
docker build -t oj-platform:latest .
```

#### 2. Docker Compose
```bash
docker-compose up -d
```

#### 3. 查看日志
```bash
docker-compose logs -f app
```

### 生产环境部署

#### 1. 配置优化
```yaml
# config.yaml
server:
  port: 8080
  mode: release  # 生产模式

database:
  host: postgres
  port: 5432
  user: oj_user
  password: <strong-password>
  dbname: oj_platform
  sslmode: require

jwt:
  secret: <random-256-bit-secret>
  expire: 24

log:
  level: info
  format: json
```

#### 2. 环境变量
```bash
export DB_PASSWORD=<password>
export JWT_SECRET=<secret>
```

#### 3. 反向代理（Nginx）
```nginx
upstream oj_backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name oj.example.com;

    client_max_body_size 10M;

    location / {
        proxy_pass http://oj_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /ws {
        proxy_pass http://oj_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## 性能优化

### 1. 数据库优化
- 连接池配置
- 查询优化（避免N+1）
- 索引优化

### 2. 判题引擎优化
- Worker池复用
- 编译缓存
- 并发控制

### 3. API优化
- 响应压缩
- 缓存策略
- 限流保护

### 4. 监控指标
```go
// Prometheus指标
- http_request_duration_seconds
- judge_task_queue_size
- judge_task_processing_time
- active_goroutines
```

---

## 故障排查

### 常见问题

#### 1. 编译错误
**现象**：提交代码后返回Compile Error
**排查**：
```bash
# 查看编译输出
tail -f server.log | grep "Compile Error"
```

#### 2. 超时问题
**现象**：任务一直Pending
**排查**：
```bash
# 检查worker状态
curl http://localhost:8080/health

# 查看队列积压
# 添加监控端点
```

#### 3. 数据库连接失败
**现象**：启动时报错"failed to connect database"
**排查**：
```bash
# 检查数据库状态
docker ps | grep postgres

# 测试连接
psql -h localhost -U oj_user -d oj_platform
```

---

## 扩展开发

### 添加新语言支持

#### 1. 定义语言配置
```go
type LanguageConfig struct {
    Name        string
    CompileCmd  []string
    RunCmd      []string
    SourceFile  string
    BinaryFile  string
}
```

#### 2. 实现编译执行
```go
func (j *Judge) RunPython(code, input string) *JudgeResult {
    // Python是解释型语言，无需编译
    cmd := exec.Command("python3", "-c", code)
    // ...执行逻辑
}
```

#### 3. 注册到判题引擎
```go
func (j *Judge) Run(code, input, expected, language string) *JudgeResult {
    switch language {
    case "go":
        return j.RunGo(code, input, expected)
    case "python":
        return j.RunPython(code, input, expected)
    // ...
    }
}
```

### 添加WebSocket实时推送

#### 1. 集成Centrifuge
```go
import "github.com/centrifugal/centrifuge"

func setupWebSocket(r *gin.Engine) {
    node := centrifuge.New()
    r.GET("/ws", func(c *gin.Context) {
        // WebSocket处理
    })
}
```

#### 2. 推送判题结果
```go
func (s *JudgeService) updateSubmission(submissionID uint, result *TaskResult) {
    // 更新数据库
    // ...

    // 推送WebSocket消息
    s.wsClient.Publish("submission_"+userID, result)
}
```

---

## 安全加固

### 1. 代码沙箱
- 限制系统调用（seccomp）
- 限制文件访问（chroot）
- 限制网络访问

### 2. API安全
- Rate Limiting（限流）
- CORS配置
- SQL注入防护（GORM已防护）

### 3. 认证安全
- Token刷新机制
- 密码强度验证
- 防暴力破解

---

## 测试策略

### 单元测试
```bash
go test ./internal/judge -v
go test ./internal/services -v
```

### 集成测试
```bash
# 启动测试服务器
go test ./tests/integration -v
```

### 压力测试
```bash
# 使用wrk或ab
wrk -t10 -c100 -d30s http://localhost:8080/api/v1/problems
```

---

## 维护手册

### 日志管理
```bash
# 日志轮转配置
/var/log/oj-platform/*.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
}
```

### 数据备份
```bash
# SQLite备份
sqlite3 oj_platform.db ".backup backup.db"

# PostgreSQL备份
pg_dump oj_platform > backup.sql
```

### 监控告警
- Prometheus + Grafana
- 告警规则：
  - 错误率 > 5%
  - 响应时间 > 1s
  - 队列积压 > 100

---

## 版本历史

### v1.0.0 (2026-03-18)
- ✅ 基础架构搭建
- ✅ 用户认证系统
- ✅ Go代码判题引擎
- ✅ 任务队列系统
- ✅ RESTful API

### 计划功能
- [ ] 多语言支持（Python, Java）
- [ ] WebSocket实时推送
- [ ] 前端界面
- [ ] 题目管理后台
- [ ] 排行榜系统

---

## 联系方式

- 项目地址：https://github.com/your-org/oj-platform
- 问题反馈：https://github.com/your-org/oj-platform/issues
- 技术文档：https://docs.oj-platform.com
