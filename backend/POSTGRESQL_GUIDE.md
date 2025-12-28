# PostgreSQL 数据库使用指南

本文档详细说明了 Go-Flutter Mall 项目中 PostgreSQL 数据库的配置、模型设计、常用操作及最佳实践。

## 1. 简介与环境配置

本项目使用 **PostgreSQL** 作为核心关系型数据库，存储用户、商品、订单等结构化数据。ORM 框架采用 **GORM** 进行数据操作与迁移。

### 1.1 环境要求
- **版本**: PostgreSQL 13+
- **驱动**: `gorm.io/driver/postgres`
- **辅助库**: `github.com/lib/pq` (用于数组等高级类型支持)

### 1.2 连接配置 (`backend/config/database.go`)
目前数据库连接信息在代码中配置（生产环境建议迁移至环境变量）：

```go
dsn := "host=localhost user=postgres password=postgres dbname=go_flutter_mall port=5432 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

**参数说明**:
- `sslmode=disable`: 开发环境关闭 SSL，生产环境建议开启 (`require` 或 `verify-full`)。
- `TimeZone=Asia/Shanghai`: 确保写入的时间与业务时区一致，避免 8 小时误差。

---

## 2. 核心数据模型 (Schema)

项目使用 GORM 的 `AutoMigrate` 功能自动同步表结构。所有模型均嵌入了 `gorm.Model`，自动包含 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` (软删除)。

### 2.1 核心表概览

| 表名 (Table) | 模型 (Model) | 描述 | 关键字段 |
|---|---|---|---|
| `users` | `User` | 用户表 | `username`, `password` (hash), `role` |
| `products` | `Product` | 商品主表 | `price`, `stock`, `images` (Array), `category_id` |
| `product_skus` | `ProductSKU` | 商品规格 | `product_id`, `specs`, `price`, `stock` |
| `orders` | `Order` | 订单主表 | `order_number`, `user_id`, `total_amount`, `status` |
| `order_items` | `OrderItem` | 订单明细 | `order_id`, `product_id`, `quantity`, `price` |
| `cart_items` | `CartItem` | 购物车 | `user_id`, `sku_id`, `quantity`, `selected` |
| `addresses` | `Address` | 收货地址 | `user_id`, `name`, `phone`, `detail` |

### 2.2 特色类型支持
在 `Product` 模型中，使用了 PostgreSQL 特有的数组类型存储图片列表：

```go
import "github.com/lib/pq"

type Product struct {
    gorm.Model
    // ...
    // 对应数据库中的 text[] 类型
    Images pq.StringArray `gorm:"type:text[]" json:"images"` 
}
```

---

## 3. 常用开发场景与示例

### 3.1 创建数据 (Create)

**场景**: 上架一个新商品，包含多个 SKU 规格。

```go
func CreateProduct(db *gorm.DB) error {
    product := models.Product{
        Name:        "高性能机械键盘",
        Description: "全键热插拔，RGB背光",
        Price:       399.00,
        Stock:       100,
        CategoryID:  1,
        Images:      pq.StringArray{"https://img.example.com/1.jpg", "https://img.example.com/2.jpg"},
        SKUs: []models.ProductSKU{
            {Name: "红轴版", Price: 399.00, Stock: 50, Specs: `{"switch": "red"}`},
            {Name: "青轴版", Price: 399.00, Stock: 50, Specs: `{"switch": "blue"}`},
        },
    }
    
    // GORM 会自动关联创建 Product 和 ProductSKU
    return db.Create(&product).Error
}
```

### 3.2 查询数据 (Read)

**场景**: 获取商品详情，并预加载 SKU 和评价列表。

```go
func GetProductDetail(db *gorm.DB, productID uint) (*models.Product, error) {
    var product models.Product
    
    // 使用 Preload 加载关联表数据
    // 避免 N+1 查询问题
    err := db.Preload("SKUs").
              Preload("Reviews").
              First(&product, productID).Error
              
    return &product, err
}
```

**场景**: 分页查询商品列表。

```go
func ListProducts(db *gorm.DB, page, pageSize int) ([]models.Product, int64, error) {
    var products []models.Product
    var total int64
    
    offset := (page - 1) * pageSize
    
    // 1. 获取总数
    db.Model(&models.Product{}).Count(&total)
    
    // 2. 获取分页数据
    // Scope 用于封装通用的分页逻辑
    err := db.Limit(pageSize).Offset(offset).Find(&products).Error
    
    return products, total, err
}
```

### 3.3 事务处理 (Transaction)

**场景**: 用户下单，需要同时创建订单、创建订单项、扣减库存、清空购物车。这是典型的 ACID 事务场景。

```go
func CreateOrder(db *gorm.DB, userID uint, items []models.CartItem) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建订单
        order := models.Order{UserID: userID, Status: 1, OrderNumber: generateOrderNo()}
        if err := tx.Create(&order).Error; err != nil {
            return err // 返回错误会自动回滚
        }
        
        // 2. 处理每个商品
        for _, item := range items {
            // 2.1 扣减库存 (使用 SQL 表达式防止并发超卖)
            // UPDATE product_skus SET stock = stock - ? WHERE id = ? AND stock >= ?
            res := tx.Model(&models.ProductSKU{}).
                Where("id = ? AND stock >= ?", item.SKUID, item.Quantity).
                UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))
                
            if res.RowsAffected == 0 {
                return fmt.Errorf("库存不足: SKU %d", item.SKUID)
            }
            
            // 2.2 创建订单项
            orderItem := models.OrderItem{
                OrderID:   order.ID,
                ProductID: item.ProductID,
                Quantity:  item.Quantity,
                Price:     item.Product.Price, // 实际应取 SKU 价格
            }
            if err := tx.Create(&orderItem).Error; err != nil {
                return err
            }
        }
        
        // 3. 清空购物车
        if err := tx.Where("user_id = ? AND selected = ?", userID, true).Delete(&models.CartItem{}).Error; err != nil {
            return err
        }
        
        return nil // 返回 nil 提交事务
    })
}
```

### 3.4 软删除 (Soft Delete)

由于使用了 `gorm.Model`，调用 `Delete` 时仅会更新 `deleted_at` 字段，数据仍在数据库中。

```go
// 软删除
db.Delete(&product) 
// UPDATE products SET deleted_at="2024-01-01..." WHERE id = 123;

// 查询时 GORM 会自动过滤已删除记录
db.Find(&products) // WHERE deleted_at IS NULL

// 如果需要查询包含已删除的记录
db.Unscoped().Find(&products)

// 物理删除 (慎用)
db.Unscoped().Delete(&product)
```

---

## 4. SQL 调试与运维

### 4.1 常用 SQL 查询 (CLI/pgAdmin)

在开发过程中，有时直接使用 SQL 更高效。

**查看表大小**:
```sql
SELECT 
    relname AS "relation", 
    pg_size_pretty(pg_total_relation_size(C.oid)) AS "total_size" 
FROM pg_class C 
LEFT JOIN pg_namespace N ON (N.oid = C.relnamespace) 
WHERE nspname = 'public' 
ORDER BY pg_total_relation_size(C.oid) DESC;
```

**检查死锁/活跃查询**:
```sql
SELECT pid, usename, state, query 
FROM pg_stat_activity 
WHERE state != 'idle';
```

**手动验证库存**:
```sql
SELECT id, name, stock FROM product_skus WHERE id = 1;
```

### 4.2 备份与恢复

**备份整个数据库**:
```bash
# 格式: pg_dump -U [user] -h [host] [dbname] > [filename]
pg_dump -U postgres -h localhost go_flutter_mall > backup_$(date +%Y%m%d).sql
```

**恢复数据库**:
```bash
# 1. 确保数据库存在 (如果重建)
createdb -U postgres go_flutter_mall

# 2. 导入数据
psql -U postgres -h localhost -d go_flutter_mall < backup_20251228.sql
```

---

## 5. 性能优化建议

1.  **索引 (Indexing)**:
    *   GORM 的 `AutoMigrate` 会自动创建主键和外键索引。
    *   对于频繁查询的字段（如 `orders.status`, `products.category_id`），应在模型 tag 中添加 `index`。
    *   示例: `Status int `gorm:"index"``

2.  **避免 N+1 查询**:
    *   在循环中查询数据库是性能杀手。
    *   **错误**: 循环 OrderList，在循环内查询 User。
    *   **正确**: 使用 `Preload("User")` 一次性加载。

3.  **连接池配置**:
    *   在 `config/database.go` 中可以进一步配置连接池参数，适应高并发。
    ```go
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)  // 空闲连接数
    sqlDB.SetMaxOpenConns(100) // 最大连接数
    sqlDB.SetConnMaxLifetime(time.Hour)
    ```

## 6. 常见问题排查

*   **Error: "pq: role 'postgres' does not exist"**:
    *   检查 `dsn` 中的用户名是否正确，Mac Homebrew 安装的 Postgres 默认用户可能为当前系统用户名。
*   **Error: "connection refused"**:
    *   确保 Postgres 服务已启动 (`brew services list` 或 `docker ps`)。
*   **时间显示不正确**:
    *   检查 DSN 中的 `TimeZone` 设置。
    *   检查 Go 结构体中的时间字段类型是否为 `time.Time`。

