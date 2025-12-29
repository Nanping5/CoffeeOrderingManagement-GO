package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	Username                string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email                   string         `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password                string         `gorm:"size:255;not null" json:"-"` // 不返回密码
	Phone                   string         `gorm:"size:20" json:"phone"`
	Role                    string         `gorm:"size:20;default:'user'" json:"role"` // admin, user
	FirstName               string         `gorm:"size:50" json:"first_name"`
	LastName                string         `gorm:"size:50" json:"last_name"`
	AvatarURL               string         `gorm:"size:255" json:"avatar_url"`
	Gender                  *string        `gorm:"type:enum('male','female','other')" json:"gender"`
	BirthDate               *time.Time     `json:"birth_date"`
	IsVerified              bool           `gorm:"default:false" json:"is_verified"`
	VerificationToken       string         `gorm:"size:255" json:"-"`
	VerificationExpiresAt   *time.Time     `json:"-"`
	ResetPasswordToken      string         `gorm:"size:255" json:"-"`
	ResetPasswordExpiresAt  *time.Time     `json:"-"`
	IsActive                bool           `gorm:"default:true" json:"is_active"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`

	// 关联
	UserPoints       *UserPoints        `gorm:"foreignKey:UserID" json:"user_points,omitempty"`
	Orders           []Order            `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	PointTransactions []PointTransaction `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Phone     *string    `json:"phone"`
	Gender    *string    `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
	AvatarURL *string    `json:"avatar_url"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	FullName    string     `json:"full_name"`
	AvatarURL   string     `json:"avatar_url"`
	Gender      *string    `json:"gender"`
	BirthDate   *time.Time `json:"birth_date"`
	IsVerified  bool       `json:"is_verified"`
	IsActive    bool       `json:"is_active"`
	MemberLevel string     `json:"member_level,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	fullName := u.FirstName
	if u.LastName != "" {
		if fullName != "" {
			fullName += " "
		}
		fullName += u.LastName
	}

	resp := UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Phone:      u.Phone,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		FullName:   fullName,
		AvatarURL:  u.AvatarURL,
		Gender:     u.Gender,
		BirthDate:  u.BirthDate,
		IsVerified: u.IsVerified,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
	}

	if u.UserPoints != nil {
		resp.MemberLevel = string(u.UserPoints.MemberLevel)
	}

	return resp
}
