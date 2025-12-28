# Go Flutter Mall Backend

这是 Go Flutter Mall 项目的后端服务，基于 Golang 和 Gin 框架开发，提供了一套完整的电商系统 API。

## 技术栈

本项目采用了现代化的 Go 语言技术栈，构建高性能、可扩展的后端服务。

### 核心框架与语言
*   **语言**: [Go 1.24+](https://go.dev/) - 高性能、并发友好的编程语言。
*   **Web 框架**: [Gin](https://github.com/gin-gonic/gin) - 轻量级、高性能的 HTTP Web 框架。

### 数据存储
*   **关系型数据库**: [PostgreSQL](https://www.postgresql.org/) - 用于存储核心业务数据（用户、商品、订单、评价等）。
*   **ORM 框架**: [GORM](https://gorm.io/) - 强大的 Go 对象关系映射库，支持 PostgreSQL。
*   **文档数据库**: [MongoDB](https://www.mongodb.com/) - 用于存储非结构化数据或高频写入数据（如聊天记录、日志）。
*   **缓存**: [Redis](https://redis.io/) - 用于缓存热点数据、会话管理等（`go-redis/v9`）。创建订单并发处理，待支付处理

### 消息队列与异步处理
*   **消息队列**: [Kafka](https://kafka.apache.org/) - 用于处理高吞吐量的消息传递和异步任务解耦（`IBM/sarama`）。消息队列定制通知，系统通知
*   **WebSocket**: [Gorilla WebSocket](https://github.com/gorilla/websocket) - 实现实时双向通信（用于聊天和实时通知）。

### 认证与安全
*   **身份认证**: [JWT (JSON Web Tokens)](https://github.com/golang-jwt/jwt) - 用于无状态的用户身份验证。
*   **加密**: `bcrypt` (via `golang.org/x/crypto`) - 用于用户密码的安全哈希存储。

### 文档与工具
*   **API 文档**: [Swagger](https://swagger.io/) - 自动生成交互式 API 文档 (`swaggo/gin-swagger`)。

## 详细文档

更多详细信息请参考以下文档：

*   [部署与启动指南](docs/project_docs/DEPLOY.md)
*   [技术栈详解](docs/project_docs/TECH_STACK.md)
*   [功能逻辑与实现细节](docs/project_docs/FEATURES.md)

## 功能模块

后端服务包含以下核心功能模块：

### 1. 用户系统
*   **用户注册/登录**: 支持用户名/密码注册和登录，颁发 JWT Token。
*   **管理员系统**: 独立的管理员账户体系，用于后台管理。
*   **个人信息管理**: 头像、收货地址管理。

### 2. 商品系统
*   **商品管理**: 商品的增删改查 (CRUD)。
*   **分类管理**: 商品分类层级。
*   **SKU 管理**: 支持多规格商品（颜色、尺寸等）。
*   **搜索**: 商品关键字搜索及搜索历史记录。

### 3. 交易系统
*   **购物车**: 添加商品、修改数量、删除商品。
*   **订单管理**:
    *   创建订单、支付模拟。
    *   订单状态流转（待支付 -> 待发货 -> 待收货 -> 已完成/售后）。
    *   管理员订单管理（发货、查看详情）。

### 4. 互动与反馈
*   **评价系统**: 用户购买商品后进行评分和评论，支持图片。
*   **即时通讯 (Chat)**:
    *   基于 WebSocket 的实时聊天。
    *   用户与客服（管理员）的一对一对话。
    *   消息持久化存储 (MongoDB)。
*   **系统通知**: 基于 WebSocket 的实时系统通知推送。

### 5. 运维与监控
*   **数据填充**: 提供 `scripts/seed.go` 脚本，快速初始化测试数据。
*   **API 文档**: 集成 Swagger UI，方便接口调试。

## 快速开始

### 前置要求
确保本地已安装并运行以下服务：
- PostgreSQL (默认端口 5432)
- MongoDB (默认端口 27017, 项目内置了启动脚本)
- Redis (默认端口 6379)
- Kafka (可选，如未安装相关功能可能受限)

### 运行步骤

1.  **启动 MongoDB** (如果使用项目内置):
    ```bash
    # 在 backend 目录下
    ./mongo_bin/bin/mongod --dbpath ./data/db --bind_ip 127.0.0.1 --port 27017
    ```

2.  **启动后端服务**:
    ```bash
    cd backend
    go mod tidy
    go run main.go
    ```
    服务默认运行在 `http://localhost:8080`。

3.  **初始化数据** (可选):
    ```bash
    go run scripts/seed.go
    ```
    此命令会重置数据库并填充测试用户、商品和订单数据。

4.  **访问接口文档**:
    打开浏览器访问: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
