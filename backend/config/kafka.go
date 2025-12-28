package config

import (
	"log"

	"github.com/IBM/sarama"
)

var KafkaProducer sarama.SyncProducer

// ConnectKafka 初始化 Kafka 连接
func ConnectKafka() {
	// Kafka 配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// 连接到 Kafka Broker (假设运行在 localhost:9092)
	// 如果无法连接，不会 Panic，而是打印日志并禁用 Kafka 功能
	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v. Kafka features will be disabled.", err)
		return
	}

	KafkaProducer = producer
	log.Println("Connected to Kafka successfully!")
}
