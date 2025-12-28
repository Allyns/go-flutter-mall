package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/pkg/kafka"

	"github.com/redis/go-redis/v9"
)

const DelayQueueKey = "order:delay_queue"

// AddToDelayQueue 添加订单到延时队列
// delay: 延时时间 (例如 30 分钟)
func AddToDelayQueue(orderID uint, userID uint, delay time.Duration) error {
	ctx := context.Background()
	score := float64(time.Now().Add(delay).Unix())
	member := fmt.Sprintf("%d:%d", orderID, userID) // Store OrderID:UserID

	return config.RedisClient.ZAdd(ctx, DelayQueueKey, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// StartScheduler 启动调度器
func StartScheduler() {
	if config.RedisClient == nil {
		log.Println("Redis is disabled, scheduler will not start.")
		return
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		ctx := context.Background()

		for range ticker.C {
			now := float64(time.Now().Unix())

			// 获取已过期的任务 (Score <= Now)
			// 注意：如果 Redis 连接失败，这里会报错并阻塞，需要处理错误
			vals, err := config.RedisClient.ZRangeByScore(ctx, DelayQueueKey, &redis.ZRangeBy{
				Min:    "-inf",
				Max:    fmt.Sprintf("%f", now),
				Offset: 0,
				Count:  10, // 每次处理 10 个，防止阻塞
			}).Result()

			if err != nil {
				// 减少日志噪音，每隔一段时间打印一次
				// log.Printf("Failed to poll delay queue: %v", err)
				continue
			}

			if len(vals) > 0 {
				for _, val := range vals {
					// 解析 OrderID 和 UserID
					var orderID, userID uint
					fmt.Sscanf(val, "%d:%d", &orderID, &userID)

					// 从 Redis 中移除
					// 使用 ZRem 并检查返回值，确保只有一个消费者处理该任务
					removed, err := config.RedisClient.ZRem(ctx, DelayQueueKey, val).Result()
					if err != nil {
						log.Printf("Failed to remove from delay queue: %v", err)
						continue
					}
					if removed == 0 {
						// 已经被其他实例移除，跳过
						continue
					}

					// 发送超时事件到 Kafka
					// 这里的逻辑是：调度器只负责触发，具体的业务逻辑（取消订单）交给 Kafka 消费者
					// 这样实现了 延时（Redis）+ 异步处理（Kafka）的分离
					event := kafka.OrderEvent{
						OrderID:   orderID,
						UserID:    userID,
						EventType: "timeout",
					}

					// 如果 Kafka 不可用，这里可以直接调用处理函数降级处理
					if config.KafkaProducer != nil {
						kafka.SendOrderEvent(event)
					} else {
						// Fallback: manually trigger timeout logic if Kafka is missing
						// kafka.HandleOrderTimeout(event) // Need to export if used here
						log.Printf("Kafka missing, order %d timeout event lost (or implement fallback)", orderID)
					}
				}
			}
		}
	}()

	log.Println("Delay Queue Scheduler started...")
}
