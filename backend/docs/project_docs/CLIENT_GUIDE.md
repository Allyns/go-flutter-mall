# 客户端综合说明文档

面向客户端（App/Web/Flutter）的一体化文档，涵盖功能说明、技术栈、安装与环境准备、项目启动与接入方式。阅读本文件即可完成后端对接与本地运行。

## 项目概览
- 后端服务基于 Go + Gin 提供 REST API 与 WebSocket 实时通信
- 默认服务地址: `http://localhost:8080`
- API 文档: `http://localhost:8080/swagger/index.html`
- API 基础路径: `/api`

## 技术栈
- 语言与框架
  - Go 1.24.x
  - Gin v1.11.0（路由、中间件、请求处理）
  - Gorilla WebSocket v1.5.3（实时通信）
- 数据与持久化
  - PostgreSQL（核心业务数据）
  - GORM v1.31.1（ORM，自动迁移、事务、预加载）
  - MongoDB（可选，用于非结构化数据或聊天/日志）
- 缓存与消息
  - Redis + go-redis v9.17.2（缓存、分布式锁、延时队列）
  - Kafka + Sarama v1.46.3（异步消息事件，如订单创建）
- 安全与文档
  - JWT v5.3.0（认证）
  - Swagger（自动生成 API 文档）

## 安装与环境准备（macOS）
以下为本机安装方案；也可使用 Docker 运行依赖服务。

### 使用 Homebrew 安装
```bash
brew install go
brew install postgresql
brew install redis
brew install mongodb-community@7
brew install kafka
```
启动服务（如使用 brew services）：
```bash
brew services start postgresql
brew services start redis
brew services start mongodb-community@7
brew services start kafka
```

初始化 PostgreSQL（与项目默认配置保持一致）：
```bash
# 设置 postgres 用户密码为 postgres
psql -U postgres -c "ALTER USER postgres WITH PASSWORD 'postgres';"
# 创建数据库
createdb -U postgres go_flutter_mall
```

### 使用 Docker 安装（可选）
```bash
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
docker exec -it postgres psql -U postgres -c "CREATE DATABASE go_flutter_mall;"

docker run -d --name redis -p 6379:6379 redis
docker run -d --name mongo -p 27017:27017 mongo
docker run -d --name kafka -p 9092:9092 --env KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
  bitnami/kafka:latest
```
> Kafka 为可选，未启动时系统会降级禁用相关异步通知功能。

## 启动项目
在仓库根目录进入 `backend`：
```bash
cd backend
go mod download
go run main.go
```
- 默认监听端口: `8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- 可选数据初始化：
```bash
go run scripts/seed.go
```

## 客户端接入指南
### 认证与鉴权
- 登录成功后返回 JWT Token
- 客户端请求需在 Header 中携带：
```
Authorization: Bearer <token>
```
- 获取当前用户信息：`GET /api/auth/me`

### WebSocket
- 连接地址：`ws://localhost:8080/api/ws`
- 用于聊天与系统通知的实时推送

## 功能说明（模块与接口）
以下为核心模块的功能与接口入口，详细参数请参考 Swagger。

### 1. 用户与认证
- 注册：`POST /api/auth/register`
- 登录：`POST /api/auth/login`
- 管理员登录：`POST /api/auth/admin/login`
- 我的信息：`GET /api/auth/me`（需授权）

### 2. 商品
- 列表：`GET /api/products`
- 详情：`GET /api/products/:id`
- 评价列表：`GET /api/products/:id/reviews`
- 管理（需后台权限）：创建/更新/删除

### 3. 购物车（需授权）
- 查询：`GET /api/cart`
- 数量：`GET /api/cart/count`
- 添加：`POST /api/cart`
- 更新：`PUT /api/cart/:id`
- 删除：`DELETE /api/cart/:id`

### 4. 订单（需授权）
- 创建：`POST /api/orders`（从选中购物车项生成订单）
- 列表：`GET /api/orders`
- 数量统计：`GET /api/orders/counts`
- 详情：`GET /api/orders/:id`
- 支付：`POST /api/orders/:id/pay`
- 确认收货：`PUT /api/orders/:id/receipt`
- 评价：`POST /api/orders/:id/review`
- 售后：`POST /api/orders/:id/after-sales`
- 管理员接口：更新状态、删除订单、查询全部

订单状态码：
- 0: 待支付
- 1: 已支付
- 2: 待收货
- 3: 待评价
- 4: 已完成
- 5: 售后中

并发与库存控制（创建订单）：
- Redis 分布式锁防止超卖（`lock:product:{id}`）
- 数据库条件更新：仅在 `stock >= quantity` 时扣减库存
- 成功后清空已购买购物车项
- 异步处理：
  - 发送 Kafka 订单事件（可选）
  - Redis 延时队列 30 分钟未支付自动取消
  - 创建站内通知

### 5. 聊天与通知（需授权）
- 聊天用户列表：`GET /api/chat/users`
- 消息记录：`GET /api/chat/messages/:userId`
- 标记已读：`PUT /api/chat/read`
- 发送系统消息：`POST /api/chat/notification`
- 通知列表：`GET /api/notifications`
- 标记已读：`PUT /api/notifications/:id/read`
- 未读数量：`GET /api/notifications/unread-count`
- 管理员通知查询：`/api/notifications/admin/*`

### 6. 搜索历史（需授权）
- 添加：`POST /api/search/history`
- 查询：`GET /api/search/history`
- 清空：`DELETE /api/search/history`

### 7. 地址（需授权）
- 列表：`GET /api/addresses`
- 创建：`POST /api/addresses`
- 更新：`PUT /api/addresses/:id`
- 删除：`DELETE /api/addresses/:id`

## 配置说明
默认配置位于源码中（开发环境友好），如需自定义：
- PostgreSQL DSN：`config/database.go`
- MongoDB：`config/mongo.go`
- Redis：`config/redis.go`
- Kafka：`config/kafka.go`

## 常见问题
- 未携带 Token：返回 401，请先登录并在请求头中添加 Authorization。
- Redis 未启动：订单并发控制与延时队列受限，建议在开发环境也启动。
- Kafka 未启动：系统自动禁用异步事件，不影响主流程。
- 跨域问题：服务已启用 CORS，允许常见方法与头，生产环境建议限制允许域名。

## 参考入口
- Swagger 文档：`/swagger/index.html`
- WebSocket：`/api/ws`
- 根健康检查：`/`

