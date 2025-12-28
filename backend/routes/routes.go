package routes

import (
	"go-flutter-mall/backend/controllers"
	"go-flutter-mall/backend/controllers/admin"
	"go-flutter-mall/backend/controllers/cart"
	"go-flutter-mall/backend/controllers/chat"
	"go-flutter-mall/backend/controllers/notification"
	"go-flutter-mall/backend/controllers/order"
	"go-flutter-mall/backend/controllers/product"
	"go-flutter-mall/backend/controllers/search"
	"go-flutter-mall/backend/middleware"
	"go-flutter-mall/backend/pkg/websocket"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-flutter-mall/backend/docs" // Import generated docs
)

// SetupRoutes 配置所有的 API 路由
// r: Gin 引擎实例
func SetupRoutes(r *gin.Engine, hub *websocket.Hub) {
	// 添加根路由健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Go Flutter Mall API is running",
			"version": "1.0.0",
		})
	})

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建 /api 路由组
	api := r.Group("/api")
	{
		// WebSocket 路由
		api.GET("/ws", func(c *gin.Context) {
			websocket.ServeWs(hub, c)
		})

		// Auth 路由 (公开)
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register) // 用户注册
			auth.POST("/login", controllers.Login)       // 用户登录
			auth.POST("/admin/login", admin.Login)       // 管理员登录

			// Admin Stats (临时公开，实际应加 AdminMiddleware)
			auth.GET("/admin/stats", admin.GetDashboardStats)

			// 需认证的 Auth 路由
			auth.GET("/me", middleware.AuthMiddleware(), controllers.GetUserProfile) // 获取当前用户信息
		}

		// 商品路由 (公开)
		products := api.Group("/products")
		{
			products.GET("", product.GetProducts)                   // 获取商品列表
			products.GET("/:id", product.GetProductDetail)          // 获取商品详情
			products.GET("/:id/reviews", product.GetProductReviews) // 获取商品评价

			// 管理员接口 (需认证)
			// TODO: Add AdminMiddleware
			products.POST("", product.CreateProduct)       // 创建商品
			products.PUT("/:id", product.UpdateProduct)    // 更新商品
			products.DELETE("/:id", product.DeleteProduct) // 删除商品
		}

		// 购物车路由 (需认证)
		// 使用 middleware.AuthMiddleware() 保护该组下的所有路由
		cartGroup := api.Group("/cart", middleware.AuthMiddleware())
		{
			cartGroup.GET("", cart.GetCart)               // 获取购物车
			cartGroup.GET("/count", cart.GetCartCount)    // 获取购物车数量
			cartGroup.POST("", cart.AddToCart)            // 添加到购物车
			cartGroup.PUT("/:id", cart.UpdateCartItem)    // 更新购物车项
			cartGroup.DELETE("/:id", cart.DeleteCartItem) // 删除购物车项
		}

		// 地址路由 (需认证)
		addressGroup := api.Group("/addresses", middleware.AuthMiddleware())
		{
			addressGroup.GET("", controllers.GetAddresses)         // 获取地址列表
			addressGroup.POST("", controllers.CreateAddress)       // 创建地址
			addressGroup.PUT("/:id", controllers.UpdateAddress)    // 更新地址
			addressGroup.DELETE("/:id", controllers.DeleteAddress) // 删除地址
		}

		// 订单路由 (需认证)
		orderGroup := api.Group("/orders", middleware.AuthMiddleware())
		{
			orderGroup.POST("", order.CreateOrder)                     // 创建订单
			orderGroup.GET("", order.GetOrders)                        // 获取订单列表
			orderGroup.GET("/counts", order.GetOrderCounts)            // 获取订单数量统计
			orderGroup.GET("/:id", order.GetOrderDetail)               // 获取订单详情
			orderGroup.POST("/:id/pay", order.PayOrder)                // 支付订单
			orderGroup.PUT("/:id/receipt", order.ConfirmReceipt)       // 确认收货
			orderGroup.POST("/:id/review", order.ReviewOrder)          // 评价订单
			orderGroup.POST("/:id/after-sales", order.ApplyAfterSales) // 申请售后

			// 管理员接口 (临时放在这里，实际应有 AdminMiddleware)
			orderGroup.PUT("/:id/status", order.UpdateOrderStatus) // 更新订单状态
			orderGroup.DELETE("/:id", order.DeleteOrder)           // 删除订单
			orderGroup.GET("/admin/all", order.GetAllOrders)       // 管理员获取所有订单
		}
		// 聊天路由 (需认证)
		chatGroup := api.Group("/chat", middleware.AuthMiddleware())
		{
			chatGroup.GET("/users", chat.GetChatUsers)                   // 获取聊天用户列表
			chatGroup.GET("/messages/:userId", chat.GetMessages)         // 获取聊天记录
			chatGroup.PUT("/read", chat.MarkMessagesAsRead)              // 标记所有管理员消息为已读
			chatGroup.POST("/notification", chat.SendSystemNotification) // 发送系统消息
		}

		// 搜索历史路由 (需认证)
		searchGroup := api.Group("/search", middleware.AuthMiddleware())
		{
			searchGroup.POST("/history", search.AddSearchHistory)     // 添加搜索记录
			searchGroup.GET("/history", search.GetSearchHistory)      // 获取搜索记录
			searchGroup.DELETE("/history", search.ClearSearchHistory) // 清空搜索记录
		}

		// 消息通知路由 (需认证)
		notificationGroup := api.Group("/notifications", middleware.AuthMiddleware())
		{
			notificationGroup.GET("", notification.GetNotifications)            // 获取消息列表
			notificationGroup.PUT("/:id/read", notification.MarkAsRead)         // 标记已读
			notificationGroup.GET("/unread-count", notification.GetUnreadCount) // 获取未读数量

			// 管理员接口
			notificationGroup.GET("/admin/all", notification.GetAllSystemNotifications)           // 获取所有系统通知
			notificationGroup.GET("/admin/user/:userId", notification.GetUserSystemNotifications) // 获取特定用户的通知
		}
	}
}
