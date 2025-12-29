package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllOrders 获取所有订单（管理员）
func GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := database.GetDB()
	var orders []models.Order
	var total int64

	query := db.Model(&models.Order{})

	// 状态筛选
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// 日期范围筛选
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			// 加一天以包含结束日期的全天
			t = t.Add(24 * time.Hour)
			query = query.Where("created_at < ?", t)
		}
	}

	// 获取总数
	query.Count(&total)

	// 分页
	offset := (page - 1) * perPage
	query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders)

	// 格式化订单数据
	orderList := make([]gin.H, 0)
	for _, order := range orders {
		// 获取订单项数量和计算总价
		var itemCount int64
		var orderItems []models.OrderItem
		db.Model(&models.OrderItem{}).Where("order_id = ?", order.ID).Find(&orderItems)
		itemCount = int64(len(orderItems))

		// 计算订单总价
		totalPrice := 0.0
		for _, item := range orderItems {
			totalPrice += item.GetSubtotal()
		}

		// 计算最终支付金额
		finalPrice := totalPrice - order.PointsDeductionAmount

		orderList = append(orderList, gin.H{
			"id":                      order.ID,
			"order_number":            order.OrderNumber,
			"pickup_code":             order.PickupCode,
			"original_total_price":    totalPrice,
			"points_deduction_amount": order.PointsDeductionAmount,
			"final_payment_amount":    finalPrice,
			"status":                  order.Status,
			"notes":                   order.Notes,
			"item_count":              itemCount,
			"created_at":              order.CreatedAt,
			"updated_at":              order.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"orders":   orderList,
			"total":    total,
			"page":     page,
			"per_page": perPage,
			"pages":    (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// UpdateOrderStatus 更新订单状态（管理员）
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"请求参数错误: " + err.Error()},
		})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{
		"pending":   true,
		"preparing": true,
		"ready":     true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"无效的订单状态"},
		})
		return
	}

	db := database.GetDB()
	var order models.Order

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"订单不存在"},
		})
		return
	}

	// 更新状态
	order.Status = models.OrderStatus(req.Status)
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"更新订单状态失败: " + err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订单状态更新成功",
		"order": gin.H{
			"id":           order.ID,
			"order_number": order.OrderNumber,
			"pickup_code":  order.PickupCode,
			"status":       order.Status,
			"updated_at":   order.UpdatedAt,
		},
	})
}

// GetOrderStatistics 获取订单统计（管理员）
func GetOrderStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := database.GetDB()
	query := db.Model(&models.Order{})

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

	// 总订单数
	var totalOrders int64
	query.Count(&totalOrders)

	// 各状态订单数
	var statusCounts []struct {
		Status string
		Count  int64
	}
	db.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts)

	statusMap := make(map[string]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}

	// 总收入（所有非取消订单）
	var totalRevenue float64
	db.Model(&models.Order{}).
		Where("status != ?", "cancelled").
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&totalRevenue)

	// 今日订单数
	today := time.Now().Truncate(24 * time.Hour)
	var todayOrders int64
	db.Model(&models.Order{}).
		Where("created_at >= ?", today).
		Count(&todayOrders)

	// 今日收入（所有非取消订单）
	var todayRevenue float64
	db.Model(&models.Order{}).
		Where("status != ? AND created_at >= ?", "cancelled", today).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&todayRevenue)

	// 会员统计
	var totalUsers int64
	db.Model(&models.User{}).Where("role = ?", "user").Count(&totalUsers)

	// 各等级会员数
	var levelCounts []struct {
		MemberLevel string
		Count       int64
	}
	db.Model(&models.UserPoints{}).
		Select("member_level, COUNT(*) as count").
		Group("member_level").
		Scan(&levelCounts)

	memberLevelMap := make(map[string]int64)
	for _, lc := range levelCounts {
		memberLevelMap[lc.MemberLevel] = lc.Count
	}

	// 热销商品Top5 - 使用小写json字段名
	type TopProduct struct {
		MenuID   uint    `json:"menu_id"`
		MenuName string  `json:"menu_name"`
		Quantity int64   `json:"quantity"`
		Revenue  float64 `json:"revenue"`
	}
	var topProducts []TopProduct
	db.Table("order_items").
		Select("order_items.menu_item_id, menu_items.name as menu_name, SUM(order_items.quantity) as quantity, SUM(order_items.subtotal) as revenue").
		Joins("LEFT JOIN menu_items ON order_items.menu_item_id = menu_items.id").
		Group("order_items.menu_item_id, menu_items.name").
		Order("quantity DESC").
		Limit(5).
		Scan(&topProducts)

	// 最近7天订单趋势 - 使用小写json字段名
	type DailyOrder struct {
		Date    string  `json:"date"`
		Count   int64   `json:"count"`
		Revenue float64 `json:"revenue"`
	}
	var dailyOrders []DailyOrder
	sevenDaysAgo := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)
	db.Model(&models.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as count, COALESCE(SUM(total_price), 0) as revenue").
		Where("created_at >= ?", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyOrders)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_orders":    totalOrders,
			"total_revenue":   totalRevenue,
			"today_orders":    todayOrders,
			"today_revenue":   todayRevenue,
			"status_counts":   statusMap,
			"pending_count":   statusMap["pending"],
			"preparing_count": statusMap["preparing"],
			"ready_count":     statusMap["ready"],
			"completed_count": statusMap["completed"],
			"cancelled_count": statusMap["cancelled"],
			"total_users":     totalUsers,
			"member_levels":   memberLevelMap,
			"top_products":    topProducts,
			"daily_orders":    dailyOrders,
		},
	})
}

// DeleteOrder 删除订单（管理员）
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var order models.Order

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"订单不存在"},
		})
		return
	}

	// 删除订单（会级联删除订单项）
	if err := db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"删除订单失败: " + err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订单删除成功",
	})
}
