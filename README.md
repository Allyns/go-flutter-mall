# Go Flutter Mall

基于 Flutter 和 Go (Gin) 的全栈电商商城系统。

## 项目结构

- **app/**: Flutter 移动端应用 (Android/iOS)
- **admin-web/**: Vue3 + Vite 管理后台前端
- **backend/**: Go (Gin) 后端 API 服务

## API 文档

后端服务启动后，可访问 Swagger API 文档：

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

## 快速开始

请参考各子目录下的 `README.md` (如有) 或直接运行各部分代码。

### 后端 (Backend)

```bash
cd backend
go run main.go
```

### 移动端 (App)

```bash
cd app
flutter run
```

### 管理后台 (Admin Web)

```bash
cd admin-web
npm install
npm run dev
```
