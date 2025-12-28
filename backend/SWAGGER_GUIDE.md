# Swagger API 文档使用指南

本文档详细说明了 Go-Flutter Mall 项目中如何使用 Swagger (Gin-Swagger) 进行 API 文档的编写、生成和查看。

## 1. 简介

Swagger 是一套用于设计、构建、记录和使用 RESTful Web 服务的工具。本项目使用 `swaggo/swag` 库，通过在代码中编写注释来自动生成 OpenAPI (Swagger 2.0) 规范文档。

### 1.1 核心组件
*   **swag**: 命令行工具，解析 Go 代码中的注释并生成 `docs/swagger.json` 和 `docs/swagger.yaml`。
*   **gin-swagger**: Gin 中间件，用于在 Web 页面中渲染 Swagger UI。

---

## 2. 安装与环境配置

### 2.1 安装 Swag CLI
确保你的环境已安装 Go，然后运行以下命令安装 `swag` 命令行工具：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

安装完成后，确保 `$GOPATH/bin` 在你的系统 `PATH` 环境变量中，以便可以直接运行 `swag` 命令。

---

## 3. 生成文档

每次修改 API 注释后，都需要重新生成文档。

### 3.1 生成命令
在 `backend` 根目录下执行：

```bash
swag init
```

执行成功后，`backend/docs` 目录下会生成/更新以下文件：
*   `docs.go`
*   `swagger.json`
*   `swagger.yaml`

---

## 4. 编写注释规范

Swagger 文档完全依赖于代码中的特定格式注释。

### 4.1 全局配置 (`main.go`)
在 `main.go` 的 `main` 函数上方配置项目信息：

```go
// @title           Go Flutter Mall API
// @version         1.0
// @description     This is the backend API for the Go Flutter Mall application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### 4.2 控制器注释 (`controllers/*.go`)
在每个 Handler 函数上方添加注释。

**示例：登录接口 (`controllers/auth_controller.go`)**

```go
// Login 处理用户登录请求
// @Summary      User Login
// @Description  Login with email/username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginInput     true  "Login Credentials"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Router       /auth/login [post]
func Login(c *gin.Context) {
    // ...
}
```

**关键注解说明**:
*   `@Summary`: 简短摘要。
*   `@Description`: 详细描述。
*   `@Tags`: API 分组（如 Auth, Product, Order）。
*   `@Accept`: 请求内容类型 (通常为 json)。
*   `@Produce`: 响应内容类型 (通常为 json)。
*   `@Param`: 参数定义。格式: `参数名 参数类型 数据类型 是否必填 "描述"`
    *   `path`: 路径参数 (如 `/products/{id}`)
    *   `query`: 查询参数 (如 `/products?page=1`)
    *   `body`: 请求体 (POST/PUT)
*   `@Success`: 成功响应。格式: `状态码 {数据类型} 返回结构体`
*   `@Failure`: 失败响应。
*   `@Router`: 路由路径和 HTTP 方法。
*   `@Security`: 安全认证 (如 `BearerAuth`)。

**示例：需要认证的接口**

```go
// @Security BearerAuth
// @Router /orders [post]
```

### 4.3 模型注释 (`models/*.go`)
Swagger 会解析结构体字段的 tag 来生成参数说明。

```go
type LoginInput struct {
    Email    string `json:"email" example:"user@example.com"` // 示例值
    Password string `json:"password" binding:"required" example:"123456"`
}
```

---

## 5. 访问与测试

### 5.1 启动服务
```bash
go run main.go
```

### 5.2 访问 Swagger UI
在浏览器中打开：
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### 5.3 在线调试
1.  在 Swagger UI 页面，点击具体的 API 接口。
2.  点击 **Try it out** 按钮。
3.  填写参数（如有）。
4.  点击 **Execute** 发送请求。
5.  查看 **Responses** 中的结果。

**认证测试**:
1.  先调用 `/auth/login` 接口获取 Token。
2.  点击页面顶部的 **Authorize** 按钮。
3.  在 Value 框中输入 `Bearer <your_token_string>`。
4.  点击 **Authorize** -> **Close**。
5.  现在可以测试带有锁图标（需要认证）的接口了。

---

## 6. 常见问题与最佳实践

1.  **文档未更新**: 必须重新运行 `swag init` 并重启 Go 服务，更改才会生效。
2.  **解析错误**: `swag init` 会严格检查注释格式。如果报错，请仔细检查报错信息中的行号和语法。
3.  **结构体引用**: 在 `@Success` 或 `@Param` 中引用结构体时，如果结构体在不同包，需要加上包名（如 `models.User`）。
4.  **隐藏字段**: 在结构体字段中使用 `swaggerignore:"true"` 可以让该字段不显示在文档中。

## 7. 常用注解速查表

| 注解 | 说明 | 示例 |
| :--- | :--- | :--- |
| `@Summary` | 接口标题 | Get User Profile |
| `@Description` | 接口详细描述 | Get detailed information of the current user |
| `@Tags` | 接口分组标签 | User |
| `@Accept` | 请求数据格式 | json, mpfd (multipart/form-data) |
| `@Produce` | 响应数据格式 | json |
| `@Param` | 请求参数 | id path int true "User ID" |
| `@Success` | 成功响应 | 200 {object} models.User |
| `@Failure` | 失败响应 | 400 {object} map[string]string |
| `@Router` | 路由定义 | /users/{id} [get] |
| `@Security` | 认证方式 | BearerAuth |
