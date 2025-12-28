# 技术栈文档 (Technology Stack)

Go Flutter Mall 后端项目采用现代化的 Go 语言生态构建，注重高性能、并发处理和可扩展性。

## 1. 核心语言与框架
*   **语言**: [Go](https://go.dev/) (Version 1.24.0)
    *   利用 Go 的 Goroutine 实现高并发处理。
*   **Web 框架**: [Gin](https://github.com/gin-gonic/gin) (v1.11.0)
    *   高性能 HTTP Web 框架。
    *   用于路由管理、中间件处理 (Auth, CORS)、请求绑定与响应。

## 2. 数据存储 (Persistence)
*   **关系型数据库**: [PostgreSQL](https://www.postgresql.org/)
    *   **ORM**: [GORM](https://gorm.io/) (v1.31.1)
    *   用于存储核心业务数据（用户、订单、商品）。
    *   特性：自动迁移、事务支持、预加载 (Preload)。
*   **文档型数据库**: [MongoDB](https://www.mongodb.com/)
    *   **Driver**: [mongo-driver](https://go.mongodb.org/mongo-driver) (v1.17.6)
    *   用于存储非结构化数据或特定业务场景（代码中已集成连接）。

## 3. 缓存与中间件 (Cache & Middleware)
*   **缓存/KV 存储**: [Redis](https://redis.io/)
    *   **Client**: [go-redis](https://github.com/redis/go-redis/v9) (v9.17.2)
    *   **用途**:
        *   **分布式锁**: 解决订单创建时的库存扣减并发问题 (`SetNX`)。
        *   **延时队列**: 订单超时未支付自动取消 (ZSet)。
*   **消息队列**: [Kafka](https://kafka.apache.org/)
    *   **Client**: [Sarama](github.com/IBM/sarama) (v1.46.3)
    *   **用途**: 异步消息通知（如 `OrderCreated` 事件），解耦业务逻辑。

## 4. 实时通信 (Real-time)
*   **WebSocket**: [Gorilla WebSocket](https://github.com/gorilla/websocket) (v1.5.3)
    *   用于实现即时聊天功能 (Chat Controller) 和实时消息推送。

## 5. 安全与认证 (Security & Auth)
*   **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt/v5) (v5.3.0)
    *   用于用户身份认证，生成和验证 Access Token。
*   **密码加密**: [bcrypt](golang.org/x/crypto)
    *   用于用户密码的安全哈希存储。

## 6. API 文档 (Documentation)
*   **Swagger**: [swag](https://github.com/swaggo/swag)
    *   通过代码注释自动生成 OpenAPI 规范文档。
    *   集成 `gin-swagger` 提供在线文档界面。

## 7. 工具与库 (Tools)
*   **UUID/NanoID**: 用于生成唯一标识符 (订单号等)。
*   **Cors**: `gin-contrib/cors` 处理跨域请求。
