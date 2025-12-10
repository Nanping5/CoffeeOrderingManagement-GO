package models

import (
	"fmt"
	"math/rand"
	"time"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order 订单模型
type Order struct {
	ID                    uint         `gorm:"primaryKey" json:"id"`
	UserID                *uint        `gorm:"index" json:"user_id"` // 用户ID，可为空（游客订单）
	OrderNumber           string       `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	PickupCode            string       `gorm:"size:10;index;not null" json:"pickup_code"`
	TotalPrice            float64      `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status                OrderStatus  `gorm:"type:enum('pending','preparing','ready','completed','cancelled');default:'pending';index" json:"status"`
	Notes                 string       `gorm:"type:text" json:"notes"`
	
	// 积分相关字段
	CustomerPointsUsed    int          `gorm:"default:0;index" json:"customer_points_used"`           // 使用的积分数量
	PointsDeductionAmount float64      `gorm:"type:decimal(10,2);default:0.00" json:"points_deduction_amount"` // 积分抵扣金额
	PointsEarned          int          `gorm:"default:0;index" json:"points_earned"`                  // 获得的积分数量
	OriginalTotalPrice    float64      `gorm:"type:decimal(10,2)" json:"original_total_price"`        // 原始总价（积分抵扣前）
	FinalPaymentAmount    float64      `gorm:"type:decimal(10,2)" json:"final_payment_amount"`        // 实际支付金额
	MemberLevelAtTime     *MemberLevel `gorm:"type:enum('bronze','silver','gold','platinum')" json:"member_level_at_time"` // 下单时会员等级
	
	CreatedAt             time.Time    `gorm:"index" json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`

	// 关联
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// GenerateOrderNumber 生成订单号
func GenerateOrderNumber() string {
	datePrefix := time.Now().Format("20060102")
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("CO%s%04d", datePrefix, randomNum)
}

// GeneratePickupCode 生成取餐码
func GeneratePickupCode() string {
	hour := time.Now().Hour()
	letter := rune(65 + (hour % 26)) // A-Z
	randomNum := rand.Intn(1000)
	return fmt.Sprintf("%c%03d", letter, randomNum)
}
