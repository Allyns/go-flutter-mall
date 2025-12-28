package product

import (
	"net/http"
	"strconv"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
)

// GetProducts 获取商品列表
// @Summary      Get Product List
// @Description  Get a list of products with pagination and search
// @Tags         Product
// @Produce      json
// @Param        page       query     int     false  "Page number" default(1)
// @Param        page_size  query     int     false  "Page size" default(10)
// @Param        search     query     string  false  "Search keyword"
// @Success      200        {array}   models.Product
// @Failure      500        {object}  map[string]interface{}
// @Router       /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 查询数据库，预加载 SKU 信息
	query := config.DB.Preload("SKUs")

	// 搜索功能
	search := c.Query("search")
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Limit 和 Offset 用于实现分页
	result := query.Limit(pageSize).Offset(offset).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductDetail 获取商品详情
// @Summary      Get Product Detail
// @Description  Get detailed information of a product by ID
// @Tags         Product
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func GetProductDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// 根据 ID 查询商品，并预加载 SKU 信息
	if err := config.DB.Preload("SKUs").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ProductSKUInput 商品 SKU 输入
type ProductSKUInput struct {
	Name  string  `json:"name" binding:"required"`
	Specs string  `json:"specs" binding:"required"` // JSON string
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}

// CreateProductInput 创建商品输入
type CreateProductInput struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Price       float64           `json:"price" binding:"required"`
	Stock       int               `json:"stock" binding:"required"`
	CoverImage  string            `json:"cover_image"`
	CategoryID  uint              `json:"category_id" binding:"required"`
	SKUs        []ProductSKUInput `json:"skus"` // 商品 SKU 列表
}

// CreateProduct 创建商品
// @Summary      Create Product
// @Description  Create a new product (Admin only)
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        input  body      CreateProductInput  true  "Product Info"
// @Success      201    {object}  models.Product
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /products [post]
func CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		CoverImage:  input.CoverImage,
		CategoryID:  input.CategoryID,
		Status:      1, // 默认上架
	}

	// 开启事务
	tx := config.DB.Begin()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// 创建 SKU
	if len(input.SKUs) > 0 {
		for _, skuInput := range input.SKUs {
			sku := models.ProductSKU{
				ProductID: product.ID,
				Name:      skuInput.Name,
				Specs:     skuInput.Specs,
				Price:     skuInput.Price,
				Stock:     skuInput.Stock,
				// Image: skuInput.Image, // 暂时没有图片输入
			}
			if err := tx.Create(&sku).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product SKU"})
				return
			}
		}
	}

	tx.Commit()

	// 重新查询以包含 SKUs
	config.DB.Preload("SKUs").First(&product, product.ID)

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct 更新商品
// @Summary      Update Product
// @Description  Update an existing product (Admin only)
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id     path      int                 true  "Product ID"
// @Param        input  body      CreateProductInput  true  "Product Info"
// @Success      200    {object}  models.Product
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.CoverImage = input.CoverImage
	product.CategoryID = input.CategoryID

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductReviews 获取商品评价
// @Summary      Get Product Reviews
// @Description  Get a list of reviews for a specific product
// @Tags         Product
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {array}   models.Review
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id}/reviews [get]
func GetProductReviews(c *gin.Context) {
	id := c.Param("id")
	var reviews []models.Review

	// 查询评价，关联用户信息
	// 临时方案：移除 Status=1 限制，允许管理员查看所有评价
	// 或者如果想保留前台逻辑，可以在这里判断是否是 Admin 路由
	// 但最简单的修复是：先展示所有
	query := config.DB.Preload("User").Where("product_id = ?", id).Order("created_at desc")

	// 如果是 Admin，可能需要查看所有。如果是前台，只看 Status=1
	// 这里通过 query 参数来区分，或者为了演示，全部返回
	if c.Query("status") == "active" {
		query = query.Where("status = 1")
	}

	if err := query.Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// DeleteProduct 删除商品
// @Summary      Delete Product
// @Description  Delete a product by ID (Admin only)
// @Tags         Product
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
