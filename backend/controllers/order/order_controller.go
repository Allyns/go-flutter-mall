package order

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"
	"go-flutter-mall/backend/pkg/kafka"
	"go-flutter-mall/backend/pkg/scheduler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateOrderInput 创建订单的输入参数
type CreateOrderInput struct {
	AddressID uint `json:"address_id" binding:"required"` // 收货地址 ID
}

// CreateOrder 创建新订单
// @Summary      Create Order
// @Description  Create a new order from selected cart items
// @Tags         Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      CreateOrderInput  true  "Order Info"
// @Success      201    {object}  models.Order
// @Failure      400    {object}  map[string]interface{}
// @Failure      409    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /orders [post]
func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input CreateOrderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开启数据库事务
	tx := config.DB.Begin()

	// 1. 获取购物车中选中的商品
	var cartItems []models.CartItem
	if err := tx.Preload("Product").Where("user_id = ? AND selected = ?", userID, true).Find(&cartItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	if len(cartItems) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items selected in cart"})
		return
	}

	// Redis 并发控制 (分布式锁)
	ctx := context.Background()
	for _, item := range cartItems {
		// 锁键: lock:product:{id}
		lockKey := fmt.Sprintf("lock:product:%d", item.ProductID)
		locked, err := config.RedisClient.SetNX(ctx, lockKey, 1, 5*time.Second).Result()
		if err != nil {
			// Redis 连接失败，降级处理：跳过锁检查，直接依赖 DB 的事务和行锁 (InnoDB)
			// 或者记录日志并继续。这里选择记录日志并继续，假设 DB 能处理并发 (乐观锁或行锁)
			fmt.Printf("Redis SetNX failed (connection refused?): %v. Proceeding without distributed lock.\n", err)
			// 不要 return，继续执行
		} else if !locked {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Server busy, please try again"}) // 并发冲突
			return
		} else {
			// 只有获取锁成功才需要在结束后释放
			defer config.RedisClient.Del(ctx, lockKey)
		}
	}

	// 2. 计算总金额并准备订单项
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range cartItems {
		price := item.Product.Price
		// TODO: 如果有 SKU，应该取 SKU 的价格
		// if item.SKUID != 0 { ... }

		totalAmount += price * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID:    item.ProductID,
			ProductName:  item.Product.Name,
			ProductImage: item.Product.CoverImage,
			SKUID:        item.SKUID,
			Price:        price,
			Quantity:     item.Quantity,
		})

		// 3. 扣减库存
		// 关键修复：使用 WHERE 条件检查库存是否充足 (stock >= quantity)
		result := tx.Model(&models.Product{}).
			Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
			UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))

		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
			return
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Insufficient stock for product: %s", item.Product.Name)})
			return
		}
	}

	// 4. 创建订单记录
	order := models.Order{
		OrderNo:     fmt.Sprintf("%d%d", time.Now().UnixNano(), userID.(uint)), // 生成唯一订单号
		UserID:      userID.(uint),
		TotalAmount: totalAmount,
		Status:      0, // 待支付
		AddressID:   input.AddressID,
		Items:       orderItems,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Create order failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create order: %v", err)})
		return
	}

	// 5. 清空购物车中已购买的商品
	if err := tx.Where("user_id = ? AND selected = ?", userID, true).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	// 提交事务
	tx.Commit()

	// 6. 发送消息通知 (Kafka)
	// 生产 "OrderCreated" 消息
	kafka.SendOrderEvent(kafka.OrderEvent{
		OrderID:   order.ID,
		UserID:    order.UserID,
		EventType: "created",
	})

	// 7. 添加到延时队列 (Redis ZSet)
	// 30分钟后触发超时
	// 如果 Redis 不可用，这里可能会失败，记录错误但不影响主流程
	if err := scheduler.AddToDelayQueue(order.ID, order.UserID, 30*time.Minute); err != nil {
		fmt.Printf("Warning: Failed to add to delay queue (Redis down?): %v\n", err)
	}

	// 8. 创建本地通知 (DB) - 也可以移到 Kafka Consumer 中处理
	notification := models.Notification{
		UserID:  userID.(uint),
		Title:   "订单创建成功",
		Content: fmt.Sprintf("您的订单 %s 已成功创建，请尽快支付。", order.OrderNo),
		IsRead:  false,
	}
	config.DB.Create(&notification)

	c.JSON(http.StatusCreated, order)
}

// GetOrders 获取订单列表
// @Summary      Get Order List
// @Description  Get a list of orders for the authenticated user
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        status  query     int  false  "Order Status"
// @Success      200     {array}   models.Order
// @Failure      500     {object}  map[string]interface{}
// @Router       /orders [get]
func GetOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	var orders []models.Order
	status := c.Query("status")

	query := config.DB.Preload("Items").Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询用户的订单，按创建时间倒序排列，并预加载订单项
	if err := query.Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderCounts 获取各状态订单数量
// @Summary      Get Order Counts
// @Description  Get the count of orders for each status
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[int]int
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/counts [get]
func GetOrderCounts(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 定义结果结构
	type StatusCount struct {
		Status int `json:"status"`
		Count  int `json:"count"`
	}

	var results []StatusCount

	// 执行聚合查询
	// SELECT status, count(*) as count FROM orders WHERE user_id = ? GROUP BY status
	if err := config.DB.Model(&models.Order{}).Select("status, count(*) as count").Where("user_id = ?", userID).Group("status").Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order counts"})
		return
	}

	// 将结果转换为 map 方便前端使用
	// counts: { "0": 2, "1": 5, ... }
	counts := make(map[int]int)
	for _, r := range results {
		counts[r.Status] = r.Count
	}

	c.JSON(http.StatusOK, counts)
}

// GetOrderDetail 获取订单详情
// @Summary      Get Order Detail
// @Description  Get detailed information of a specific order
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  models.Order
// @Failure      404  {object}  map[string]interface{}
// @Router       /orders/{id} [get]
func GetOrderDetail(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	var order models.Order

	// 查询特定订单，确保只能查看自己的订单
	if err := config.DB.Preload("Items").Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// PayOrder 模拟支付订单
// @Summary      Pay Order
// @Description  Simulate payment for an order
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id}/pay [post]
func PayOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// 更新订单状态为已支付 (1)
	// 必须确保订单状态为 0 (待支付) 才能支付
	result := config.DB.Model(&models.Order{}).Where("id = ? AND user_id = ? AND status = ?", id, userID, 0).Update("status", 1)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pay order"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be paid (maybe cancelled or already paid)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order paid successfully"})
}

// ConfirmReceipt 确认收货
// @Summary      Confirm Receipt
// @Description  Confirm receipt of an order
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id}/receipt [put]
func ConfirmReceipt(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// 只有状态为 2 (待收货) 的订单才能确认收货，更新为 3 (待评价)
	if err := config.DB.Model(&models.Order{}).Where("id = ? AND user_id = ? AND status = ?", id, userID, 2).Update("status", 3).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Receipt confirmed successfully"})
}

// ReviewOrderInput 评价输入
type ReviewOrderInput struct {
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

// ReviewOrder 评价订单
// @Summary      Review Order
// @Description  Submit a review for an order
// @Tags         Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int               true  "Order ID"
// @Param        input  body      ReviewOrderInput  true  "Review Content"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /orders/{id}/review [post]
func ReviewOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	var input ReviewOrderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开启事务
	tx := config.DB.Begin()

	// 1. 创建评价记录
	// 遍历订单项，为每个商品创建评价 (这里简化为只评价第一个商品，或者需要前端传 ProductID)
	// 实际上 ReviewOrderInput 应该包含 ProductID 或者我们只针对订单评价
	// 假设 ReviewOrder 是针对整个订单的，但 Review 模型关联了 ProductID
	// 这里我们查询订单包含的商品，将评价关联到第一个商品，或者改进 API 让用户针对每个商品评价
	// 简单起见，我们关联到订单的第一个商品
	var order models.Order
	if err := tx.Preload("Items").First(&order, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if len(order.Items) > 0 {
		review := models.Review{
			UserID:    userID.(uint),
			ProductID: order.Items[0].ProductID, // 默认关联第一个商品
			OrderID:   order.ID,
			Content:   input.Content,
			Rating:    input.Rating,
			Status:    1,
		}
		if err := tx.Create(&review).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
			return
		}
	}

	// 2. 更新订单状态
	if err := tx.Model(&models.Order{}).Where("id = ? AND user_id = ? AND status = ?", id, userID, 3).Update("status", 4).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Order reviewed successfully"})
}

// ApplyAfterSales 申请售后
// @Summary      Apply After-Sales
// @Description  Apply for after-sales service for an order
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id}/after-sales [post]
func ApplyAfterSales(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// 只有状态为 4 (已完成) 的订单才能申请售后，申请后更新为 5 (售后中)
	if err := config.DB.Model(&models.Order{}).Where("id = ? AND user_id = ? AND status = ?", id, userID, 4).Update("status", 5).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply for after-sales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "After-sales applied successfully"})
}

// UpdateOrderStatusInput 更新订单状态的输入参数
type UpdateOrderStatusInput struct {
	Status int `json:"status" binding:"required"`
}

// UpdateOrderStatus 管理员更新订单状态
// @Summary      Update Order Status
// @Description  Update the status of an order (Admin only)
// @Tags         Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                     true  "Order ID"
// @Param        input  body      UpdateOrderStatusInput  true  "New Status"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var input UpdateOrderStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查订单是否存在
	var order models.Order
	if err := config.DB.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 更新状态
	if err := config.DB.Model(&order).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "order": order})
}

// GetAllOrders 管理员获取所有订单
// @Summary      Get All Orders
// @Description  Get a list of all orders (Admin only)
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Order
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/admin/all [get]
func GetAllOrders(c *gin.Context) {
	var orders []models.Order

	// 查询所有订单，按创建时间倒序排列
	if err := config.DB.Preload("Items").Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// DeleteOrder 管理员删除订单 (软删除)
// @Summary      Delete Order
// @Description  Delete an order (Admin only)
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	// 检查订单是否存在
	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 软删除
	if err := config.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
