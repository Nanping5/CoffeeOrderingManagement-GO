package models

import (
	"time"
)

// MemberLevelConfig 会员等级配置模型
type MemberLevelConfig struct {
	ID                    uint        `gorm:"primaryKey" json:"id"`
	LevelName             MemberLevel `gorm:"type:enum('bronze','silver','gold','platinum');uniqueIndex;not null" json:"level_name"`
	LevelDisplayName      string      `gorm:"size:50;not null" json:"level_display_name"`
	MinPoints             int         `gorm:"not null;index" json:"min_points"`
	PointsEarningRate     float64     `gorm:"type:decimal(5,4);default:1.0000" json:"points_earning_rate"`
	PointsDiscountRate    float64     `gorm:"type:decimal(5,4);default:0.0000" json:"points_discount_rate"`
	MaxDiscountPercentage float64     `gorm:"type:decimal(5,2);default:20.00" json:"max_discount_percentage"`
	Benefits              string      `gorm:"type:json" json:"benefits"` // JSON格式
	LevelIcon             string      `gorm:"size:255" json:"level_icon"`
	IsActive              bool        `gorm:"default:true;index" json:"is_active"`
	SortOrder             int         `gorm:"default:0;index" json:"sort_order"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
}

// TableName 指定表名
func (MemberLevelConfig) TableName() string {
	return "member_levels"
}

// MemberLevelHistory 会员等级升级记录
type MemberLevelHistory struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	UserID          uint         `gorm:"not null;index" json:"user_id"`
	FromLevel       *MemberLevel `gorm:"type:enum('bronze','silver','gold','platinum')" json:"from_level"`
	ToLevel         MemberLevel  `gorm:"type:enum('bronze','silver','gold','platinum');not null" json:"to_level"`
	PointsAtUpgrade int          `gorm:"not null" json:"points_at_upgrade"`
	UpgradeReason   string       `gorm:"size:255" json:"upgrade_reason"`
	CreatedAt       time.Time    `gorm:"index" json:"created_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (MemberLevelHistory) TableName() string {
	return "member_level_history"
}
