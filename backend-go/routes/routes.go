package routes

import (
	"coffee-ordering-backend/handlers"
	"coffee-ordering-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "咖啡自助点餐系统API",
			"version": "1.0.0-go",
			"mode":    "自助点餐模式",
			"endpoints": gin.H{
				"menu":   "/api/menu",
				"orders": "/api/orders",
				"admin":  "/api/admin",
			},
		})
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
			"app":      "running",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 认证路由（公开）
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/refresh", handlers.RefreshToken)
			auth.POST("/admin/login", handlers.AdminLogin)
		}

		// 菜单路由（公开）
		menu := api.Group("/menu")
		{
			menu.GET("", handlers.GetMenuItems)
			menu.GET("/categories", handlers.GetCategories)
			menu.GET("/:id", handlers.GetMenuItem)
		}

		// 订单路由（支持可选认证）
		orders := api.Group("/orders")
		orders.Use(middleware.OptionalUserAuth())
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("/:id", handlers.GetOrder)
			orders.GET("/pickup/:pickup_code", handlers.GetOrderByPickupCode)
			orders.POST("/points-calculation", middleware.UserAuthRequired(), handlers.CalculatePointsForOrder)
		}

		// 用户路由（需要认证）
		user := api.Group("/user")
		user.Use(middleware.UserAuthRequired())
		{
			user.GET("/profile", handlers.GetUserProfile)
			user.PUT("/profile", handlers.UpdateUserProfile)
			user.PUT("/password", handlers.ChangePassword)

			// 订单相关
			user.GET("/orders", handlers.GetUserOrders)

			// 积分相关
			user.GET("/points", handlers.GetUserPoints)
			user.GET("/points/transactions", handlers.GetPointTransactions)
		}

		// 管理员路由（需要认证）
		admin := api.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			// 菜单管理
			adminMenu := admin.Group("/menu")
			{
				adminMenu.GET("", handlers.GetAllMenuItemsAdmin)
				adminMenu.POST("", handlers.CreateMenuItem)
				adminMenu.PUT("/:id", handlers.UpdateMenuItem)
				adminMenu.DELETE("/:id", handlers.DeleteMenuItem)
				adminMenu.PATCH("/:id/toggle", handlers.ToggleMenuItemAvailability)
			}

			// 订单管理
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", handlers.GetAllOrders)
				adminOrders.PUT("/:id/status", handlers.UpdateOrderStatus)
				adminOrders.DELETE("/:id", handlers.DeleteOrder)
				adminOrders.GET("/statistics", handlers.GetOrderStatistics)
			}
		}
	}
}
