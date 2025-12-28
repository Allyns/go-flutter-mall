package cart

import (
	"net/http"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
)

// AddToCartInput 定义添加到购物车的请求参数
type AddToCartInput struct {
	ProductID uint `json:"product_id" binding:"required"`     // 商品 ID
	SKUID     uint `json:"sku_id"`                            // SKU ID (可选)
	Quantity  int  `json:"quantity" binding:"required,min=1"` // 数量，至少为 1
}

// GetCart 获取当前用户的购物车列表
// @Summary      Get Cart
// @Description  Get list of items in the user's shopping cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.CartItem
// @Failure      500  {object}  map[string]interface{}
// @Router       /cart [get]
func GetCart(c *gin.Context) {
	// 从上下文中获取用户 ID (由 Auth 中间件设置)
	userID, _ := c.Get("userID")

	var cartItems []models.CartItem
	// 查询该用户的购物车项，并预加载商品信息
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}

// AddToCart 添加商品到购物车
// @Summary      Add to Cart
// @Description  Add a product to the shopping cart or update quantity if exists
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      AddToCartInput  true  "Cart Item Info"
// @Success      200    {object}  models.CartItem
// @Success      201    {object}  models.CartItem
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /cart [post]
func AddToCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input AddToCartInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查该商品是否已在购物车中
	var existingItem models.CartItem
	// 使用 Limit(1).Find 避免 First 抛出 record not found 错误日志
	result := config.DB.Where("user_id = ? AND product_id = ? AND sku_id = ?", userID, input.ProductID, input.SKUID).Limit(1).Find(&existingItem)

	if result.Error != nil {
		// 真实的数据库错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking cart"})
		return
	}

	if result.RowsAffected > 0 {
		// 如果已存在，则更新数量
		existingItem.Quantity += input.Quantity
		if err := config.DB.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
		c.JSON(http.StatusOK, existingItem)
	} else {
		// 如果不存在，则创建新条目
		newItem := models.CartItem{
			UserID:    userID.(uint),
			ProductID: input.ProductID,
			SKUID:     input.SKUID,
			Quantity:  input.Quantity,
			Selected:  true,
		}
		if err := config.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart item"})
			return
		}
		c.JSON(http.StatusCreated, newItem)
	}
}

// UpdateCartItem 更新购物车项 (如数量、选中状态)
// @Summary      Update Cart Item
// @Description  Update quantity or selected status of a cart item
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int              true  "Cart Item ID"
// @Param        input  body      models.CartItem  true  "Updated Info"
// @Success      200    {object}  models.CartItem
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Router       /cart/{id} [put]
func UpdateCartItem(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var item models.CartItem
	// 确保只能更新自己的购物车项
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DeleteCartItem 删除购物车项
// DELETE /api/cart/:id
func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// 执行物理删除或软删除 (取决于 Model 定义)
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}

// GetCartCount 获取购物车商品总数量
// GET /api/cart/count
func GetCartCount(c *gin.Context) {
	userID, _ := c.Get("userID")
	var count int64
	// 统计该用户的购物车条目数
	config.DB.Model(&models.CartItem{}).Where("user_id = ?", userID).Count(&count)

	// 注意: 这里简单返回条目数，也可以改为返回商品总件数 (Sum(Quantity))
	c.JSON(http.StatusOK, gin.H{"count": count})
}
