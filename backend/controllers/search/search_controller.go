package search

import (
	"context"
	"net/http"
	"time"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddSearchHistoryInput struct {
	Keyword string `json:"keyword"`
}

// AddSearchHistory 添加搜索记录
// @Summary      Add Search History
// @Description  Add a keyword to the user's search history
// @Tags         Search
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      AddSearchHistoryInput  true  "Search Keyword"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /search/history [post]
func AddSearchHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input struct {
		Keyword string `json:"keyword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := config.MongoDB.Collection("search_history")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查是否已存在相同的关键词，如果存在则更新时间
	filter := bson.M{"user_id": userID, "keyword": input.Keyword}
	update := bson.M{
		"$set": bson.M{
			"created_at": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save search history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Search history saved"})
}

// GetSearchHistory 获取搜索历史
// @Summary      Get Search History
// @Description  Get the user's recent search history
// @Tags         Search
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.SearchHistory
// @Failure      500  {object}  map[string]interface{}
// @Router       /search/history [get]
func GetSearchHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	collection := config.MongoDB.Collection("search_history")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 按时间倒序查询前 10 条
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(10)
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		// Log the actual error for debugging
		println("Error fetching search history:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch search history"})
		return
	}
	defer cursor.Close(ctx)

	var history []models.SearchHistory
	if err = cursor.All(ctx, &history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode search history"})
		return
	}

	// 如果为空，返回空切片而不是 null
	if history == nil {
		history = []models.SearchHistory{}
	}

	c.JSON(http.StatusOK, history)
}

// ClearSearchHistory 清空搜索历史
// @Summary      Clear Search History
// @Description  Clear all search history for the user
// @Tags         Search
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /search/history [delete]
func ClearSearchHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	collection := config.MongoDB.Collection("search_history")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear search history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Search history cleared"})
}
