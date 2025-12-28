# Go 语言基础与实战教程

本文档旨在为 `go-flutter-mall` 项目的开发者提供一份全面而详细的 Go 语言指南，涵盖从基础语法到 Web 开发实战的核心知识点，并结合本项目代码进行讲解。

## 1. Go 语言简介

Go (Golang) 是 Google 开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言。

**核心特点**：
- **语法简洁**：去除了 C++ 中复杂的继承、泛型（Go 1.18 后引入泛型）、指针运算等。
- **原生并发**：Goroutine 和 Channel 让并发编程变得极易上手。
- **高性能**：编译速度快，运行效率接近 C/C++。
- **标准库强大**：`net/http` 库让构建 Web 服务变得极其简单。

---

## 2. 基础语法 (Basics)

### 2.1 变量与常量

Go 推荐使用驼峰命名法。

```go
package main

import "fmt"

// 常量定义
const AppName = "GoMall"

func main() {
    // 1. var 声明 (会自动初始化为零值，int 为 0，string 为 "")
    var x int
    x = 10

    // 2. 类型推导
    var y = 20

    // 3. 简短声明 (仅限函数内部)
    z := 30 
    
    // 多变量声明
    name, age := "Alice", 25

    fmt.Printf("Name: %s, Age: %d\n", name, age)
}
```

### 2.2 核心数据结构：切片 (Slice) 与 映射 (Map)

在 Go 中，数组是固定长度的，**切片 (Slice)** 才是动态数组，使用最为频繁。

```go
// --- Slice (切片) ---
// 创建一个初始长度为 0，容量为 5 的字符串切片
tags := make([]string, 0, 5)

// 添加元素
tags = append(tags, "Electronics", "Books")

// 切片操作 [start:end] (左闭右开)
subTags := tags[0:1] // ["Electronics"]

// --- Map (映射/字典) ---
// 创建 Map: key 为 string, value 为 int
cart := make(map[string]int)
cart["iPhone 15"] = 2
cart["MacBook"] = 1

// 检查 Key 是否存在
count, exists := cart["iPad"]
if !exists {
    fmt.Println("iPad not in cart")
}
```

### 2.3 流程控制

Go 只有 `for` 循环，没有 `while`。

```go
// 1. 标准 For 循环
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// 2. 类似 While 的循环
n := 0
for n < 5 {
    n++
}

// 3. 遍历 Slice 或 Map (Range)
items := []string{"A", "B", "C"}
for index, value := range items {
    fmt.Printf("%d: %s\n", index, value)
}

// 4. Switch (默认不需要 break)
status := "pending"
switch status {
case "pending":
    fmt.Println("等待支付")
case "paid":
    fmt.Println("已支付")
default:
    fmt.Println("未知状态")
}
```

---

## 3. 函数与面向对象 (Functions & OOP)

Go 没有 `class`，而是通过 **Struct (结构体)** 和 **Method (方法)** 来实现面向对象。

### 3.1 函数定义

Go 函数支持多返回值，这在错误处理中非常常用。

```go
// 参数类型在变量名之后，返回类型在参数列表之后
func div(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// 调用
res, err := div(10, 2)
if err != nil {
    // 处理错误
}
```

### 3.2 结构体 (Struct) 与 方法 (Method)

这是 Go 模拟“类”的方式。

```go
// 定义 User 结构体
type User struct {
    ID       uint
    Username string
    Email    string
}

// 定义 User 的方法
// (u *User) 称为接收者 (Receiver)，类似 Python 的 self 或 Java 的 this
func (u *User) Notify(message string) {
    fmt.Printf("Sending email to %s: %s\n", u.Email, message)
}

// 使用
func main() {
    user := User{ID: 1, Username: "allyn", Email: "test@example.com"}
    user.Notify("Welcome!")
}
```

### 3.3 接口 (Interface)

接口定义了一组行为规范。只要某个类型实现了接口中的所有方法，它就实现了该接口（隐式实现）。

```go
// 定义接口
type Notifier interface {
    Notify(message string)
}

// 多态调用
func SendAlert(n Notifier, msg string) {
    n.Notify(msg)
}

// User 实现了 Notify 方法，所以 User 也是 Notifier
// SendAlert(&user, "Alert!")
```

---

## 4. 并发编程 (Concurrency)

Go 的杀手锏特性。

### 4.1 Goroutine

轻量级线程。只需在函数调用前加 `go` 关键字。

```go
func processTask(id int) {
    fmt.Printf("Processing task %d\n", id)
}

func main() {
    for i := 0; i < 10; i++ {
        go processTask(i) // 启动 10 个并发任务
    }
    // 注意：主线程退出，所有 Goroutine 也会立即终止
    // 实际开发中需配合 WaitGroup 使用
}
```

### 4.2 Channel

用于 Goroutine 之间的通信，“不要通过共享内存来通信，而要通过通信来共享内存”。

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("Worker", id, "processing job", j)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // 启动 3 个 Worker
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // 发送 5 个任务
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs) // 关闭通道，通知 Worker 没有更多任务了

    // 获取结果
    for a := 1; a <= 5; a++ {
        <-results
    }
}
```

---

## 5. 本项目实战解析 (Go-Flutter Mall)

结合 `backend` 目录下的代码，讲解实际开发中的模式。

### 5.1 Web 框架：Gin

本项目使用 Gin (`github.com/gin-gonic/gin`) 作为 Web 框架。

**文件参考**: `backend/routes/routes.go`

```go
// 典型路由注册示例
func SetupRouter() *gin.Engine {
    r := gin.Default()

    // 中间件使用
    r.Use(middleware.CorsMiddleware())

    api := r.Group("/api")
    {
        // 绑定 GET 请求到控制器方法
        api.GET("/products", controllers.GetProducts)
        
        // 路径参数绑定
        api.GET("/products/:id", controllers.GetProductDetail)
        
        // 需要鉴权的路由组
        auth := api.Group("/")
        auth.Use(middleware.AuthMiddleware())
        {
            auth.POST("/cart", controllers.AddToCart)
        }
    }
    return r
}
```

### 5.2 数据库 ORM：GORM

本项目使用 GORM (`gorm.io/gorm`) 操作 PostgreSQL。

**模型定义**: `backend/models/product.go`
```go
type Product struct {
    gorm.Model        // 自动包含 ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string      `json:"name"`
    Price float64     `json:"price"`
    SKUs  []ProductSKU `gorm:"foreignKey:ProductID"` // 一对多关联
}
```

**控制器逻辑**: `backend/controllers/product/product_controller.go`
```go
func GetProducts(c *gin.Context) {
    var products []models.Product
    
    // 获取查询参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize := 10
    
    // GORM 链式调用
    // Offset/Limit 实现分页
    // Preload 实现关联预加载
    result := config.DB.Preload("SKUs").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&products)
        
    if result.Error != nil {
        c.JSON(500, gin.H{"error": "Database error"})
        return
    }
    
    c.JSON(200, gin.H{"data": products})
}
```

### 5.3 错误处理最佳实践

在 Go 中，`if err != nil` 是标准范式。在本项目中，我们通常在 Controller 层统一处理错误并返回 JSON。

```go
// 示例：创建订单
func CreateOrder(c *gin.Context) {
    var req OrderRequest
    // 1. 绑定并校验 JSON 参数
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid params"})
        return
    }

    // 2. 调用业务逻辑 (可能会返回 error)
    err := services.CreateOrderService(req)
    if err != nil {
        // 根据错误类型返回不同状态码
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Success"})
}
```

### 5.4 常用命令

在 `backend` 目录下：

```bash
# 1. 初始化模块 (如果是一个新项目)
go mod init go-flutter-mall/backend

# 2. 整理依赖 (自动下载 import 的包，移除未使用的包)
go mod tidy

# 3. 运行项目
go run main.go

# 4. 编译成二进制文件
go build -o server main.go
```

## 6. 学习资源推荐

1.  **A Tour of Go**: 官方交互式教程，必看 (https://tour.go-zh.org/)
2.  **Go by Example**: 通过例子学 Go (https://gobyexample-cn.github.io/)
3.  **Gin 官方文档**: (https://gin-gonic.com/zh-cn/docs/)
4.  **GORM 官方文档**: (https://gorm.io/zh_CN/docs/)
