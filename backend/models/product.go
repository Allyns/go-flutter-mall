package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Category 表示商品分类
type Category struct {
	gorm.Model
	Name      string `json:"name"`       // 分类名称
	SortOrder int    `json:"sort_order"` // 排序权重
	Icon      string `json:"icon"`       // 分类图标
}

// Product 表示商品信息
// 包含商品的基本属性和关联的 SKU
type Product struct {
	gorm.Model
	Name        string         `json:"name"`                                          // 商品名称
	Description string         `json:"description"`                                   // 商品描述
	Price       float64        `json:"price"`                                         // 商品基础价格
	Stock       int            `json:"stock"`                                         // 总库存
	CoverImage  string         `json:"cover_image"`                                   // 封面图片 URL
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`                     // 商品轮播图列表 (PostgreSQL 数组类型)
	CategoryID  uint           `json:"category_id"`                                   // 分类 ID
	Status      int            `gorm:"default:1" json:"status"`                       // 商品状态: 1-上架, 0-下架
	SKUs        []ProductSKU   `gorm:"foreignKey:ProductID" json:"skus"`              // 关联的 SKU 列表
	Reviews     []Review       `gorm:"foreignKey:ProductID" json:"reviews,omitempty"` // 关联的评价列表
}

// ProductSKU 表示商品的库存量单位 (Stock Keeping Unit)
// 用于管理商品的不同规格 (如颜色、尺寸)
type ProductSKU struct {
	gorm.Model
	ProductID uint    `json:"product_id"` // 关联的商品 ID
	Name      string  `json:"name"`       // SKU 名称 (如 "红色 XL")
	Specs     string  `json:"specs"`      // 规格详情 JSON 字符串
	Price     float64 `json:"price"`      // SKU 价格
	Stock     int     `json:"stock"`      // SKU 库存
	Image     string  `json:"image"`      // SKU 图片
}

// CartItem 表示购物车中的一项
type CartItem struct {
	gorm.Model
	UserID    uint    `json:"user_id"`                               // 关联的用户 ID
	ProductID uint    `json:"product_id"`                            // 关联的商品 ID
	Product   Product `json:"product"`                               // 预加载的商品信息
	SKUID     uint    `gorm:"column:sku_id;default:0" json:"sku_id"` // 关联的 SKU ID (如果商品有规格)
	Quantity  int     `json:"quantity"`                              // 购买数量
	Selected  bool    `gorm:"default:true" json:"selected"`          // 是否选中
}
