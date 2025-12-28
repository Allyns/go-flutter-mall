package config

import (
	"fmt"
	"log"

	"go-flutter-mall/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 是全局数据库连接实例
// 在整个应用程序中用于执行数据库操作
var DB *gorm.DB

// ConnectDatabase 初始化数据库连接
// 它使用 PostgreSQL 驱动程序连接到数据库，并执行自动迁移
func ConnectDatabase() {
	// 数据库连接字符串 (DSN)
	// 包含主机、用户、密码、数据库名称、端口和 SSL 设置
	// 注意: 在生产环境中，这些敏感信息应从环境变量中读取
	dsn := "host=localhost user=postgres password=postgres dbname=go_flutter_mall port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	// 打开数据库连接
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果连接失败，记录错误并终止程序
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库架构
	// GORM 会自动创建或更新表结构以匹配 Go 结构体定义
	// 这包括 User, Product, CartItem, Order, OrderItem 等模型
	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.ProductSKU{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Address{},
		&models.AdminUser{},
		&models.ChatMessage{},
		&models.Notification{},
		&models.Review{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 将连接实例赋值给全局变量 DB
	DB = database
	fmt.Println("Database connected and migrated successfully")
}
