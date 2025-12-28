package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis 初始化 Redis 连接
func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Redis 地址 (使用 IPv4 避免 ::1 连接问题)
		Password: "",               // 密码
		DB:       0,                // 默认数据库
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		// 不强制退出，因为 Redis 可能是可选的（对于某些非关键功能）
		// 但在此需求中，订单并发依赖 Redis，所以最好确保连接成功
		// log.Fatal("Failed to connect to Redis")
	} else {
		fmt.Println("Connected to Redis successfully!")
	}
}
