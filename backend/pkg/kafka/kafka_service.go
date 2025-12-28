package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

// OrderEvent 订单事件消息结构
type OrderEvent struct {
	OrderID   uint   `json:"order_id"`
	UserID    uint   `json:"user_id"`
	EventType string `json:"event_type"` // "created", "timeout", "cancelled"
}

// SendOrderEvent 发送订单事件到 Kafka
func SendOrderEvent(event OrderEvent) error {
	if config.KafkaProducer == nil {
		return fmt.Errorf("Kafka producer is not initialized")
	}

	bytes, _ := json.Marshal(event)

	msg := &sarama.ProducerMessage{
		Topic: "order-events",
		Value: sarama.StringEncoder(bytes),
	}

	partition, offset, err := config.KafkaProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

// StartConsumer 启动 Kafka 消费者
func StartConsumer() {
	if config.KafkaProducer == nil { // 简单检查 Kafka 是否可用
		log.Println("Kafka is disabled, consumer will not start.")
		return
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Printf("Failed to start Kafka consumer: %v", err)
		return
	}
	// defer consumer.Close() // Keep running

	partitionConsumer, err := consumer.ConsumePartition("order-events", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to consume partition: %v", err)
		return
	}

	log.Println("Kafka Consumer started...")

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				handleMessage(msg)
			case err := <-partitionConsumer.Errors():
				log.Printf("Consumer error: %v", err)
			}
		}
	}()
}

func handleMessage(msg *sarama.ConsumerMessage) {
	var event OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Received event: %s for Order %d", event.EventType, event.OrderID)

	// 根据事件类型处理
	switch event.EventType {
	case "timeout":
		// 处理订单超时取消逻辑
		handleOrderTimeout(event)
	case "created":
		// 可以在这里处理其他逻辑，例如发送邮件等
		log.Printf("Order %d created, waiting for payment...", event.OrderID)
	}
}

func handleOrderTimeout(event OrderEvent) {
	// 1. 检查订单状态，如果仍为 0 (待支付)，则取消订单
	var order models.Order
	if err := config.DB.First(&order, event.OrderID).Error; err != nil {
		log.Printf("Order %d not found", event.OrderID)
		return
	}

	if order.Status == 0 {
		// 开启事务
		tx := config.DB.Begin()

		// 更新状态为 -1 (已取消)
		if err := tx.Model(&order).Update("status", -1).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to cancel order %d: %v", event.OrderID, err)
			return
		}

		// 恢复库存
		// 需要预加载 Items 以获取 ProductID 和 Quantity
		var orderWithItems models.Order
		if err := tx.Preload("Items").First(&orderWithItems, order.ID).Error; err == nil {
			for _, item := range orderWithItems.Items {
				// 增加库存
				tx.Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
			}
		}

		tx.Commit()

		// 2. 发送消息通知
		notification := models.Notification{
			UserID:  event.UserID,
			Title:   "订单已取消",
			Content: fmt.Sprintf("您的订单 %s 因超时未支付已自动取消。", order.OrderNo),
			IsRead:  false,
		}
		config.DB.Create(&notification)

		log.Printf("Order %d cancelled due to timeout.", event.OrderID)
	} else {
		log.Printf("Order %d status is %d, skip cancellation.", event.OrderID, order.Status)
	}
}
