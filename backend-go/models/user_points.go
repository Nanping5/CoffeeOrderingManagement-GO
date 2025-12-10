package models

import (
	"time"
)

// MemberLevel 会员等级类型
type MemberLevel string

const (
	MemberLevelBronze   MemberLevel = "bronze"
	MemberLevelSilver   MemberLevel = "silver"
	MemberLevelGold     MemberLevel = "gold"
	MemberLevelPlatinum MemberLevel = "platinum"
)

// UserPoints 用户积分模型
type UserPoints struct {
	ID               uint        `gorm:"primaryKey" json:"id"`
	UserID           uint        `gorm:"uniqueIndex;not null" json:"user_id"`
	TotalPoints      int         `gorm:"default:0" json:"total_points"`           // 当前积分余额
	AvailablePoints  int         `gorm:"default:0" json:"available_points"`       // 可用积分
	FrozenPoints     int         `gorm:"default:0" json:"frozen_points"`          // 冻结积分
	LifetimePoints   int         `gorm:"default:0" json:"lifetime_points"`        // 历史累计积分
	MemberLevel      MemberLevel `gorm:"type:enum('bronze','silver','gold','platinum');default:'bronze'" json:"member_level"`
	LevelUpgradeDate *time.Time  `json:"level_upgrade_date"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserPoints) TableName() string {
	return "user_points"
}

// UserPointsResponse 积分信息响应
type UserPointsResponse struct {
	TotalPoints      int                    `json:"total_points"`
	AvailablePoints  int                    `json:"available_points"`
	FrozenPoints     int                    `json:"frozen_points"`
	LifetimePoints   int                    `json:"lifetime_points"`
	MemberLevel      string                 `json:"member_level"`
	NextLevel        *NextLevelInfo         `json:"next_level,omitempty"`
	LevelBenefits    *MemberLevelBenefits   `json:"level_benefits"`
}

// NextLevelInfo 下一等级信息
type NextLevelInfo struct {
	Name           string `json:"name"`
	RequiredPoints int    `json:"required_points"`
	PointsNeeded   int    `json:"points_needed"`
}

// MemberLevelBenefits 会员等级权益
type MemberLevelBenefits struct {
	PointsEarningRate   float64  `json:"points_earning_rate"`
	MaxDiscountRate     float64  `json:"max_discount_rate"`
	Benefits            []string `json:"benefits"`
}
