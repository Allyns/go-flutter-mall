package notification

import (
	"net/http"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
)

// GetNotifications 获取用户的消息列表
// @Summary      Get Notifications
// @Description  Get a list of notifications for the authenticated user
// @Tags         Notification
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Notification
// @Failure      500  {object}  map[string]interface{}
// @Router       /notifications [get]
func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("userID")
	var notifications []models.Notification

	// 按时间倒序查询
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead 标记消息为已读
// @Summary      Mark Notification as Read
// @Description  Mark a specific notification as read
// @Tags         Notification
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Notification ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /notifications/{id}/read [put]
func MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	if err := config.DB.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

// GetUnreadCount 获取未读消息数量
// @Summary      Get Unread Count
// @Description  Get the count of unread notifications and chat messages
// @Tags         Notification
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]int64
// @Failure      500  {object}  map[string]interface{}
// @Router       /notifications/unread-count [get]
func GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userID")
	var count int64

	// 统计未读通知
	if err := config.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count unread notifications"})
		return
	}

	// 统计未读聊天消息 (SenderType='admin', ReceiverID=userID, IsRead=false)
	var chatCount int64
	if err := config.DB.Model(&models.ChatMessage{}).Where("receiver_id = ? AND sender_type = ? AND is_read = ?", userID, "admin", false).Count(&chatCount).Error; err != nil {
		// 如果聊天消息表没有 is_read 字段或查询失败，暂不计入总数或打日志
		// 但根据 models/admin.go，ChatMessage 确实有 IsRead 字段
	}

	total := count + chatCount
	c.JSON(http.StatusOK, gin.H{"unread_count": total, "notification_count": count, "chat_count": chatCount})
}

// GetAllSystemNotifications 管理员获取所有系统通知
// @Summary      Get All Notifications
// @Description  Get a list of all system notifications (Admin only)
// @Tags         Notification
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Notification
// @Failure      500  {object}  map[string]interface{}
// @Router       /notifications/admin/all [get]
func GetAllSystemNotifications(c *gin.Context) {
	var notifications []models.Notification

	// 预加载 User 信息，以便显示接收者
	if err := config.DB.Preload("User").Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUserSystemNotifications 管理员获取特定用户的通知
// @Summary      Get User Notifications
// @Description  Get notifications for a specific user (Admin only)
// @Tags         Notification
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path      int  true  "User ID"
// @Success      200     {array}   models.Notification
// @Failure      500     {object}  map[string]interface{}
// @Router       /notifications/admin/user/{userId} [get]
func GetUserSystemNotifications(c *gin.Context) {
	userID := c.Param("userId")
	var notifications []models.Notification

	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
