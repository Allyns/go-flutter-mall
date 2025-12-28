# Go-Flutter 商城项目规划文档

## 1. 项目概述
本项目是一个功能完善的电商商城系统，采用前后端分离架构，注重性能与用户体验。
- **后端**: Go (Gin + GORM)
- **数据库**: PostgreSQL
- **客户端**: Flutter (Android/iOS)

## 2. 技术栈详细选型

### 2.1 后端 (Backend)
- **语言**: Go (Golang)
- **Web 框架**: Gin (高性能、轻量级)
- **ORM 框架**: GORM (PostgreSQL)，自动迁移、事务、预加载
- **数据库**: PostgreSQL（核心数据），MongoDB（聊天记录、日志、非结构化数据）
- **缓存/队列**: Redis（缓存、分布式锁、延时队列），Kafka（异步事件）
- **鉴权**: JWT (Access Token)，基于中间件的路由保护与角色控制
- **实时**: Gorilla WebSocket（聊天与系统通知）
- **配置**: Viper（环境配置管理）
- **文档**: Swagger/OpenAPI（代码注释生成，集成 Swagger UI）

### 2.2 数据库 (Database)
- **核心**: PostgreSQL
- **缓存**: Redis（用于购物车、Session、热点商品、分布式锁、延时队列）
- **文档**: MongoDB（用于聊天记录、日志、非结构化数据）

### 2.3 客户端 (Client)
- **框架**: Flutter
- **状态管理**: **Riverpod** (编译时安全、可测试性强) + **Flutter Hooks** (简化 Widget 生命周期管理)
- **网络请求**: Dio (拦截器、全局错误处理)
- **路由**: GoRouter (声明式路由)
- **本地存储**: SharedPreferences (简单配置) / Hive (本地缓存)
- **UI 组件库**: Material 3 设计风格
- **国际化**: Flutter Intl / easy_localization
- **主题**: 明暗主题切换，支持动态色（Android 12+）
- **架构建议**: 分层（data/domain/presentation），Repository + UseCase，统一错误码与异常映射
- **网络层细节**: Token 自动续期、重试策略、全局 Loading/Toast、统一日志埋点

### 2.4 Web 管理端 (Admin Web)
- **框架**: React / Vue 3（任选其一，推荐 Vue 3 + Vite）
- **UI 库**: Ant Design（React）/ Element Plus（Vue）
- **状态管理**: Zustand/Redux（React）或 Pinia（Vue）
- **路由**: React Router / Vue Router
- **网络**: Axios，集中式拦截器与错误处理
- **构建与部署**: Vite + CI/CD
- **架构建议**: 模块化（商品/订单/用户/通知），组件库二次封装，权限指令/高阶组件

## 3. 客户端功能详细规划

### 3.1 首页 (Home)
- **沉浸式顶栏**: 搜索框、消息入口。
- **轮播图 (Banner)**: 展示活动、热门商品。
- **金刚区 (Grid Nav)**: 核心分类快捷入口。
- **商品推荐流**: 猜你喜欢（瀑布流布局）。

### 3.2 分类 (Category)
- **双栏联动**: 左侧一级分类，右侧二级/三级分类网格。
- **筛选排序**: 价格、销量、新品筛选。

### 3.3 商品详情 (Product Detail)
- **商品展示**: 多图轮播、视频展示。
- **SKU 选择**: 规格选择器（颜色、尺寸等），库存联动。
- **评价系统**: 用户评分、图文评价列表。
- **底部操作栏**: 收藏、加入购物车、立即购买。

### 3.4 购物车 (Cart)
- **商品管理**: 数量增减、左滑删除、规格修改。
- **选中逻辑**: 单选、全选、店铺维度选择。
- **价格计算**: 实时计算总价、优惠金额。

### 3.5 结算与支付 (Checkout)
- **订单确认**: 收货地址选择、配送方式、备注。
- **支付方式**: 模拟 支付宝/微信 支付流程。

### 3.6 个人中心 (Profile)
- **订单管理**: 全部、待付款、待发货、待收货、待评价。
- **常用工具**: 收货地址管理、我的收藏、浏览历史。
- **设置**: 个人信息修改、深色模式切换、多语言支持。

### 3.7 即时通讯与通知 (Chat & Notification)
- **聊天**: 与客服/管理员一对一聊天，WebSocket 实时更新
- **消息存储**: MongoDB（或 Postgres）持久化
- **系统通知**: 订单状态变更、活动推送，支持未读计数与已读状态
- **实现要点**: 心跳与重连、离线消息补偿、未读合并、消息去重

### 3.8 搜索历史 (Search)
- **记录**: 保存用户搜索关键词
- **管理**: 查看、清空搜索历史
- **实现要点**: 关键词脱敏与频率限制，结合热门推荐

## 4. 数据库设计概览 (核心表)

- `users`: 用户基础信息
- `user_addresses`: 收货地址
- `categories`: 多级商品分类
- `products`: 商品主表
- `product_skus`: 商品规格库存表 (SKU)
- `cart_items`: 购物车
- `orders`: 订单主表
- `order_items`: 订单快照
- `comments`: 商品评价
- `notifications`: 系统通知
- `search_histories`: 搜索历史
- `admin_users`: 管理员账户
- `chat_messages`: 聊天消息

> 实际实现中，表结构与模型定义参考后端代码：用户、商品、SKU、购物车、订单、订单项、地址、评价、通知等。

## 5. 开发计划

### Phase 1: 基础架构 (Infrastructure)
- [x] 项目初始化 (Go mod, Flutter create)
- [ ] 后端数据库模型定义 (GORM) & 自动迁移
- [ ] 客户端基础封装 (Dio, Riverpod Provider, Router)

### Phase 2: 用户与商品 (Core Features)
- [ ] 用户注册/登录 (JWT)
- [ ] 首页 UI 与接口联调
- [ ] 分类与商品列表页
- [ ] 商品详情页 (含 SKU 选择)

### Phase 3: 交易闭环 (Transaction)
- [ ] 购物车逻辑 (本地 + 云端同步)
- [ ] 订单创建与状态流转
- [ ] 地址管理

### Phase 4: 完善与优化 (Polish)
- [ ] 评价系统
- [ ] 个人中心完善
- [ ] 性能优化与深色模式适配

## 6. 后端功能与实现详解

- **认证与授权**
  - 注册/登录发放 JWT，后端中间件校验 `Authorization: Bearer <token>`
  - 角色区分：`user` / `admin`，管理员接口需额外中间件
  - 接口参考：[routes.go](../backend/routes/routes.go)
- **商品与分类**
  - 商品 CRUD，SKU 规格与库存管理
  - 列表/详情/评价接口，支持分页与筛选
  - 模型参考：[models/product.go](../backend/models/product.go)
- **购物车**
  - 添加/更新/删除，选中状态，数量变更与合并逻辑
  - 存储于 Postgres，接口在 `/api/cart/*`
- **订单系统（核心）**
  - 事务开启：创建订单整体在事务内执行
  - 并发控制：Redis `SetNX` 以 `lock:product:{id}` 防止超卖
  - 库存扣减：仅当 `stock >= quantity` 条件满足时扣减，行级锁保证原子性
  - 生成订单与订单项，清空已购买购物车项
  - 异步处理：
    - Kafka 发送 `OrderCreated` 事件
    - Redis 延时队列 30 分钟未支付自动取消并回滚库存
    - 创建站内通知（Postgres）
  - 参考实现：[order_controller.go](../backend/controllers/order/order_controller.go)
- **聊天与通知**
  - WebSocket Hub 管理客户端连接与消息广播
  - 聊天记录与系统通知存储，未读计数与标记已读
  - 参考实现：`pkg/websocket/*` 与 `controllers/chat/*`、`controllers/notification/*`
- **搜索历史与地址管理**
  - 搜索关键词记录与维护
  - 地址 CRUD，默认地址设置
- **API 文档**
  - Swagger 注释生成，在线文档 `/swagger/index.html`

## 7. Web 管理端功能与实现规划

- **登录与权限**
  - 管理员登录，基于 JWT 与前端路由守卫
- **仪表盘**
  - 订单/用户/商品统计与图表
- **商品管理**
  - 商品与 SKU 的 CRUD、图片与库存管理
- **订单管理**
  - 查看订单详情、发货、状态流转（后台维度）
- **用户与通知**
  - 用户列表、禁用/启用
  - 系统通知创建与下发（与后端通知接口打通）
- **技术实现要点**
  - Axios 拦截器统一处理鉴权与错误
  - 表格分页、筛选、导出
  - 组件化与状态管理统一封装
  - 权限粒度控制（菜单、按钮级别），审计日志与操作记录

## 8. 环境与部署建议

- **开发环境**
  - 使用 Homebrew 或 Docker 启动 Postgres/Redis/Mongo/Kafka
  - 后端 `go run main.go`，前端 Flutter/Web 独立运行
- **测试环境**
  - 引入种子数据脚本，模拟用户与商品
  - E2E 联调（客户端与后端）
- **生产环境**
  - 数据库备份与监控
 - Redis/Kafka 高可用
  - 反向代理与 HTTPS
  - 环境配置分层（dev/staging/prod），灰度发布与回滚
 
## 9. 风险与优化方向
 
 - **库存并发与超卖**：Redis 锁 + 条件更新，必要时引入乐观锁
 - **消息一致性**：Kafka 开启幂等与重试策略，降级方案健壮
 - **性能优化**：热点数据缓存，慢查询优化，批量接口
 - **安全性**：Token 过期与刷新策略、接口速率限制、输入校验
 - **前端体验**：骨架屏、占位与错误提示统一规范

## 10. 环境搭建（macOS）

- 安装组件（Homebrew）：
  ```bash
  brew install go
  brew install postgresql
  brew install redis
  brew install mongodb-community@7
  brew install kafka
  ```
- 启动服务：
  ```bash
  brew services start postgresql
  brew services start redis
  brew services start mongodb-community@7
  brew services start kafka
  ```
- 初始化 PostgreSQL（与后端默认配置一致）：
  ```bash
  psql -U postgres -c "ALTER USER postgres WITH PASSWORD 'postgres';"
  createdb -U postgres go_flutter_mall
  ```
- Docker 方案：
  ```bash
  docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
  docker exec -it postgres psql -U postgres -c "CREATE DATABASE go_flutter_mall;"

  docker run -d --name redis -p 6379:6379 redis
  docker run -d --name mongo -p 27017:27017 mongo
  docker run -d --name kafka -p 9092:9092 --env KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 bitnami/kafka:latest
  ```
- 端口约定：
  - PostgreSQL: 5432
  - Redis: 6379
  - MongoDB: 27017
  - Kafka: 9092
  - 后端服务: 8080

## 11. 启动方式

### 11.1 后端启动
- 进入后端目录：
  ```bash
  cd backend
  go mod download
  go run main.go
  ```
- 编译运行：
  ```bash
  go build -o server main.go
  ./server
  ```
- 服务入口与端口：`http://localhost:8080`
- Swagger 文档：`http://localhost:8080/swagger/index.html`
- 依赖服务需已启动（Postgres/Redis/Mongo/Kafka），与 `config/` 中连接一致：
  - PostgreSQL DSN 在 `config/database.go`
  - Redis 在 `config/redis.go`
  - MongoDB 在 `config/mongo.go`
  - Kafka 在 `config/kafka.go`
- 初始化测试数据（可用时）：
  ```bash
  go run scripts/seed.go
  ```

### 11.2 客户端（Flutter）启动
- 进入 Flutter 项目目录：
  ```bash
  flutter pub get
  ```
- 配置后端地址（开发）：
  - 基础地址：`http://localhost:8080/api`
  - WebSocket 地址：`ws://localhost:8080/api/ws`
  - 在网络层或配置文件中设置 `baseUrl` 与 Token 注入（Authorization: Bearer）。
- 运行：
  ```bash
  flutter run -d ios   # 或 android
  flutter run -d chrome  # Flutter Web
  ```
- 网络层建议：
  - 全局拦截器：注入 Token、处理 401、重试策略
  - 统一错误码映射与异常处理
  - 日志与埋点

### 11.3 Web 管理端启动
- 进入 Web 管理端目录（React/Vue）：
  ```bash
  npm install
  npm run dev
  ```
- 环境变量：
  - `VITE_API_BASE_URL=http://localhost:8080/api`
  - Token 存储与路由守卫（登录态校验）
- 生产构建：
  ```bash
  npm run build
  ```
- 关键配置：
  - Axios 拦截器：Authorization 注入、错误统一处理
  - 权限控制：菜单/按钮级别
  - 组件与状态管理统一封装
