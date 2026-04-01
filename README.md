# Go 后台开发模板

基于 Gin + GORM + PostgreSQL + JWT 的后台开发模板，包含完整的用户认证、CRUD 操作、Docker 支持和 Swagger 文档。

## 技术栈

- **Web 框架**: Gin
- **数据库**: PostgreSQL + GORM
- **认证**: JWT (golang-jwt/jwt)
- **日志**: Zap
- **配置**: Viper
- **文档**: Swagger (swaggo/swag)
- **容器化**: Docker + Docker Compose

## 特性

- ✅ RESTful API 设计
- ✅ JWT 认证授权
- ✅ 数据库迁移
- ✅ Swagger 文档
- ✅ 结构化日志
- ✅ 优雅关闭
- ✅ 健康检查（含数据库连接检查）
- ✅ 版本管理
- ✅ 环境变量配置
- ✅ Docker 支持
- ✅ CORS 支持

## 项目结构

```txt
server/
├── cmd/server/          # 应用入口
├── config/              # 配置管理
├── internal/            # 内部代码
│   ├── handler/         # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   └── router/          # 路由
├── pkg/                 # 公共工具
│   ├── database/        # 数据库连接
│   ├── logger/          # 日志工具
│   └── response/        # 统一响应
├── docs/                # Swagger 文档
├── config.yaml          # 配置文件
└── Makefile             # 构建脚本
```

## 快速开始

### 方式一：使用 Docker Compose（推荐）

一键启动应用和数据库：

```bash
docker-compose up -d
```

查看日志：

```bash
make docker-logs
```

停止服务：

```bash
make docker-down
```

### 方式二：本地开发

### 1. 安装依赖

```bash
make deps
```

### 2. 配置数据库

修改 `config.yaml` 中的数据库配置：

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: server_db
  sslmode: disable
```

确保 PostgreSQL 已启动并创建对应数据库：

```bash
createdb server_db
```

### 3. 生成 Swagger 文档

```bash
make swagger
```

### 4. 运行项目

```bash
make run
```

服务将在 `http://localhost:8080` 启动。

## 环境变量配置

支持通过环境变量覆盖 `config.yaml` 配置：

```bash
export SERVER_PORT=8080
export DATABASE_HOST=localhost
export DATABASE_PASSWORD=your_password
export JWT_SECRET=your-secret-key
```

完整环境变量列表：

- `SERVER_PORT` - 服务端口
- `SERVER_MODE` - 运行模式 (debug/release)
- `DATABASE_HOST` - 数据库地址
- `DATABASE_PORT` - 数据库端口
- `DATABASE_USER` - 数据库用户
- `DATABASE_PASSWORD` - 数据库密码
- `DATABASE_DBNAME` - 数据库名
- `DATABASE_SSLMODE` - SSL 模式
- `JWT_SECRET` - JWT 密钥
- `JWT_EXPIRE_HOURS` - Token 过期时间
- `LOG_LEVEL` - 日志级别
- `LOG_FILE` - 日志文件路径

## API 文档

启动服务后访问：`http://localhost:8080/swagger/index.html`

## API 端点

### 系统相关

- `GET /health` - 健康检查（含数据库连接状态）
- `GET /version` - 获取版本信息

### 认证相关

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 用户管理（需要 JWT Token）

- `GET /api/v1/users/me` - 获取当前用户信息
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户
- `PATCH /api/v1/users/:id/status` - 更新用户状态

## 使用示例

### 1. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456",
    "email": "admin@example.com"
  }'
```

### 2. 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

### 3. 使用 Token 访问受保护接口

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Makefile 命令

```bash
make help      # 显示帮助信息
make deps      # 安装依赖
make swagger   # 生成 Swagger 文档
make build     # 编译项目
make run       # 运行项目
make dev       # 生成文档并运行
make clean     # 清理编译文件
make test      # 运行测试
make fmt       # 格式化代码

# Docker 相关
make docker-build  # 构建 Docker 镜像
make docker-up     # 启动 Docker Compose
make docker-down   # 停止 Docker Compose
make docker-logs   # 查看 Docker 日志
```

## Docker 部署

### 构建镜像（带版本信息）

```bash
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t server:1.0.0 .
```

### 生产环境部署

```bash
docker run -d \
  --name server \
  -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e DATABASE_PASSWORD=your-db-password \
  -e JWT_SECRET=your-jwt-secret \
  -e SERVER_MODE=release \
  server:1.0.0
```

## 配置说明

`config.yaml` 配置项：

```yaml
server:
  port: 8080 # 服务端口
  mode: debug # 运行模式: debug/release

database:
  host: localhost # 数据库地址
  port: 5432 # 数据库端口
  user: postgres # 数据库用户
  password: postgres # 数据库密码
  dbname: server_db # 数据库名
  sslmode: disable # SSL 模式

jwt:
  secret: your-secret-key-change-in-production # JWT 密钥（生产环境请修改）
  expire_hours: 24 # Token 过期时间（小时）

log:
  level: info # 日志级别: debug/info/warn/error
  file: logs/app.log # 日志文件路径
```

## 开发建议

1. **生产环境部署前**：
   - 修改 `jwt.secret` 为强密码
   - 设置 `server.mode` 为 `release`
   - 配置合适的日志级别
   - 使用环境变量管理敏感信息

2. **数据库迁移**：
   - 项目启动时会自动执行 `AutoMigrate`
   - 生产环境建议使用专业的迁移工具

3. **安全建议**：
   - 使用 HTTPS
   - 实施 API 限流
   - 添加请求参数验证
   - 定期更新依赖包
   - 不要在代码中硬编码密钥

4. **监控和运维**：
   - 使用 `/health` 端点进行健康检查
   - 使用 `/version` 端点查看版本信息
   - 配置日志收集和分析
   - 设置告警机制

## License

MIT
