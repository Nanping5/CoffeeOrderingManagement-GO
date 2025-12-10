package database

import (
	"coffee-ordering-backend/config"
	"coffee-ordering-backend/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	dsn := config.AppConfig.GetDSN()

	// 配置 GORM
	gormConfig := &gorm.Config{}
	if config.AppConfig.GinMode == "debug" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移（可选，如果数据库已存在则注释掉）
	// AutoMigrate()
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.MenuItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
