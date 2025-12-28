package models

import "gorm.io/gorm"

// Review 商品评价
type Review struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"` // 关联用户信息
	ProductID uint   `json:"product_id"`
	OrderID   uint   `json:"order_id"`
	Content   string `json:"content"`         // 评价内容
	Rating    int    `json:"rating"`          // 评分 1-5
	Images    string `json:"images"`          // 评价图片(JSON array string)
	Status    int    `json:"status"`          // 0:隐藏 1:显示 (默认1)
	Reply     string `json:"reply,omitempty"` // 商家回复
}
