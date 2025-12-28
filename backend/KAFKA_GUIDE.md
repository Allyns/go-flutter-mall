# Apache Kafka 使用指南

本文档详细说明了 Go-Flutter Mall 项目中 **Kafka** 消息队列的集成方式、配置说明、核心业务流程以及最佳实践。

## 1. 简介

Kafka 在本项目中作为核心的**异步消息中间件**，主要承担以下职责：
1.  **系统解耦**：将订单创建、支付成功等核心业务与通知发送、报表统计等辅助业务解耦。
2.  **流量削峰**：在高并发下单场景下，缓冲写入压力。
3.  **异步处理**：配合 Redis 延时队列，处理订单超时自动取消等耗时任务。

### 1.1 技术选型
*   **客户端库**: `github.com/IBM/sarama` (Go 社区最成熟的 Kafka 客户端)
*   **版本兼容**: 支持 Kafka 2.0+

---

## 2. 环境与连接配置

### 2.1 安装依赖
```bash
go get -u github.com/IBM/sarama
```

### 2.2 连接初始化 (`backend/config/kafka.go`)

系统启动时（`main.go`）会调用 `config.ConnectKafka()` 初始化全局生产者。

```go
package config

import (
    "github.com/IBM/sarama"
    "log"
)

var KafkaProducer sarama.SyncProducer

func ConnectKafka() {
    config := sarama.NewConfig()
    
    // 生产者配置
    config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有 ISR 副本确认，数据最安全
    config.Producer.Return.Successes = true          // 成功交付也要返回，便于同步发送
    config.Producer.Retry.Max = 5                    // 失败重试次数

    // 连接 Broker (生产环境应从配置读取)
    brokers := []string{"localhost:9092"}
    
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        log.Printf("Failed to connect to Kafka: %v. Features disabled.", err)
        return
    }

    KafkaProducer = producer
}
```

---

## 3. 消息模型与生产者

### 3.1 消息结构定义
为了规范消息格式，我们定义了统一的事件结构体 (`backend/pkg/kafka/kafka_service.go`)。

```go
type OrderEvent struct {
    OrderID   uint   `json:"order_id"`
    UserID    uint   `json:"user_id"`
    EventType string `json:"event_type"` // 枚举值: "created", "paid", "timeout", "cancelled"
}
```

### 3.2 发送消息 (Producer)
使用同步生产者 (`SyncProducer`) 确保消息发送成功后再返回，适用于订单等关键业务。

```go
// backend/pkg/kafka/kafka_service.go

func SendOrderEvent(event OrderEvent) error {
    if config.KafkaProducer == nil {
        return fmt.Errorf("Kafka not initialized")
    }

    // 1. 序列化
    bytes, _ := json.Marshal(event)

    // 2. 构建消息
    msg := &sarama.ProducerMessage{
        Topic: "order-events",
        Value: sarama.StringEncoder(bytes),
        // Key: sarama.StringEncoder(fmt.Sprint(event.OrderID)), // 可选：设置 Key 保证同一订单有序
    }

    // 3. 发送
    partition, offset, err := config.KafkaProducer.SendMessage(msg)
    if err != nil {
        log.Printf("Send Kafka message failed: %v", err)
        return err
    }
    
    return nil
}
```

---

## 4. 消费者与业务处理

消费者在后台 Goroutine 中运行，监听特定 Topic 的消息。

### 4.1 启动消费者
当前实现为简单的单分区消费者。生产环境建议使用 **Consumer Group** 以支持多实例负载均衡。

```go
// backend/pkg/kafka/kafka_service.go

func StartConsumer() {
    // ... 初始化 Consumer ...
    
    // 消费 "order-events" Topic 的 0 号分区
    partitionConsumer, _ := consumer.ConsumePartition("order-events", 0, sarama.OffsetNewest)

    go func() {
        for msg := range partitionConsumer.Messages() {
            handleMessage(msg)
        }
    }()
}
```

### 4.2 消息处理逻辑 (Router)
根据 `EventType` 分发处理逻辑。

```go
func handleMessage(msg *sarama.ConsumerMessage) {
    var event OrderEvent
    json.Unmarshal(msg.Value, &event)

    switch event.EventType {
    case "timeout":
        handleOrderTimeout(event) // 处理订单超时
    case "created":
        // 发送创建通知
    case "paid":
        // 更新统计数据
    }
}
```

---

## 5. 核心业务场景：订单超时取消

这是 Kafka 在本项目中最典型的应用场景，结合了 Redis 延时队列和 Kafka 异步处理。

**流程**:
1.  用户下单，Redis 写入延时任务 (Score = 当前时间 + 30分钟)。
2.  Scheduler 轮询 Redis，发现到期任务，发送 `timeout` 事件到 Kafka。
3.  Kafka Consumer 收到消息，执行 `handleOrderTimeout`。

**代码实现 (`handleOrderTimeout`)**:

```go
func handleOrderTimeout(event OrderEvent) {
    // 1. 幂等性检查：查询数据库，确认订单状态仍为 "待支付" (0)
    var order models.Order
    if err := config.DB.First(&order, event.OrderID).Error; err != nil {
        return
    }

    if order.Status == 0 {
        tx := config.DB.Begin() // 开启事务

        // 2. 更新订单状态为 "已取消" (-1)
        if err := tx.Model(&order).Update("status", -1).Error; err != nil {
            tx.Rollback(); return
        }

        // 3. 恢复库存 (关键步骤)
        var orderWithItems models.Order
        tx.Preload("Items").First(&orderWithItems, order.ID)
        for _, item := range orderWithItems.Items {
             tx.Model(&models.Product{}).
                Where("id = ?", item.ProductID).
                UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
        }

        tx.Commit()

        // 4. 创建站内信通知
        createNotification(event.UserID, "订单已取消", "因超时未支付...")
    }
}
```

---

## 6. 生产环境建议

1.  **Consumer Group**: 目前使用 `ConsumePartition` 仅适合单机开发。生产环境必须使用 `sarama.ConsumerGroup`，以便多个后端实例共同消费，自动重平衡 (Rebalance)。
2.  **Graceful Shutdown**: 确保在服务停止时调用 `producer.Close()` 和 `consumer.Close()`，避免消息丢失或重复消费。
3.  **消息幂等性**: 消费者必须设计为幂等的。网络波动可能导致消息重复投递，务必在业务逻辑中检查状态（如 `if order.Status == 0`）。
4.  **死信队列 (DLQ)**: 对于处理失败的消息（如 JSON 解析错误、数据库异常），应记录到专门的死信 Topic 或日志中，避免阻塞主队列。
