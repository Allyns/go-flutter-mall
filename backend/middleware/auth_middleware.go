package middleware

import (
	"net/http"
	"strings"

	"go-flutter-mall/backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 是 JWT 认证中间件
// 它会拦截请求，检查 Authorization 头中的 Token 是否有效
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // 终止后续处理
			return
		}

		// 2. 解析 Bearer Token
		// 格式通常为: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 验证 Token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 4. 将用户 ID 存入上下文
		// 后续的控制器可以通过 c.Get("userID") 获取当前登录用户的 ID
		c.Set("userID", claims.UserID)

		// 继续处理请求
		c.Next()
	}
}
