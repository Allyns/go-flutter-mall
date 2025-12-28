package models

import "gorm.io/gorm"

// AdminUser 管理员用户
type AdminUser struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `json:"-"`    // 存储哈希后的密码
	Role     string `json:"role"` // admin, support, stock_manager
	Avatar   string `json:"avatar"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	gorm.Model
	SenderID   uint   `json:"sender_id"`   // 发送者 ID (如果是客服，则是 AdminUser ID；如果是用户，则是 User ID)
	ReceiverID uint   `json:"receiver_id"` // 接收者 ID
	SenderType string `json:"sender_type"` // user, admin
	Content    string `json:"content"`
	IsRead     bool   `json:"is_read"`
	Type       string `json:"type"` // text, image
}
