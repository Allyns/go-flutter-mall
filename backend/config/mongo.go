package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// ConnectMongoDB initializes the MongoDB connection
func ConnectMongoDB() {
	// 默认连接本地 MongoDB，端口 27017
	// 如果在 Docker 中运行或有密码，请修改 URI
	// 例如: "mongodb://user:password@localhost:27017"
	uri := "mongodb://127.0.0.1:27017"

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v. MongoDB features will be disabled.", err)
		return
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v. MongoDB features will be disabled.", err)
		return
	}

	MongoClient = client
	MongoDB = client.Database("go_flutter_mall")

	fmt.Println("Connected to MongoDB successfully!")
}
