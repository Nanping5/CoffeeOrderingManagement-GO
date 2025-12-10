package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMenuItem 创建菜单项（管理员）
func CreateMenuItem(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Category    string  `json:"category" binding:"required"`
		ImageURL    string  `json:"image_url"`
		IsAvailable bool    `json:"is_available"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"请求参数错误: " + err.Error()},
		})
		return
	}

	db := database.GetDB()

	menuItem := models.MenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		IsAvailable: req.IsAvailable,
	}

	if err := db.Create(&menuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"创建菜单项失败: " + err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   true,
		"message":   "菜单项创建成功",
		"menu_item": menuItem,
	})
}

// UpdateMenuItem 更新菜单项（管理员）
func UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
		Category    *string  `json:"category"`
		ImageURL    *string  `json:"image_url"`
		IsAvailable *bool    `json:"is_available"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"请求参数错误: " + err.Error()},
		})
		return
	}

	db := database.GetDB()
	var menuItem models.MenuItem

	if err := db.First(&menuItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"菜单项不存在"},
		})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  []string{"价格必须大于0"},
			})
			return
		}
		updates["price"] = *req.Price
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	if err := db.Model(&menuItem).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"更新菜单项失败: " + err.Error()},
		})
		return
	}

	// 重新查询以获取更新后的数据
	db.First(&menuItem, id)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "菜单项更新成功",
		"menu_item": menuItem,
	})
}

// DeleteMenuItem 删除菜单项（管理员）
func DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var menuItem models.MenuItem

	if err := db.First(&menuItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"菜单项不存在"},
		})
		return
	}

	// 检查是否有关联的订单项
	var count int64
	db.Model(&models.OrderItem{}).Where("menu_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []string{"该菜单项有关联订单，无法删除"},
		})
		return
	}

	if err := db.Delete(&menuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"删除菜单项失败: " + err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "菜单项删除成功",
	})
}

// ToggleMenuItemAvailability 切换菜单项可用状态（管理员）
func ToggleMenuItemAvailability(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var menuItem models.MenuItem

	if err := db.First(&menuItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"errors":  []string{"菜单项不存在"},
		})
		return
	}

	// 切换状态
	menuItem.IsAvailable = !menuItem.IsAvailable
	if err := db.Save(&menuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  []string{"更新状态失败: " + err.Error()},
		})
		return
	}

	statusText := "上架"
	if !menuItem.IsAvailable {
		statusText = "下架"
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "菜单项" + statusText + "成功",
		"menu_item": menuItem,
	})
}

// GetAllMenuItemsAdmin 获取所有菜单项（管理员，包括下架商品）
func GetAllMenuItemsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	db := database.GetDB()
	var items []models.MenuItem
	var total int64

	query := db.Model(&models.MenuItem{})

	// 分类筛选
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页
	offset := (page - 1) * perPage
	query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":    items,
			"total":    total,
			"page":     page,
			"per_page": perPage,
			"pages":    (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}
