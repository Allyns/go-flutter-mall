package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchHistory 搜索历史模型
// 存储在 MongoDB 中
type SearchHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    uint               `bson:"user_id" json:"user_id"`
	Keyword   string             `bson:"keyword" json:"keyword"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
