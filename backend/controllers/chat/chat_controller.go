package chat

import (
	"net/http"
	"strconv"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
)

// GetChatUsers 获取最近联系的用户列表
// @Summary      Get Chat Users
// @Description  Get a list of users who have chatted with the admin
// @Tags         Chat
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]interface{}
// @Router       /chat/users [get]
func GetChatUsers(c *gin.Context) {
	users := []models.User{} // Initialize as empty slice

	// 查找所有发过消息的用户 (sender_type='user')
	var senderIDs []uint
	config.DB.Model(&models.ChatMessage{}).Where("sender_type = ?", "user").Distinct("sender_id").Pluck("sender_id", &senderIDs)

	// 查找所有收到过消息的用户 (receiver_id where sender_type='admin')
	var receiverIDs []uint
	config.DB.Model(&models.ChatMessage{}).Where("sender_type = ?", "admin").Distinct("receiver_id").Pluck("receiver_id", &receiverIDs)

	// 合并并去重
	idMap := make(map[uint]bool)
	for _, id := range senderIDs {
		idMap[id] = true
	}
	for _, id := range receiverIDs {
		idMap[id] = true
	}

	var allUserIDs []uint
	for id := range idMap {
		allUserIDs = append(allUserIDs, id)
	}

	if len(allUserIDs) > 0 {
		if err := config.DB.Where("id IN ?", allUserIDs).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
			return
		}
	} else {
		// 如果没有任何消息，可以返回所有用户，方便管理员主动发起会话
		// 或者返回空，取决于业务需求。这里返回所有用户
		config.DB.Find(&users)
	}

	c.JSON(http.StatusOK, users)
}

// GetMessages 获取与指定用户的聊天记录
// @Summary      Get Messages
// @Description  Get chat history with a specific user
// @Tags         Chat
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path      int  true  "User ID"
// @Success      200     {array}   models.ChatMessage
// @Failure      500     {object}  map[string]interface{}
// @Router       /chat/messages/{userId} [get]
func GetMessages(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, _ := strconv.Atoi(userIDStr)

	messages := []models.ChatMessage{} // Initialize as empty slice

	// 查询 (SenderID = userID AND SenderType = 'user') OR (ReceiverID = userID AND SenderType = 'admin')
	// 假设 Admin 只有一种类型，且所有 Admin 都能看到所有消息

	if err := config.DB.Where(
		"(sender_id = ? AND sender_type = ?) OR (receiver_id = ? AND sender_type = ?)",
		userID, "user", userID, "admin",
	).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// MarkMessagesAsRead 标记与管理员的聊天消息为已读
// @Summary      Mark Messages as Read
// @Description  Mark all messages from admin as read for the current user
// @Tags         Chat
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /chat/read [put]
func MarkMessagesAsRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 更新所有发给该用户的未读消息 (SenderType='admin', ReceiverID=userID)
	if err := config.DB.Model(&models.ChatMessage{}).
		Where("receiver_id = ? AND sender_type = ? AND is_read = ?", userID, "admin", false).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark messages as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}

type SendNotificationInput struct {
	UserID  uint   `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// SendSystemNotification 发送系统消息
// @Summary      Send System Notification
// @Description  Send a system notification to a user
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      SendNotificationInput  true  "Notification Info"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /chat/notification [post]
func SendSystemNotification(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := models.Notification{
		UserID:  input.UserID,
		Title:   input.Title,
		Content: input.Content,
		IsRead:  false,
	}

	if err := config.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}
