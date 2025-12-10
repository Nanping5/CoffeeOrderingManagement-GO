package models

import (
	"time"
)

// TransactionType 积分交易类型
type TransactionType string

const (
	TransactionTypeEarned        TransactionType = "earned"
	TransactionTypeUsed          TransactionType = "used"
	TransactionTypeExpired       TransactionType = "expired"
	TransactionTypeRefunded      TransactionType = "refunded"
	TransactionTypeSignupBonus   TransactionType = "signup_bonus"
	TransactionTypeBirthdayBonus TransactionType = "birthday_bonus"
	TransactionTypeReferralBonus TransactionType = "referral_bonus"
)

// PointTransaction 积分变动记录模型
type PointTransaction struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	UserID          uint            `gorm:"not null;index" json:"user_id"`
	OrderID         *uint           `gorm:"index" json:"order_id"`
	TransactionType TransactionType `gorm:"type:enum('earned','used','expired','refunded','signup_bonus','birthday_bonus','referral_bonus');not null" json:"transaction_type"`
	PointsChange    int             `gorm:"not null" json:"points_change"`      // 正数为获得，负数为使用
	PointsBalance   int             `gorm:"not null" json:"points_balance"`     // 变动后余额
	Description     string          `gorm:"size:255;not null" json:"description"`
	ReferenceID     string          `gorm:"size:100;index" json:"reference_id"`
	Metadata        *string         `gorm:"type:json" json:"metadata,omitempty"` // JSON格式额外数据
	ExpiresAt       *time.Time      `json:"expires_at"`
	IsActive        bool            `gorm:"default:true;index" json:"is_active"`
	CreatedAt       time.Time       `gorm:"index" json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// 关联
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (PointTransaction) TableName() string {
	return "point_transactions"
}

// PointTransactionResponse 积分交易响应
type PointTransactionResponse struct {
	ID              uint      `json:"id"`
	TransactionType string    `json:"transaction_type"`
	PointsChange    int       `json:"points_change"`
	PointsBalance   int       `json:"points_balance"`
	Description     string    `json:"description"`
	OrderID         *uint     `json:"order_id,omitempty"`
	OrderNumber     string    `json:"order_number,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
