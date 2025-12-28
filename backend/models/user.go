package models

import (
	"gorm.io/gorm"
)

// User 表示系统中的用户
// 包含用户的基本信息和认证凭据
type User struct {
	gorm.Model        // 嵌入 GORM 模型，包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"uniqueIndex;not null" json:"username"` // 用户名，必须唯一且不能为空
	Email      string `gorm:"uniqueIndex;not null" json:"email"`    // 邮箱，必须唯一且不能为空
	Password   string `gorm:"not null" json:"-"`                    // 密码，存储哈希值，JSON 序列化时忽略
	Avatar     string `json:"avatar"`                               // 用户头像 URL
	Role       string `gorm:"default:'user'" json:"role"`           // 用户角色 (user/admin)，默认为 user
}

// Address 表示用户的收货地址
type Address struct {
	gorm.Model
	UserID        uint   `json:"user_id"`        // 关联的用户 ID
	ReceiverName  string `json:"receiver_name"`  // 收货人姓名
	Phone         string `json:"phone"`          // 联系电话
	Province      string `json:"province"`       // 省份
	City          string `json:"city"`           // 城市
	District      string `json:"district"`       // 区/县
	DetailAddress string `json:"detail_address"` // 详细地址
	IsDefault     bool   `json:"is_default"`     // 是否为默认地址
}
