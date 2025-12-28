# Gin Web Framework 使用指南

本文档详细说明了 Go-Flutter Mall 项目中使用的 Web 框架 **Gin** 的核心概念、路由设计、中间件机制及常用 API 开发模式。

## 1. 简介

Gin 是一个用 Go (Golang) 编写的 HTTP Web 框架。它以高性能（基于 httprouter）、中间件支持和简洁的 API 著称。

### 1.1 核心特性
*   **高性能**: 极快的路由匹配。
*   **中间件支持**: 方便地拦截请求进行日志记录、认证、CORS 处理等。
*   **JSON 验证**: 内置请求参数绑定与验证。
*   **路由组**: 方便管理 API 版本和模块。

---

## 2. 项目结构集成

在本项目中，Gin 的集成主要分布在以下文件：

*   **入口**: `backend/main.go` - 初始化引擎、全局中间件、启动服务。
*   **路由**: `backend/routes/routes.go` - 定义 API 路径与处理函数的映射。
*   **控制器**: `backend/controllers/*` - 具体的业务逻辑处理。
*   **中间件**: `backend/middleware/*` - 认证等通用处理。

### 2.1 初始化 (`main.go`)

```go
func main() {
    // 1. 创建默认引擎 (包含 Logger 和 Recovery 中间件)
    r := gin.Default()

    // 2. 注册全局中间件 (如 CORS)
    r.Use(cors.New(cors.Config{...}))

    // 3. 注册路由
    routes.SetupRoutes(r, hub)

    // 4. 启动服务 (默认 8080)
    r.Run(":8080")
}
```

---

## 3. 路由 (Routing)

路由定义在 `backend/routes/routes.go` 中，使用了 **路由组 (Group)** 来清晰地组织 API。

### 3.1 路由组与层级
```go
func SetupRoutes(r *gin.Engine) {
    // 创建 /api 根组
    api := r.Group("/api")
    {
        // 商品模块 (公开)
        products := api.Group("/products")
        {
            products.GET("", product.GetProducts)
            products.GET("/:id", product.GetProductDetail) // 路径参数
        }

        // 购物车模块 (需认证)
        // 在组级别应用 AuthMiddleware
        cartGroup := api.Group("/cart", middleware.AuthMiddleware())
        {
            cartGroup.POST("", cart.AddToCart)
        }
    }
}
```

### 3.2 路径参数
使用 `:` 获取 URL 路径中的参数。

*   **定义**: `/products/:id`
*   **请求**: `/products/123`
*   **获取**: `c.Param("id")` -> "123"

---

## 4. 请求处理 (Request Handling)

在 Controller 中，我们通常需要解析三种类型的参数：

### 4.1 路径参数 (Path Parameters)
用于资源定位。

```go
// GET /products/:id
func GetProductDetail(c *gin.Context) {
    id := c.Param("id") // 返回 string
    // ...
}
```

### 4.2 查询参数 (Query Parameters)
用于分页、搜索、过滤。

```go
// GET /products?page=1&search=apple
func GetProducts(c *gin.Context) {
    // 获取参数，若不存在则返回空字符串
    search := c.Query("search")
    
    // 获取参数，若不存在则使用默认值
    page := c.DefaultQuery("page", "1") 
}
```

### 4.3 请求体 (JSON Body)
用于 POST/PUT 请求的数据提交。通常配合 Struct Tag (`json`, `binding`) 使用。

**定义结构体**:
```go
type CreateProductInput struct {
    Name  string  `json:"name" binding:"required"` // 必填
    Price float64 `json:"price" binding:"required,gt=0"` // 必填且大于0
}
```

**解析**:
```go
func CreateProduct(c *gin.Context) {
    var input CreateProductInput
    
    // ShouldBindJSON 会自动根据 Content-Type 解析并校验字段
    if err := c.ShouldBindJSON(&input); err != nil {
        // 校验失败，返回 400 错误
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 业务逻辑...
}
```

---

## 5. 响应处理 (Response)

Gin 提供了多种响应格式，本项目主要使用 JSON。

### 5.1 JSON 响应
使用 `gin.H` (即 `map[string]interface{}`) 快速构建响应，或直接返回结构体。

```go
// 返回 Map
c.JSON(http.StatusOK, gin.H{
    "message": "Success",
    "id":      123,
})

// 返回结构体 (字段需大写导出)
c.JSON(http.StatusOK, userModel)

// 返回错误
c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
```

### 5.2 常用状态码
*   `http.StatusOK` (200): 成功。
*   `http.StatusCreated` (201): 资源创建成功。
*   `http.StatusBadRequest` (400): 参数错误。
*   `http.StatusUnauthorized` (401): 未登录。
*   `http.StatusForbidden` (403): 无权限。
*   `http.StatusNotFound` (404): 资源不存在。
*   `http.StatusInternalServerError` (500): 服务器内部错误。

---

## 6. 中间件 (Middleware)

中间件用于拦截请求，可以在处理具体业务逻辑 **之前** 或 **之后** 执行代码。

### 6.1 定义中间件
参考 `backend/middleware/auth_middleware.go`。

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // --- 请求处理前 ---
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatus(401) // 终止请求，不进入 Controller
            return
        }

        // 上下文传值
        c.Set("userID", 1001)

        c.Next() // 执行后续的处理函数 (Controller)

        // --- 请求处理后 ---
        // (例如计算请求耗时)
    }
}
```

### 6.2 获取中间件设置的值
在 Controller 中获取：

```go
func GetProfile(c *gin.Context) {
    // Get 返回 (interface{}, bool)
    val, exists := c.Get("userID")
    if !exists {
        // handle error
    }
    userID := val.(uint) // 类型断言
}
```

---

## 7. 常用开发模板

### 7.1 标准 Controller 模板

```go
// 1. 定义请求参数结构体
type UpdateProfileInput struct {
    Nickname string `json:"nickname" binding:"required"`
}

// 2. 编写处理函数
func UpdateProfile(c *gin.Context) {
    // A. 绑定参数
    var input UpdateProfileInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // B. 获取上下文信息 (如当前用户)
    userID, _ := c.Get("userID")

    // C. 业务逻辑 (调用 Service 或 DB)
    if err := db.UpdateUser(userID, input); err != nil {
        c.JSON(500, gin.H{"error": "Update failed"})
        return
    }

    // D. 返回结果
    c.JSON(200, gin.H{"status": "ok"})
}
```

### 7.2 获取 Header
```go
userAgent := c.GetHeader("User-Agent")
contentType := c.ContentType()
```

### 7.3 Cookie 操作
```go
// 设置 Cookie
c.SetCookie("label", "value", 3600, "/", "localhost", false, true)

// 获取 Cookie
cookie, err := c.Cookie("label")
```
