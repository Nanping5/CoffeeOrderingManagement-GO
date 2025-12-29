package handlers

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/models"
	"coffee-ordering-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req models.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
			"errors":  err.Error(),
		})
		return
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "邮箱已被注册",
		})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "密码加密失败",
		})
		return
	}

	// 创建用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		IsActive:  true,
	}

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户创建失败",
		})
		return
	}

	// 检查积分账户是否已存在（触发器可能已创建）
	var userPoints models.UserPoints
	if err := tx.Where("user_id = ?", user.ID).First(&userPoints).Error; err != nil {
		// 不存在则创建（触发器应该已经创建，但以防触发器失败）
		userPoints = models.UserPoints{
			UserID:         user.ID,
			TotalPoints:    50, // 注册奖励50积分
			LifetimePoints: 50,
			MemberLevel:    models.MemberLevelBronze,
		}
		if err := tx.Create(&userPoints).Error; err != nil {
			// 如果创建失败，检查是否是触发器已经创建
			if err := tx.Where("user_id = ?", user.ID).First(&userPoints).Error; err == nil {
				// 触发器已创建，继续使用
			} else {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "积分账户创建失败",
				})
				return
			}
		}
	}

	// 更新积分账户（无论是否是新创建的，都设置为初始积分）
	userPoints.TotalPoints = 50
	userPoints.LifetimePoints = 50
	if err := tx.Save(&userPoints).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "积分账户更新失败",
		})
		return
	}

	// 记录注册奖励积分
	pointTransaction := models.PointTransaction{
		UserID:          user.ID,
		TransactionType: models.TransactionTypeSignupBonus,
		PointsChange:    50,
		PointsBalance:   50,
		Description:     "注册奖励",
	}

	if err := tx.Create(&pointTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "积分记录创建失败",
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "注册失败",
		})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "令牌生成失败",
		})
		return
	}

	// 加载用户积分信息
	db.Preload("UserPoints").First(&user, user.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "注册成功，获得50积分奖励",
		"data": gin.H{
			"user":  user.ToResponse(),
			"token": token,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	db := database.GetDB()

	// 查找用户
	var user models.User
	if err := db.Preload("UserPoints").Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "邮箱或密码错误",
		})
		return
	}

	// 检查账户状态
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "账户已被禁用",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "邮箱或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "令牌生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user":       user.ToResponse(),
			"token":      token,
			"expires_in": 7200, // 2小时
		},
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 简化版本，实际应该将token加入黑名单
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未提供认证令牌",
		})
		return
	}

	parts := []string{}
	if authHeader != "" {
		parts = []string{authHeader[:6], authHeader[7:]}
	}
	if len(parts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "令牌格式错误",
		})
		return
	}

	newToken, err := utils.RefreshToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "令牌刷新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "令牌刷新成功",
		"data": gin.H{
			"token":      newToken,
			"expires_in": 7200,
		},
	})
}


// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	db := database.GetDB()

	// 查找管理员用户
	var user models.User
	if err := db.Where("username = ? AND role = ?", req.Username, "admin").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查账户状态
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "账户已被禁用",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "令牌生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
			"token":      token,
			"expires_in": 7200,
		},
	})
}
