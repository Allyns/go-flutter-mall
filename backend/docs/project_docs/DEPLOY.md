# 项目部署与启动文档 (Deployment & Startup)

本文档详细介绍了 Go Flutter Mall 后端项目的部署环境要求、配置说明以及启动步骤。

## 1. 环境要求 (Prerequisites)

在运行本项目之前，请确保您的开发环境已安装以下软件：

*   **Go**: 版本 1.24.0 或更高。
*   **PostgreSQL**: 主要关系型数据库，用于存储用户、商品、订单等核心数据。
*   **MongoDB**: (可选) 用于日志、非结构化数据或特定功能（目前代码中已连接，若未启动会禁用相关功能）。
*   **Redis**: 用于缓存、分布式锁（订单并发控制）和延时队列。
*   **Kafka**: (可选) 用于异步消息通知（如订单创建后的消息推送）。

## 2. 配置文件 (Configuration)

项目的配置主要位于 `backend/config/` 目录下。在启动前，请根据您的本地环境修改相应的连接配置。

### 2.1 数据库配置 (`config/database.go`)
默认连接信息如下，请根据实际情况修改：
- Host: `localhost`
- Port: `5432`
- User: `postgres`
- Password: `postgres`
- DB Name: `go_flutter_mall`

### 2.2 MongoDB 配置 (`config/mongo.go`)
默认连接 URI: `mongodb://127.0.0.1:27017`
数据库名: `go_flutter_mall`

### 2.3 Redis 配置 (`config/redis.go`)
默认地址: `127.0.0.1:6379`
密码: 无
DB: `0`

### 2.4 Kafka 配置 (`config/kafka.go`)
默认 Broker: `localhost:9092`

## 3. 启动步骤 (Startup Steps)

### 3.1 获取依赖
在 `backend` 目录下执行：
```bash
go mod download
```

### 3.2 启动基础设施
确保 PostgreSQL, Redis, MongoDB (可选), Kafka (可选) 已启动。

**示例 (使用 Docker 启动 Redis 和 Mongo):**
```bash
docker run -d --name redis -p 6379:6379 redis
docker run -d --name mongo -p 27017:27017 mongo
```

### 3.3 运行项目
在 `backend` 目录下执行：
```bash
go run main.go
```

或者编译后运行：
```bash
go build -o server main.go
./server
```

项目启动成功后，默认运行在 `8080` 端口（Gin 默认端口，除非在环境变量或代码中另行指定）。

## 4. 数据库迁移 (Database Migration)
项目启动时 (`main.go` -> `config.ConnectDatabase()`) 会自动执行 GORM 的 `AutoMigrate`，自动创建或更新以下表结构：
- users
- products
- categories
- product_skus
- cart_items
- orders
- order_items
- addresses
- admin_users
- chat_messages
- notifications
- reviews

无需手动运行 SQL 脚本。

## 5. API 文档
项目启动后，可以访问 Swagger 文档查看所有 API 接口：
地址: `http://localhost:8080/swagger/index.html`

## 6. 常见问题
- **Redis 连接失败**: 订单创建功能依赖 Redis 分布式锁，如果 Redis 未连接，可能会遇到并发问题，建议开发环境也启动 Redis。
- **Kafka 连接失败**: 系统会打印日志并禁用 Kafka 相关功能（如异步通知），不影响主流程运行。
