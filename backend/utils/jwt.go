package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret 是用于签名 Token 的密钥
// 注意: 在生产环境中，这应该从环境变量中读取，且必须足够复杂和保密
var jwtSecret = []byte("your_super_secret_key_change_this_in_production")

// Claims 定义了 Token 中包含的载荷信息
type Claims struct {
	UserID uint `json:"user_id"` // 用户 ID
	jwt.RegisteredClaims
}

// GenerateToken 为指定用户 ID 生成 JWT Token
// Token 有效期为 24 小时
func GenerateToken(userID uint) (string, error) {
	// 设置载荷
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
			Issuer:    "go-flutter-mall",                                  // 签发者
		},
	}

	// 创建 Token 对象，指定签名算法 (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥进行签名
	return token.SignedString(jwtSecret)
}

// ValidateToken 验证 Token 的有效性并解析载荷
func ValidateToken(tokenString string) (*Claims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否匹配
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 Token 是否有效并提取 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
