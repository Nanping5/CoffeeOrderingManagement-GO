package middleware

import (
	"coffee-ordering-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminRequired 管理员权限验证中间件（JWT版）
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")

		// 检查是否有 token
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}

		// 提取 token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的认证信息",
			})
			c.Abort()
			return
		}

		// 兼容旧的假token（用于开发测试）
		if strings.HasPrefix(token, "fake-admin-token-") {
			c.Set("user_id", uint(0))
			c.Set("username", "admin")
			c.Set("role", "admin")
			c.Next()
			return
		}

		// JWT验证
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 检查是否是管理员
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
