# Go Flutter Mall App

这是 Go Flutter Mall 项目的移动端应用，基于 Flutter 开发，支持 iOS 和 Android。应用提供了完整的电商购物流程，包括商品浏览、购物车、下单支付、即时通讯等功能。

## 技术栈

本项目采用了 Flutter 生态中流行的现代化架构和库。

### 核心框架
*   **UI 框架**: [Flutter](https://flutter.dev/) (SDK ^3.9.2) - Google 的 UI 工具包，用于构建精美的原生应用。
*   **语言**: [Dart](https://dart.dev/)

### 状态管理与路由
*   **状态管理**: [Riverpod](https://riverpod.dev/) (`flutter_riverpod`, `hooks_riverpod`) - 编译时安全的响应式缓存和数据绑定框架。
*   **Hooks**: `flutter_hooks` - 简化 Widget 生命周期管理。
*   **路由**: [Go Router](https://pub.dev/packages/go_router) - 声明式的路由管理，支持深层链接和路由守卫。

### 网络与通信
*   **HTTP 客户端**: [Dio](https://pub.dev/packages/dio) - 强大的 Dart HTTP 客户端，支持拦截器、全局配置等。
*   **WebSocket**: `web_socket_channel` - 用于实现实时聊天和通知。

### 本地存储
*   **Key-Value 存储**: `shared_preferences` - 用于持久化存储简单的配置信息（如 Token、用户设置）。

## 功能模块

客户端包含以下核心功能模块：

### 1. 认证模块 (`features/auth`)
*   **登录/注册**: 用户注册和登录，JWT Token 管理。
*   **自动登录**: 基于本地存储的 Token 自动维持登录状态。

### 2. 首页与商品 (`features/home`, `features/product`)
*   **首页**: 展示推荐商品、分类入口。
*   **商品详情**:
    *   商品大图轮播。
    *   SKU 规格选择（颜色、尺寸等）。
    *   查看商品评价。
    *   加入购物车/立即购买。
*   **搜索**: 关键字搜索商品，支持搜索历史。

### 3. 购物车 (`features/cart`)
*   **购物车管理**: 查看购物车商品、修改数量、选中/取消选中、删除商品。
*   **结算入口**: 计算选中商品总价并跳转结算。

### 4. 订单与支付 (`features/order`)
*   **结算页**: 确认收货地址、商品清单、总金额。
*   **订单列表**: 查看不同状态的订单（待支付、待发货、待收货等）。
*   **收货地址管理**: 添加、编辑、删除收货地址。

### 5. 个人中心 (`features/profile`)
*   **个人信息**: 展示头像、用户名。
*   **功能入口**: 订单管理、地址管理、设置等。

### 6. 互动与消息 (`features/chat`, `features/notification`)
*   **客服聊天**: 与后台管理员进行实时文字聊天。
*   **系统通知**: 接收订单状态变更等系统通知。

## 快速开始

### 环境准备
请确保已安装 Flutter SDK 并配置好开发环境（Android Studio 或 VS Code）。

```bash
flutter doctor
```

### 运行

```bash
# 获取依赖
cd app
flutter pub get

# 运行 (连接模拟器或真机)
flutter run
```

### 目录结构

```
app/lib/
├── core/             # 核心基础库 (HTTP, Router, Utils)
├── features/         # 业务功能模块 (按功能划分)
│   ├── auth/         # 认证
│   ├── home/         # 首页
│   ├── product/      # 商品
│   ├── cart/         # 购物车
│   ├── order/        # 订单 & 地址
│   ├── chat/         # 聊天
│   ├── notification/ # 通知
│   └── profile/      # 个人中心
└── main.dart         # 入口文件
```

## 技术栈（补充）

### 质量与测试
- 代码规范：`flutter_lints`，规则见 [analysis_options.yaml](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/analysis_options.yaml)。
- 测试框架：`flutter_test`，运行 `flutter test`。

### 关键基础设施
- 统一路由配置：使用 GoRouter 并结合认证状态做重定向，见 [router.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/router/router.dart)。
- HTTP 封装与拦截器：单例封装自动注入认证头与日志，见 [http_client.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/http/http_client.dart)。
- 未读消息计数：统一管理消息角标，见 [unread_provider.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/providers/unread_provider.dart)。

## 功能说明（详细）

### 认证与访问控制
- 登录/注册与 Token 持久化（`shared_preferences`）。
- 路由守卫：未登录访问受限页面自动跳转登录；已登录访问登录/注册跳转首页。详见 [routerProvider 重定向](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/router/router.dart#L54-L75)。

### 商品与搜索
- 顶部搜索栏、分类筛选、搜索历史，见 [product_list_screen.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/features/product/screens/product_list_screen.dart)。
- 商品详情支持路由参数 `/product/:id`，加入购物车后可前往结算，见 [product_detail_screen.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/features/product/screens/product_detail_screen.dart)。

### 购物车
- 模型包含商品详情、数量、选中状态，见 [cart_item.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/features/cart/models/cart_item.dart)。
- 删除/更新本地状态并同步后端，见 [cart_provider.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/features/cart/providers/cart_provider.dart#L71-L83)。

### 订单与地址
- 结算页 `/checkout`；地址列表 `/addresses` 及子路由 `add`、`edit`，见 [router.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/router/router.dart#L26-L60)。
- 订单列表支持状态筛选，个人中心展示状态汇总，见 [profile_screen.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/features/profile/screens/profile_screen.dart)。

### 消息与聊天
- 客服聊天与系统通知；底部导航“消息”显示未读数，见 [home_screen.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/features/home/screens/home_screen.dart#L53-L69)。

## 启动项目文档（详细）

### 前置条件
- 安装 Flutter SDK 并完成环境配置（Android Studio 或 VS Code）。
- 至少一个可用的模拟器或已连接真机。
- 后端服务启动且可访问（默认基址见“后端连接配置”）。

```bash
flutter doctor
```

### 启动步骤

```bash
# 获取依赖
cd app
flutter pub get

# 运行 (连接模拟器或真机)
flutter run
```

### 常用平台命令
- iOS 模拟器：`flutter run -d ios`
- Android 模拟器：`flutter run -d android`
- macOS 桌面：`flutter run -d macos`（如已启用桌面支持）
- Web（可选）：`flutter run -d chrome`（未默认开启）

### 构建发布包
- Android APK：`flutter build apk`
- iOS（需签名）：`flutter build ios`

## 安装与环境搭建文档

### 安装 Flutter
- 访问 https://flutter.dev 获取最新安装包，或使用包管理器安装。
- 将 `flutter/bin` 添加到系统 `PATH`。
- 运行 `flutter doctor` 并根据提示安装缺失的依赖。

### macOS（iOS 开发）
- 安装 Xcode（App Store）。
- 安装 CocoaPods：`sudo gem install cocoapods` 或 `brew install cocoapods`。
- 在 `ios/` 下执行一次 `pod install`（通常由 Flutter 自动管理）。
- 在 Xcode 中接受证书与许可，配置团队签名。

### Android 开发
- 安装 Android Studio，安装 SDK、平台工具和至少一个虚拟设备（AVD）。
- 在 Android Studio 中启用 USB 调试，并连接真机或启动模拟器。
- 配置环境变量 `ANDROID_HOME`（如必要）。

### Windows/Linux/桌面支持（可选）
- 安装对应平台的构建工具（Windows：Visual Studio；Linux：CMake/GTK）。
- 运行 `flutter config --enable-macos-desktop`/`--enable-windows-desktop`/`--enable-linux-desktop` 开启桌面支持。

## 后端连接配置

HTTP 客户端基址在 [http_client.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/http/http_client.dart#L20-L41) 中配置：
- 默认使用局域网 IP：`http://192.168.5.165:8080/api`
- Android 模拟器请使用 `10.0.2.2` 访问宿主机服务。
- iOS 模拟器/macOS 使用 `localhost`。
- 真机需使用宿主机局域网 IP，并确保网络互通。

如需修改，请调整 `baseUrl` 并根据运行平台进行区分。

## 主题与外观

应用使用暗色主题（黑色+绿色）。主题定义在 [main.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/main.dart#L16-L55)，包括 AppBar、底部导航栏、卡片等配色。

## 路由与导航

统一路由表与重定向逻辑见 [router.dart](file:///Users/allyn/Documents/trae_projects/go-flutter-mall/app/lib/core/router/router.dart)：
- 初始路由 `/`
- 登录 `/login`、注册 `/register`
- 商品详情 `/product/:id`
- 结算 `/checkout`
- 聊天 `/chat`
- 地址列表 `/addresses`，子路由 `add` 与 `edit`
- 未登录访问受限路由时自动跳转登录，已登录访问登录/注册自动跳转首页

## 开发与调试

- 使用 VS Code 或 Android Studio 进行开发。
- 支持热重载与热重启。
- 查看网络日志：Dio 已开启详细日志拦截器。
- 代码静态检查：`flutter analyze`
- 单元/Widget 测试：`flutter test`

## 安全与隐私

- 认证 Token 存储于 `shared_preferences` 并通过请求头 `Authorization: Bearer <token>` 发送。
- 不在代码中硬编码敏感信息，后端地址需根据环境手动配置。

## 常见问题

- Android 模拟器无法访问后端：后端地址需使用 `10.0.2.2` 而非 `localhost`。
- iOS 构建失败：检查 Xcode 版本、签名证书与 `pod` 依赖。
- 真机网络不通：确保手机与开发机在同一局域网，并开放防火墙端口。
