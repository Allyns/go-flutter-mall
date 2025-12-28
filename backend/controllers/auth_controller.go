package controllers

import (
	"net/http"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"
	"go-flutter-mall/backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterInput 定义注册接口的请求参数
type RegisterInput struct {
	Username string `json:"username" binding:"required"`       // 用户名必填
	Email    string `json:"email" binding:"required,email"`    // 邮箱必填且格式需正确
	Password string `json:"password" binding:"required,min=6"` // 密码必填且至少6位
}

// LoginInput 定义登录接口的请求参数
type LoginInput struct {
	Email    string `json:"email"`                       // 邮箱 (与 Username 二选一)
	Username string `json:"username"`                    // 用户名 (与 Email 二选一)
	Password string `json:"password" binding:"required"` // 密码必填
}

// Register 处理用户注册请求
// POST /api/auth/register
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. 验证请求参数
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 检查邮箱是否已存在
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// 3. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 4. 创建新用户
	newUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 5. 创建默认收货地址
	defaultAddress := models.Address{
		UserID:        newUser.ID,
		ReceiverName:  newUser.Username,
		Phone:         "13800138000",
		Province:      "Beijing",
		City:          "Beijing",
		District:      "Chaoyang",
		DetailAddress: "Default Address",
		IsDefault:     true,
	}
	config.DB.Create(&defaultAddress)

	// 6. 返回成功响应
	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful", "user": newUser})
}

// GetUserProfile 获取当前用户信息
// GET /api/auth/me
func GetUserProfile(c *gin.Context) {
	// 从上下文获取 userID (由中间件设置)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
	})
}

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
	var input LoginInput

	// 1. 验证请求参数
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 查找用户 (支持邮箱或用户名登录)
	var user models.User
	query := config.DB

	if input.Email != "" {
		query = query.Where("email = ?", input.Email)
	} else if input.Username != "" {
		query = query.Where("username = ?", input.Username)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Username is required"})
		return
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 3. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 4. 生成 JWT Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 5. 返回 Token 和用户信息
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}
