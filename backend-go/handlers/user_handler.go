package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/middleware"
	"coffee-ordering-backend/models"
	"coffee-ordering-backend/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权",
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Preload("UserPoints").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user.ToResponse(),
	})
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("用户资料更新参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.BirthDate != nil {
		updates["birth_date"] = *req.BirthDate
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失败",
		})
		return
	}

	// 重新加载用户信息
	db.Preload("UserPoints").First(&user, userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    user.ToResponse(),
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权",
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 验证当前密码
	if !utils.CheckPassword(user.Password, req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "当前密码错误",
		})
		return
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "密码加密失败",
		})
		return
	}

	// 更新密码
	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "密码更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码修改成功",
	})
}


// GetUserOrders 获取用户订单列表
func GetUserOrders(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权",
		})
		return
	}

	db := database.GetDB()

	// 获取分页参数
	page := 1
	perPage := 10
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if val, err := strconv.Atoi(pp); err == nil && val > 0 && val <= 50 {
			perPage = val
		}
	}

	// 获取状态筛选
	status := c.Query("status")

	// 构建查询
	query := db.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var orders []models.Order
	offset := (page - 1) * perPage
	query.Preload("OrderItems.MenuItem").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&orders)

	// 格式化响应
	orderList := make([]gin.H, 0)
	for _, order := range orders {
		items := make([]gin.H, 0)
		totalPrice := 0.0
		for _, item := range order.OrderItems {
			menuName := ""
			if item.MenuItem.ID > 0 {
				menuName = item.MenuItem.Name
			}
			subtotal := item.GetSubtotal()
			totalPrice += subtotal
			items = append(items, gin.H{
				"id":         item.ID,
				"menu_id":    item.MenuID,
				"menu_name":  menuName,
				"quantity":   item.Quantity,
				"unit_price": item.UnitPrice,
				"subtotal":   subtotal,
			})
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
			"points_used":             order.CustomerPointsUsed,
			"points_earned":           order.PointsEarned,
			"status":                  order.Status,
			"notes":                   order.Notes,
			"items":                   items,
			"item_count":              len(items),
			"created_at":              order.CreatedAt,
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
