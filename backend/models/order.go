package models

import (
	"time"

	"gorm.io/gorm"
)

// Order 表示用户的订单
// 包含订单状态、支付信息和总金额
type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"` // Override ID to ensure lowercase JSON
	CreatedAt time.Time      `json:"created_at"`           // Override for custom formatting if needed, but here just for key name
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderNo     string      `gorm:"uniqueIndex;not null" json:"order_no"` // 订单编号，唯一
	UserID      uint        `json:"user_id"`                              // 关联的用户 ID
	TotalAmount float64     `json:"total_amount"`                         // 订单总金额
	Status      int         `gorm:"default:0" json:"status"`              // 订单状态: 0-待支付, 1-待发货, 2-待收货, 3-待评价(已完成), 4-已取消, 5-售后中
	AddressID   uint        `json:"address_id"`                           // 收货地址 ID
	Address     Address     `json:"address"`                              // 收货地址快照 (简化处理，实际应复制地址信息)
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items"`      // 订单包含的商品项
}

// OrderItem 表示订单中的具体商品项
type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderID      uint    `json:"order_id"`      // 关联的订单 ID
	ProductID    uint    `json:"product_id"`    // 商品 ID
	ProductName  string  `json:"product_name"`  // 商品名称 (快照)
	ProductImage string  `json:"product_image"` // 商品图片 (快照)
	SKUID        uint    `json:"sku_id"`        // SKU ID
	SKUName      string  `json:"sku_name"`      // SKU 名称 (快照)
	Price        float64 `json:"price"`         // 购买时的单价
	Quantity     int     `json:"quantity"`      // 购买数量
}
