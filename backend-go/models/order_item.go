package models

import (
	"time"
)

// OrderItem 订单项模型
type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"not null;index" json:"order_id"`
	MenuID    uint      `gorm:"not null;index" json:"menu_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal  float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	MenuItem MenuItem `gorm:"foreignKey:MenuID" json:"menu_item,omitempty"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}
