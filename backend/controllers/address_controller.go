package controllers

import (
	"net/http"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
)

// CreateAddressInput 创建/更新地址的输入参数
type CreateAddressInput struct {
	ReceiverName  string `json:"receiver_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district" binding:"required"`
	DetailAddress string `json:"detail_address" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

// GetAddresses 获取用户的收货地址列表
// @Summary      Get Addresses
// @Description  Get a list of shipping addresses for the authenticated user
// @Tags         Address
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Address
// @Failure      500  {object}  map[string]interface{}
// @Router       /addresses [get]
func GetAddresses(c *gin.Context) {
	userID, _ := c.Get("userID")
	var addresses []models.Address

	if err := config.DB.Where("user_id = ?", userID).Order("is_default desc, created_at desc").Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// CreateAddress 创建新地址
// @Summary      Create Address
// @Description  Create a new shipping address
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      CreateAddressInput  true  "Address Info"
// @Success      201    {object}  models.Address
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /addresses [post]
func CreateAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input CreateAddressInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设为默认地址，先将其他地址设为非默认
	if input.IsDefault {
		config.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	address := models.Address{
		UserID:        userID.(uint),
		ReceiverName:  input.ReceiverName,
		Phone:         input.Phone,
		Province:      input.Province,
		City:          input.City,
		District:      input.District,
		DetailAddress: input.DetailAddress,
		IsDefault:     input.IsDefault,
	}

	// 如果是第一条地址，强制设为默认
	var count int64
	config.DB.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count)
	if count == 0 {
		address.IsDefault = true
	}

	if err := config.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// UpdateAddress 更新地址
// @Summary      Update Address
// @Description  Update an existing shipping address
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                 true  "Address ID"
// @Param        input  body      CreateAddressInput  true  "Address Info"
// @Success      200    {object}  models.Address
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /addresses/{id} [put]
func UpdateAddress(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	var input CreateAddressInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var address models.Address
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if input.IsDefault && !address.IsDefault {
		config.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	address.ReceiverName = input.ReceiverName
	address.Phone = input.Phone
	address.Province = input.Province
	address.City = input.City
	address.District = input.District
	address.DetailAddress = input.DetailAddress
	address.IsDefault = input.IsDefault

	if err := config.DB.Save(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, address)
}

// DeleteAddress 删除地址
// @Summary      Delete Address
// @Description  Delete a shipping address
// @Tags         Address
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /addresses/{id} [delete]
func DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted"})
}
