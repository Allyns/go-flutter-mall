# Go Flutter Mall Admin Web

这是 Go Flutter Mall 项目的 Web 管理后台，基于 Vue 3 和 Element Plus 开发，提供了商品管理、订单管理、用户查看等功能。

## 技术栈

本项目采用了主流的 Vue 3 生态技术栈，构建现代化、响应式的管理后台。

### 核心框架
*   **Web 框架**: [Vue 3](https://vuejs.org/) (Composition API, `<script setup>`)
*   **构建工具**: [Vite](https://vitejs.dev/) - 极速的开发服务器和构建工具。
*   **UI 组件库**: [Element Plus](https://element-plus.org/) - 基于 Vue 3 的组件库，提供丰富的管理端组件。

### 状态管理与路由
*   **状态管理**: [Pinia](https://pinia.vuejs.org/) - Vue 的官方状态管理库，替代 Vuex。
    *   `useAuthStore`: 管理登录状态和 Token。
*   **路由**: [Vue Router 4](https://router.vuejs.org/) - 管理页面导航和权限控制（路由守卫）。

### 网络与样式
*   **HTTP 客户端**: [Axios](https://axios-http.com/) - 处理 RESTful API 请求，配置了拦截器处理 Token 和错误。
*   **CSS 预处理器**: [Sass](https://sass-lang.com/) - 编写更高效的 CSS。

## 功能模块

管理后台包含以下核心功能：

### 1. 认证系统
*   **管理员登录**: 使用 JWT Token 进行身份验证。
*   **路由守卫**: 未登录用户会被强制重定向到登录页。

### 2. Dashboard
*   **概览**: 展示系统基础数据（预留）。

### 3. 商品管理 (`/products`)
*   **商品列表**: 分页展示商品，支持显示封面图、价格、库存等。
*   **商品编辑/添加**: 弹窗形式的表单，支持编辑商品详情、分类、库存和图片链接。
*   **评价查看**: 查看商品的评价列表（包含用户、评分、内容、时间）。
*   **删除商品**: 软删除商品。

### 4. 订单管理 (`/orders`)
*   **订单列表**: 展示所有用户订单，包含订单号、金额、状态、用户等信息。
*   **状态管理**: 管理员可修改订单状态（如发货、完成）。
*   **删除订单**: 删除异常订单。

### 5. 即时通讯 (`/chat`)
*   **客服系统**: (开发中) 管理员与用户的实时聊天界面。

### 6. 系统通知 (`/notifications`)
*   **通知管理**: (开发中) 查看和发送系统通知。

## 快速开始

### 依赖安装

```bash
cd admin-web
npm install
```

### 开发模式运行

```bash
npm run dev
```
服务默认运行在 `http://localhost:5173`。

### 生产构建

```bash
npm run build
npm run preview
```

## 目录结构

```
admin-web/
├── src/
│   ├── api/          # API 接口封装 (可选)
│   ├── assets/       # 静态资源
│   ├── components/   # 公共组件
│   ├── router/       # 路由配置
│   ├── stores/       # Pinia 状态管理
│   ├── views/        # 页面视图
│   │   ├── chat/     # 客服聊天
│   │   ├── orders/   # 订单管理
│   │   ├── products/ # 商品管理
│   │   ├── ...
│   ├── App.vue       # 根组件
│   └── main.js       # 入口文件
├── vite.config.js    # Vite 配置
└── package.json      # 项目依赖
```
