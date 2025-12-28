package admin

import (
	"net/http"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalUsers    int64   `json:"total_users"`
	TotalOrders   int64   `json:"total_orders"`
	TotalSales    float64 `json:"total_sales"`
	TotalProducts int64   `json:"total_products"`
}

// GetDashboardStats 获取仪表盘统计数据
// @Summary      Get Dashboard Stats
// @Description  Get statistics for the admin dashboard
// @Tags         Admin
// @Produce      json
// @Success      200  {object}  DashboardStats
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/admin/stats [get]
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 1. 统计用户总数
	if err := config.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	// 2. 统计订单总数 (排除已取消的订单)
	if err := config.DB.Model(&models.Order{}).Where("status != ?", -1).Count(&stats.TotalOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count orders"})
		return
	}

	// 3. 统计总销售额 (排除已取消和待支付的订单)
	// 假设已支付、待发货、待收货、已完成都算销售额
	var result struct {
		Total float64
	}
	if err := config.DB.Model(&models.Order{}).
		Where("status IN ?", []int{1, 2, 3, 4, 5}).
		Select("sum(total_amount) as total").
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate sales"})
		return
	}
	stats.TotalSales = result.Total

	// 4. 统计商品总数
	if err := config.DB.Model(&models.Product{}).Count(&stats.TotalProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
