# JWT (Access Token) 认证使用指南

本文档详细说明了 Go-Flutter Mall 项目中基于 JWT (JSON Web Token) 的用户认证机制，包括核心原理、代码实现、使用场景及客户端对接指南。

## 1. 简介与机制

本项目使用标准 **JWT** 协议进行无状态认证。
*   **Token 格式**: `Header.Payload.Signature`
*   **签名算法**: HMAC-SHA256 (HS256)
*   **有效期**: 24 小时
*   **载荷 (Payload)**: 包含 `user_id` 和标准注册声明 (iss, exp, iat)。

### 1.1 认证流程
1.  **登录**: 用户提交用户名/密码 -> 服务器验证通过 -> 生成 JWT Token -> 返回给客户端。
2.  **请求**: 客户端在后续请求的 `Authorization` Header 中携带 Token。
3.  **验证**: 服务器中间件拦截请求 -> 解析验证 Token -> 提取 `user_id` 注入上下文 -> 放行到控制器。

---

## 2. 核心代码实现

### 2.1 Token 生成与解析 (`backend/utils/jwt.go`)

该工具包负责 Token 的签发和验签。

```go
package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// 密钥 (生产环境应从环境变量读取)
var jwtSecret = []byte("your_super_secret_key_change_this_in_production")

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

// GenerateToken: 签发 Token
func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "go-flutter-mall",
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ValidateToken: 验证 Token
func ValidateToken(tokenString string) (*Claims, error) {
    // ... 解析逻辑，详见源码 ...
}
```

### 2.2 认证中间件 (`backend/middleware/auth_middleware.go`)

中间件用于保护需要登录的 API 路由。

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取 Header: Authorization: Bearer <token>
        authHeader := c.GetHeader("Authorization")
        // ... 格式校验 ...

        // 2. 解析 Token
        tokenString := strings.Split(authHeader, " ")[1]
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // 3. 将 UserID 注入上下文，供后续 Controller 使用
        c.Set("userID", claims.UserID)
        c.Next()
    }
}
```

---

## 3. 使用场景与示例

### 3.1 场景一：用户登录 (颁发 Token)

在 `AuthController.Login` 中，验证密码成功后颁发 Token。

**文件**: `backend/controllers/auth_controller.go`

```go
func Login(c *gin.Context) {
    // ... 验证账号密码 ...
    
    // 生成 Token
    token, _ := utils.GenerateToken(user.ID)

    // 返回给客户端
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  user,
    })
}
```

### 3.2 场景二：保护路由 (应用中间件)

在路由配置中，对需要保护的路由组使用 `AuthMiddleware`。

**文件**: `backend/routes/routes.go` (示例)

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()

    // 公开路由
    auth := r.Group("/api/auth")
    {
        auth.POST("/login", controllers.Login)
        auth.POST("/register", controllers.Register)
    }

    //受保护路由 (需携带 Token)
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware()) // 挂载中间件
    {
        api.GET("/auth/me", controllers.GetUserProfile) // 获取个人信息
        api.POST("/orders", orderController.CreateOrder) // 创建订单
    }
    
    return r
}
```

### 3.3 场景三：在控制器中获取当前用户

在受保护的控制器方法中，通过 `c.Get("userID")` 获取当前登录用户的 ID。

**文件**: `backend/controllers/auth_controller.go`

```go
func GetUserProfile(c *gin.Context) {
    // 从上下文获取 userID (注意类型断言)
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 使用 userID 查询数据库
    var user models.User
    config.DB.First(&user, userID)
    
    c.JSON(http.StatusOK, gin.H{"user": user})
}
```

---

## 4. 客户端对接指南

### 4.1 存储 Token
客户端（Flutter/Web）登录成功后，应将 `token` 安全存储：
*   **Web**: `localStorage` 或 `HttpOnly Cookie`。
*   **Flutter**: `flutter_secure_storage` 或 `SharedPreferences`。

### 4.2 发送请求
访问受保护接口时，必须在 HTTP Header 中添加 `Authorization` 字段。

**格式**:
```http
Authorization: Bearer <your_token_here>
```

**Postman 示例**:
1.  Tab 栏选择 **Auth**。
2.  Type 选择 **Bearer Token**。
3.  在 Token 输入框填入登录接口返回的字符串。

**Curl 示例**:
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1..." http://localhost:8080/api/auth/me
```

### 4.3 处理过期
如果 Token 过期或无效，接口将返回 `401 Unauthorized`。
客户端应捕获此错误，并执行以下操作之一：
1.  跳转回登录页让用户重新登录。
2.  (进阶) 如果实现了 Refresh Token 机制，尝试自动刷新 Token。

---

## 5. 安全最佳实践

1.  **HTTPS**: 必须在生产环境使用 HTTPS，防止 Token 在传输中被窃听。
2.  **密钥保护**: `jwtSecret` 绝对不能提交到代码仓库，应通过环境变量 (`os.Getenv("JWT_SECRET")`) 注入。
3.  **有效期**: Token 有效期不宜过长（如本项目设为 24h），以降低 Key 泄露风险。
4.  **Payload**: 不要在 Token Payload 中存放敏感信息（如密码），因为 Payload 只是 Base64 编码，任何人都可以解码查看。
