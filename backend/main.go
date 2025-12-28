package main

import (
	"fmt"
	"log"
	"time"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/pkg/kafka"
	"go-flutter-mall/backend/pkg/scheduler"
	"go-flutter-mall/backend/pkg/websocket"
	"go-flutter-mall/backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

// main 是应用程序的入口点
// 它负责初始化数据库连接、配置 CORS、设置路由并启动 Web 服务器
func main() {
	// 1. 初始化数据库连接
	// 连接到 PostgreSQL 数据库并执行自动迁移以同步数据库架构
	config.ConnectDatabase()
	// 连接到 MongoDB
	config.ConnectMongoDB()
	// 连接到 Redis
	config.ConnectRedis()
	// 连接到 Kafka
	config.ConnectKafka()

	// 2. 初始化 Gin 路由引擎
	r := gin.Default()

	// 3. 配置 CORS (跨域资源共享)
	// 允许前端应用 (如 Flutter Web 或本地调试) 访问后端 API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许所有来源 (生产环境应限制为特定域名)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 3.5 初始化 WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 3.6 启动 Kafka 消费者和延时队列调度器
	kafka.StartConsumer()
	scheduler.StartScheduler()

	// 4. 设置路由
	// 注册所有的 API 路由组 (Auth, Product, Cart, Order 等)
	routes.SetupRoutes(r, hub)

	// 5. 启动服务器
	// 在 8080 端口监听请求
	fmt.Println("Server is running on port 8080")
	log.Fatal(r.Run(":8080"))
}
