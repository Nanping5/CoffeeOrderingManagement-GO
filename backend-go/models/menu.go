package models

import (
	"time"
)

// MenuItem 菜单项模型
type MenuItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Category    string    `gorm:"size:50;not null;default:'coffee';index" json:"category"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	IsAvailable bool      `gorm:"default:true;not null;index" json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	OrderItems []OrderItem `gorm:"foreignKey:MenuID;constraint:OnDelete:RESTRICT" json:"-"`
}

// TableName 指定表名
func (MenuItem) TableName() string {
	return "menu_items"
}
