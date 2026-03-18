# OJ Platform 部署指南

## 快速开始

### 本地开发环境

#### 1. 前置要求
- Go 1.21+
- SQLite（开发环境）

#### 2. 安装和运行
```bash
# 克隆项目
git clone <repository-url>
cd oj-platform

# 安装依赖
go mod download

# 运行
go run cmd/server/main.go

# 或使用部署脚本
chmod +x deploy.sh
./deploy.sh
```

#### 3. 访问应用
- 前端页面: http://localhost:8080
- API文档: http://localhost:8080/api/v1/
- 健康检查: http://localhost:8080/health

---

## 生产环境部署

### 方式一：Docker部署

#### 1. 构建镜像
```bash
docker build -t oj-platform:latest .
```

#### 2. 使用Docker Compose
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

#### 3. 配置说明
编辑 `docker-compose.yml`:
```yaml
services:
  app:
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=oj_user
      - DB_PASSWORD=<your-password>
      - DB_NAME=oj_platform
```

### 方式二：手动部署

#### 1. 编译
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./cmd/server
```

#### 2. 准备文件
```bash
mkdir -p deploy
cp bin/server deploy/
cp config.yaml deploy/
cp -r web deploy/
```

#### 3. 配置
编辑 `config.yaml`:
```yaml
server:
  mode: release  # 生产模式

database:
  host: <your-db-host>
  port: 5432
  user: <your-db-user>
  password: <your-db-password>
  dbname: oj_platform

jwt:
  secret: <random-256-bit-secret>
```

#### 4. 运行
```bash
cd deploy
./server
```

---

## Nginx反向代理配置

### 1. 创建Nginx配置
```nginx
# /etc/nginx/sites-available/oj-platform
upstream oj_backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name oj.example.com;

    client_max_body_size 10M;

    # 前端
    location / {
        proxy_pass http://oj_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # API
    location /api {
        proxy_pass http://oj_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 静态文件缓存
    location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg)$ {
        proxy_pass http://oj_backend;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

### 2. 启用配置
```bash
sudo ln -s /etc/nginx/sites-available/oj-platform /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. HTTPS配置（推荐）
```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d oj.example.com

# 自动续期
sudo certbot renew --dry-run
```

---

## Systemd服务配置

### 1. 创建服务文件
```bash
sudo nano /etc/systemd/system/oj-platform.service
```

内容：
```ini
[Unit]
Description=OJ Platform Service
After=network.target

[Service]
Type=simple
User=oj
WorkingDirectory=/opt/oj-platform
ExecStart=/opt/oj-platform/server
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### 2. 启动服务
```bash
sudo systemctl daemon-reload
sudo systemctl enable oj-platform
sudo systemctl start oj-platform
sudo systemctl status oj-platform
```

---

## 数据库配置

### PostgreSQL（生产环境推荐）

#### 1. 创建数据库
```sql
CREATE DATABASE oj_platform;
CREATE USER oj_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE oj_platform TO oj_user;
```

#### 2. 修改配置
```yaml
database:
  host: localhost
  port: 5432
  user: oj_user
  password: your_password
  dbname: oj_platform
  sslmode: disable
```

### SQLite（开发环境）

默认使用SQLite，无需额外配置：
```yaml
database:
  host: sqlite
```

---

## 监控和日志

### 1. 日志管理
```bash
# 查看实时日志
tail -f server.log

# 日志轮转配置
sudo nano /etc/logrotate.d/oj-platform
```

内容：
```
/opt/oj-platform/server.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 0640 oj oj
}
```

### 2. 监控指标
访问健康检查端点：
```bash
curl http://localhost:8080/health
```

---

## 性能优化

### 1. 数据库优化
- 启用连接池
- 添加索引
- 定期VACUUM（PostgreSQL）

### 2. 应用优化
- 调整worker数量
- 启用GZIP压缩
- 静态文件缓存

### 3. 系统优化
```bash
# 增加文件描述符限制
ulimit -n 65535

# 内核参数优化
sudo sysctl -w net.core.somaxconn=65535
```

---

## 故障排查

### 1. 服务无法启动
```bash
# 检查端口占用
lsof -i:8080

# 检查日志
tail -f server.log

# 检查配置文件
cat config.yaml
```

### 2. 数据库连接失败
```bash
# 测试连接
psql -h localhost -U oj_user -d oj_platform

# 检查数据库状态
sudo systemctl status postgresql
```

### 3. 权限问题
```bash
# 修改文件所有者
sudo chown -R oj:oj /opt/oj-platform

# 添加执行权限
chmod +x /opt/oj-platform/server
```

---

## 备份和恢复

### 1. 数据备份
```bash
# SQLite
sqlite3 oj_platform.db ".backup backup.db"

# PostgreSQL
pg_dump oj_platform > backup_$(date +%Y%m%d).sql
```

### 2. 数据恢复
```bash
# SQLite
cp backup.db oj_platform.db

# PostgreSQL
psql oj_platform < backup_20260318.sql
```

---

## 安全加固

### 1. 防火墙配置
```bash
# 只允许必要端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. JWT密钥
生成强随机密钥：
```bash
openssl rand -base64 32
```

### 3. 定期更新
```bash
# 更新依赖
go get -u
go mod tidy

# 重新编译部署
./deploy.sh
```

---

## 扩展阅读

- [技术实现文档](./docs/TECHNICAL_GUIDE.md)
- [API文档](./docs/API.md)
- [开发指南](./docs/DEVELOPMENT.md)
