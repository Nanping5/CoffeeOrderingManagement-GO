package main

import (
	"coffee-ordering-backend/config"
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/routes"
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.GinMode)

	// 初始化数据库
	database.InitDB()

	// 创建 Gin 引擎
	r := gin.Default()

	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	origins := strings.Split(config.AppConfig.CORSOrigins, ",")
	corsConfig.AllowOrigins = origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := ":" + config.AppConfig.Port
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
