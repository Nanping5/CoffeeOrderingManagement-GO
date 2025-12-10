package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMenuItems 获取菜单列表
func GetMenuItems(c *gin.Context) {
	// 获取查询参数
	availableOnly := c.DefaultQuery("available_only", "true") == "true"
	category := c.Query("category")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	db := database.GetDB()
	var items []models.MenuItem
	var total int64

	// 构建查询
	query := db.Model(&models.MenuItem{})

	if availableOnly {
		query = query.Where("is_available = ?", true)
	}

	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

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

// GetCategories 获取菜单分类
func GetCategories(c *gin.Context) {
	db := database.GetDB()
	var categories []string

	db.Model(&models.MenuItem{}).
		Where("is_available = ?", true).
		Distinct("category").
		Pluck("category", &categories)

	categoryList := []gin.H{
		{"value": "all", "label": "全部"},
	}

	for _, cat := range categories {
		categoryList = append(categoryList, gin.H{
			"value": cat,
			"label": cat,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categoryList,
	})
}

// GetMenuItem 获取单个菜单项详情
func GetMenuItem(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var item models.MenuItem

	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    item,
	})
}
