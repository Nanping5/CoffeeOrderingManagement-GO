package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/middleware"
	"coffee-ordering-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserPoints 获取用户积分信息
func GetUserPoints(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权",
		})
		return
	}

	db := database.GetDB()
	var userPoints models.UserPoints

	if err := db.Where("user_id = ?", userID).First(&userPoints).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "积分账户不存在",
		})
		return
	}

	// 获取会员等级配置
	var levelConfig models.MemberLevelConfig
	db.Where("level_name = ? AND is_active = ?", userPoints.MemberLevel, true).First(&levelConfig)

	// 计算下一等级信息
	var nextLevel *models.NextLevelInfo
	var nextLevelConfig models.MemberLevelConfig
	if err := db.Where("min_points > ? AND is_active = ?", userPoints.LifetimePoints, true).
		Order("min_points ASC").First(&nextLevelConfig).Error; err == nil {
		nextLevel = &models.NextLevelInfo{
			Name:           string(nextLevelConfig.LevelName),
			RequiredPoints: nextLevelConfig.MinPoints,
			PointsNeeded:   nextLevelConfig.MinPoints - userPoints.LifetimePoints,
		}
	}

	// 构建响应
	response := models.UserPointsResponse{
		TotalPoints:     userPoints.TotalPoints,
		AvailablePoints: userPoints.AvailablePoints,
		FrozenPoints:    userPoints.FrozenPoints,
		LifetimePoints:  userPoints.LifetimePoints,
		MemberLevel:     string(userPoints.MemberLevel),
		NextLevel:       nextLevel,
		LevelBenefits: &models.MemberLevelBenefits{
			PointsEarningRate: levelConfig.PointsEarningRate,
			MaxDiscountRate:   levelConfig.MaxDiscountPercentage,
			Benefits:          []string{"积分累积", "生日礼品"}, // 简化版本
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetPointTransactions 获取积分历史记录
func GetPointTransactions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权",
		})
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	transactionType := c.Query("type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := database.GetDB()
	var transactions []models.PointTransaction
	var total int64

	query := db.Model(&models.PointTransaction{}).Where("user_id = ? AND is_active = ?", userID, true)

	// 类型筛选
	if transactionType != "" {
		query = query.Where("transaction_type = ?", transactionType)
	}

	// 日期范围筛选
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			t = t.Add(24 * time.Hour)
			query = query.Where("created_at < ?", t)
		}
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * perPage
	query.Offset(offset).Limit(perPage).Order("created_at DESC").
		Preload("Order").Find(&transactions)

	// 格式化响应
	transactionList := make([]models.PointTransactionResponse, 0)
	for _, trans := range transactions {
		resp := models.PointTransactionResponse{
			ID:              trans.ID,
			TransactionType: string(trans.TransactionType),
			PointsChange:    trans.PointsChange,
			PointsBalance:   trans.PointsBalance,
			Description:     trans.Description,
			CreatedAt:       trans.CreatedAt,
		}

		if trans.OrderID != nil {
			resp.OrderID = trans.OrderID
			if trans.Order != nil {
				resp.OrderNumber = trans.Order.OrderNumber
			}
		}

		transactionList = append(transactionList, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"transactions": transactionList,
			"total":        total,
			"page":         page,
			"per_page":     perPage,
			"pages":        (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}
