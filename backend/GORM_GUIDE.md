# GORM (PostgreSQL) 使用指南

本文档详细说明了 Go-Flutter Mall 项目中 ORM 框架 **GORM** 的使用方法，包括连接配置、模型定义、CRUD 操作、关联查询及事务处理。

## 1. 简介与配置

GORM 是 Go 语言中功能最全的 ORM 库。本项目使用 GORM 连接 PostgreSQL 数据库。

### 1.1 安装依赖
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

### 1.2 初始化连接 (`backend/config/database.go`)

```go
package config

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    // 生产环境应从环境变量或 config.yaml 读取
    dsn := "host=localhost user=postgres password=postgres dbname=go_flutter_mall port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    // 自动迁移模式 (Auto Migrate)
    // 自动创建表、缺失的列和索引 (不会删除未使用的列)
    database.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
    
    DB = database
}
```

---

## 2. 模型定义 (Model Definition)

GORM 倾向于约定优于配置。

### 2.1 基础模型
通常继承 `gorm.Model`，它会自动包含 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` (软删除) 字段。

```go
type User struct {
    gorm.Model
    Username string `gorm:"unique;not null" json:"username"` // 唯一且非空
    Email    string `gorm:"unique;index" json:"email"`       // 唯一且建立索引
    Password string `json:"-"`                               // JSON 序列化时忽略密码
    Age      int    `gorm:"default:18"`                      // 默认值
}
```

### 2.2 PostgreSQL 特有类型
本项目使用了 PostgreSQL 的数组类型 (`text[]`)。

**示例 (`models/product.go`)**:
```go
import "github.com/lib/pq"

type Product struct {
    gorm.Model
    Name   string         `json:"name"`
    Images pq.StringArray `gorm:"type:text[]" json:"images"` // 定义为 text[]
}
```

### 2.3 关联关系
*   **一对多 (Has Many)**: 一个分类下有多个商品。
    ```go
    type Category struct {
        ID       uint
        Products []Product // 一对多
    }
    
    type Product struct {
        ID         uint
        CategoryID uint // 外键
    }
    ```
*   **一对一 (Has One)**: 购物车项关联一个商品。
    ```go
    type CartItem struct {
        ProductID uint
        Product   Product // 预加载使用
    }
    ```

---

## 3. CRUD 操作

以下示例基于 `backend/controllers/order/order_controller.go` 和 `product_controller.go`。

### 3.1 创建 (Create)

```go
// 创建单条记录
user := models.User{Username: "John", Email: "john@example.com"}
result := config.DB.Create(&user) // user.ID 会被自动填充

if result.Error != nil {
    // handle error
}

// 批量创建
var users = []models.User{{Name: "Jinzhu"}, {Name: "Jackson"}}
config.DB.Create(&users)
```

### 3.2 查询 (Read)

```go
// 1. 根据主键获取
var product models.Product
config.DB.First(&product, 10) // SELECT * FROM products WHERE id = 10

// 2. 条件查询 (Where)
var cartItems []models.CartItem
// SELECT * FROM cart_items WHERE user_id = ? AND selected = true
config.DB.Where("user_id = ? AND selected = ?", userID, true).Find(&cartItems)

// 3. 排序与分页
var products []models.Product
config.DB.Order("price desc").Limit(10).Offset(0).Find(&products)

// 4. 预加载关联 (Preload) - 解决 N+1 问题
// 获取购物车项的同时，加载商品信息
config.DB.Preload("Product").Find(&cartItems)

// 嵌套预加载: 获取订单 -> 订单项 -> 商品
config.DB.Preload("OrderItems.Product").Find(&orders)
```

### 3.3 更新 (Update)

```go
// 1. 保存所有字段 (Save)
product.Price = 200
config.DB.Save(&product)

// 2. 更新指定字段 (Update/Updates)
// UPDATE products SET price = 200 WHERE id = 10;
config.DB.Model(&product).Update("Price", 200)

// 更新多个字段 (Struct 形式只更新非零值字段!)
config.DB.Model(&product).Updates(models.Product{Price: 200, Stock: 50})

// 更新多个字段 (Map 形式可更新 0 或 false)
config.DB.Model(&product).Updates(map[string]interface{}{"price": 200, "status": 0})
```

### 3.4 删除 (Delete)

```go
// 软删除 (如果模型包含 DeletedAt)
config.DB.Delete(&product)
// UPDATE products SET deleted_at="2024-..." WHERE id = 10;

// 永久删除 (Unscoped)
config.DB.Unscoped().Delete(&product)

// 批量删除
config.DB.Where("user_id = ?", 10).Delete(&models.CartItem{})
```

---

## 4. 事务处理 (Transactions)

在创建订单等涉及多表操作的场景，必须使用事务。

**示例 (`controllers/order/order_controller.go`)**:

```go
func CreateOrder(c *gin.Context) {
    // 1. 开启事务
    tx := config.DB.Begin()
    
    // 务必在函数结束时处理 panic 等异常回滚
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 2. 在事务中执行操作 (使用 tx 而不是 config.DB)
    
    // 扣减库存 (乐观锁示例: 检查库存是否足够)
    if err := tx.Model(&sku).Where("stock >= ?", quantity).
        UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
        tx.Rollback()
        return
    }

    // 创建订单
    if err := tx.Create(&order).Error; err != nil {
        tx.Rollback()
        return
    }

    // 3. 提交事务
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return
    }
}
```

---

## 5. 高级特性

### 5.1 Scopes (复用查询逻辑)
```go
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
    return db.Where("amount > ?", 1000)
}

func PaidOrders(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "paid")
}

// 使用
db.Scopes(AmountGreaterThan1000, PaidOrders).Find(&orders)
```

### 5.2 Hooks (钩子函数)
在创建/更新前后自动执行逻辑。

```go
// 在创建 User 前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.Password = Hash(u.Password)
    return
}
```

### 5.3 SQL 表达式
```go
// stock = stock - 1
db.Model(&product).Update("stock", gorm.Expr("stock - ?", 1))
```

---

## 6. 常见问题排查

1.  **零值不更新**: 使用 `Updates(struct)` 时，`0`、`""`、`false` 会被忽略。解决方法：使用 `map[string]interface{}` 或 `Select/Omit` 指定字段。
2.  **关联保存**: GORM 默认会在创建/更新时自动保存关联。如果不希望这样，使用 `Set("gorm:association_autoupdate", false)` 或 `Omit("Association")`。
3.  **连接池设置**:
    ```go
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)  // 空闲连接数
    sqlDB.SetMaxOpenConns(100) // 最大连接数
    sqlDB.SetConnMaxLifetime(time.Hour)
    ```
