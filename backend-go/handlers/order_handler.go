package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/models"
	"coffee-ordering-backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Items        []OrderItemRequest `json:"items" binding:"required"`
	TotalPrice   float64            `json:"total_price" binding:"required"`
	Notes        string             `json:"notes"`
	UsePoints    bool               `json:"use_points"`     // 是否使用积分
	PointsToUse  int                `json:"points_to_use"`  // 使用的积分数量
}

// OrderItemRequest 订单项请求
type OrderItemRequest struct {
	MenuID    uint    `json:"menu_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price"`
}

// CreateOrder 创建订单（支持积分抵扣）
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"请求参数错误: " + err.Error()},
		})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"订单商品不能为空"},
		})
		return
	}

	db := database.GetDB()

	// 获取用户ID（如果已登录）
	userID, isLoggedIn := c.Get("user_id")
	var userIDPtr *uint
	var userPoints *models.UserPoints
	var memberLevel models.MemberLevel = models.MemberLevelBronze

	if isLoggedIn {
		uid := userID.(uint)
		userIDPtr = &uid

		// 获取用户积分信息
		userPoints = &models.UserPoints{}
		if err := db.Where("user_id = ?", uid).First(userPoints).Error; err == nil {
			memberLevel = userPoints.MemberLevel
		}
	}

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 生成订单号和取餐码
	orderNumber := models.GenerateOrderNumber()
	pickupCode := models.GeneratePickupCode()

	// 计算积分抵扣
	originalTotal := req.TotalPrice
	pointsDeduction := 0.0
	pointsUsed := 0
	finalPayment := originalTotal

	if req.UsePoints && req.PointsToUse > 0 && userPoints != nil {
		pointsService := services.NewPointsService()
		discount, err := pointsService.CalculatePointsDiscount(req.PointsToUse, memberLevel, originalTotal)
		if err == nil && userPoints.AvailablePoints >= req.PointsToUse {
			pointsDeduction = discount
			pointsUsed = req.PointsToUse
			finalPayment = originalTotal - pointsDeduction
		}
	}

	// 创建订单
	order := models.Order{
		UserID:                userIDPtr,
		OrderNumber:           orderNumber,
		PickupCode:            pickupCode,
		TotalPrice:            finalPayment,
		OriginalTotalPrice:    originalTotal,
		CustomerPointsUsed:    pointsUsed,
		PointsDeductionAmount: pointsDeduction,
		FinalPaymentAmount:    finalPayment,
		MemberLevelAtTime:     &memberLevel,
		Notes:                 req.Notes,
		Status:                models.OrderStatusPending,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"创建订单失败: " + err.Error()},
		})
		return
	}

	// 创建订单项
	for _, item := range req.Items {
		var menuItem models.MenuItem
		if err := tx.First(&menuItem, item.MenuID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  []string{"菜品不存在"},
			})
			return
		}

		unitPrice := item.UnitPrice
		if unitPrice == 0 {
			unitPrice = menuItem.Price
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			MenuID:    item.MenuID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
			Subtotal:  float64(item.Quantity) * unitPrice,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"errors":  []string{"创建订单项失败: " + err.Error()},
			})
			return
		}
	}

	// 处理积分
	estimatedPointsEarned := 0
	if userIDPtr != nil {
		pointsService := services.NewPointsService()

		// 使用积分
		if pointsUsed > 0 {
			description := fmt.Sprintf("积分抵扣 - 订单号: %s", orderNumber)
			if err := pointsService.UsePoints(tx, *userIDPtr, pointsUsed, order.ID, description); err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"errors":  []string{err.Error()},
				})
				return
			}
		}

		// 计算可获得积分（订单完成后才实际发放）
		earnedPoints, _ := pointsService.CalculateEarnedPoints(finalPayment, memberLevel)
		estimatedPointsEarned = earnedPoints
		order.PointsEarned = earnedPoints
		tx.Save(&order)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"提交订单失败: " + err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "订单创建成功",
		"order": gin.H{
			"id":                       order.ID,
			"order_number":             order.OrderNumber,
			"pickup_code":              order.PickupCode,
			"original_total_price":     originalTotal,
			"points_deduction_amount":  pointsDeduction,
			"final_payment_amount":     finalPayment,
			"points_used":              pointsUsed,
			"estimated_points_earned":  estimatedPointsEarned,
			"status":                   order.Status,
			"created_at":               order.CreatedAt,
		},
	})
}

// GetOrder 获取订单详情
func GetOrder(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var order models.Order

	if err := db.Preload("OrderItems.MenuItem").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"订单不存在"},
		})
		return
	}

	// 格式化订单项
	orderItems := make([]gin.H, 0)
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, gin.H{
			"id":         item.ID,
			"menu_id":    item.MenuID,
			"menu_name":  item.MenuItem.Name,
			"quantity":   item.Quantity,
			"unit_price": item.UnitPrice,
			"subtotal":   item.Subtotal,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"order": gin.H{
			"id":           order.ID,
			"order_number": order.OrderNumber,
			"pickup_code":  order.PickupCode,
			"total_price":  order.TotalPrice,
			"status":       order.Status,
			"notes":        order.Notes,
			"items":        orderItems,
			"created_at":   order.CreatedAt,
		},
	})
}

// GetOrderByPickupCode 通过取餐码查询订单
func GetOrderByPickupCode(c *gin.Context) {
	pickupCode := c.Param("pickup_code")

	db := database.GetDB()
	var order models.Order

	if err := db.Preload("OrderItems.MenuItem").Where("pickup_code = ?", pickupCode).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"取餐码不存在"},
		})
		return
	}

	// 格式化订单项
	orderItems := make([]gin.H, 0)
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, gin.H{
			"id":         item.ID,
			"menu_id":    item.MenuID,
			"menu_name":  item.MenuItem.Name,
			"quantity":   item.Quantity,
			"unit_price": item.UnitPrice,
			"subtotal":   item.Subtotal,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"order": gin.H{
			"id":           order.ID,
			"order_number": order.OrderNumber,
			"pickup_code":  order.PickupCode,
			"total_price":  order.TotalPrice,
			"status":       order.Status,
			"notes":        order.Notes,
			"items":        orderItems,
			"created_at":   order.CreatedAt,
		},
	})
}


// CalculatePointsRequest 积分计算请求
type CalculatePointsRequest struct {
	Items       []OrderItemRequest `json:"items" binding:"required"`
	PointsToUse int                `json:"points_to_use"`
}

// CalculatePointsForOrder 积分使用预估
func CalculatePointsForOrder(c *gin.Context) {
	var req CalculatePointsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	db := database.GetDB()

	// 获取用户积分信息
	var userPoints models.UserPoints
	if err := db.Where("user_id = ?", userID.(uint)).First(&userPoints).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "积分账户不存在",
		})
		return
	}

	// 计算订单总价
	var originalTotal float64
	for _, item := range req.Items {
		var menuItem models.MenuItem
		if err := db.First(&menuItem, item.MenuID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "菜品不存在",
			})
			return
		}

		unitPrice := item.UnitPrice
		if unitPrice == 0 {
			unitPrice = menuItem.Price
		}
		originalTotal += float64(item.Quantity) * unitPrice
	}

	pointsService := services.NewPointsService()

	// 获取最大可用积分
	maxUsablePoints, _ := pointsService.GetMaxUsablePoints(originalTotal, userPoints.MemberLevel, userPoints.AvailablePoints)

	// 计算积分抵扣
	pointsToUse := req.PointsToUse
	if pointsToUse > maxUsablePoints {
		pointsToUse = maxUsablePoints
	}

	pointsValue := 0.0
	if pointsToUse > 0 {
		pointsValue, _ = pointsService.CalculatePointsDiscount(pointsToUse, userPoints.MemberLevel, originalTotal)
	}

	finalTotal := originalTotal - pointsValue

	// 计算可获得积分
	estimatedPointsEarned, _ := pointsService.CalculateEarnedPoints(finalTotal, userPoints.MemberLevel)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_total":          originalTotal,
			"max_usable_points":       maxUsablePoints,
			"points_to_use":           pointsToUse,
			"points_value":            pointsValue,
			"final_total":             finalTotal,
			"estimated_points_earned": estimatedPointsEarned,
			"user_points_balance":     userPoints.AvailablePoints,
			"points_after_usage":      userPoints.AvailablePoints - pointsToUse,
		},
	})
}
