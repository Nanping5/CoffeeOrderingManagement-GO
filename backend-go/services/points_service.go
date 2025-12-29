package services

import (
	"coffee-ordering-backend/database"
	"coffee-ordering-backend/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PointsService 积分服务
type PointsService struct{}

// CalculatePointsDiscount 计算积分可抵扣金额
func (s *PointsService) CalculatePointsDiscount(pointsToUse int, memberLevel models.MemberLevel, orderTotal float64) (float64, error) {
	if pointsToUse < 100 {
		return 0, errors.New("最少使用100积分")
	}

	db := database.GetDB()
	var levelConfig models.MemberLevelConfig

	if err := db.Where("level_name = ? AND is_active = ?", memberLevel, true).First(&levelConfig).Error; err != nil {
		return 0, errors.New("会员等级配置不存在")
	}

	// 计算最大抵扣金额
	maxDiscountAmount := orderTotal * (levelConfig.MaxDiscountPercentage / 100.0)

	// 计算积分价值（100积分=1元）
	pointsValue := float64(pointsToUse) / 100.0

	// 返回较小值
	if pointsValue <= maxDiscountAmount {
		return pointsValue, nil
	}

	return maxDiscountAmount, nil
}

// CalculateEarnedPoints 计算订单可获得积分
func (s *PointsService) CalculateEarnedPoints(paymentAmount float64, memberLevel models.MemberLevel) (int, error) {
	db := database.GetDB()
	var levelConfig models.MemberLevelConfig

	if err := db.Where("level_name = ? AND is_active = ?", memberLevel, true).First(&levelConfig).Error; err != nil {
		return 0, errors.New("会员等级配置不存在")
	}

	// 基础积分：每消费1元获得1积分
	basePoints := int(paymentAmount)

	// 应用会员等级倍率
	earnedPoints := int(float64(basePoints) * levelConfig.PointsEarningRate)

	return earnedPoints, nil
}

// UsePoints 使用积分
func (s *PointsService) UsePoints(tx *gorm.DB, userID uint, pointsToUse int, orderID uint, description string) error {
	var userPoints models.UserPoints

	if err := tx.Where("user_id = ?", userID).First(&userPoints).Error; err != nil {
		return errors.New("积分账户不存在")
	}

	// 检查可用积分
	if userPoints.TotalPoints < pointsToUse {
		return fmt.Errorf("可用积分不足，当前可用: %d", userPoints.TotalPoints)
	}

	// 扣减积分
	userPoints.TotalPoints -= pointsToUse

	if err := tx.Save(&userPoints).Error; err != nil {
		return errors.New("积分扣减失败")
	}

	// 记录积分变动
	transaction := models.PointTransaction{
		UserID:          userID,
		OrderID:         &orderID,
		TransactionType: models.TransactionTypeUsed,
		PointsChange:    -pointsToUse,
		PointsBalance:   userPoints.TotalPoints,
		Description:     description,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		return errors.New("积分记录创建失败")
	}

	return nil
}

// EarnPoints 获得积分
func (s *PointsService) EarnPoints(tx *gorm.DB, userID uint, pointsToEarn int, orderID *uint, transactionType models.TransactionType, description string) error {
	var userPoints models.UserPoints

	if err := tx.Where("user_id = ?", userID).First(&userPoints).Error; err != nil {
		return errors.New("积分账户不存在")
	}

	// 增加积分
	userPoints.TotalPoints += pointsToEarn
	userPoints.LifetimePoints += pointsToEarn

	// 检查是否需要升级会员等级
	newLevel := s.CalculateMemberLevel(userPoints.LifetimePoints)
	if newLevel != userPoints.MemberLevel {
		userPoints.MemberLevel = newLevel
		now := time.Now()
		userPoints.LevelUpgradeDate = &now
	}

	if err := tx.Save(&userPoints).Error; err != nil {
		return errors.New("积分增加失败")
	}

	// 记录积分变动
	transaction := models.PointTransaction{
		UserID:          userID,
		OrderID:         orderID,
		TransactionType: transactionType,
		PointsChange:    pointsToEarn,
		PointsBalance:   userPoints.TotalPoints,
		Description:     description,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		return errors.New("积分记录创建失败")
	}

	return nil
}

// CalculateMemberLevel 根据累计积分计算会员等级
func (s *PointsService) CalculateMemberLevel(lifetimePoints int) models.MemberLevel {
	if lifetimePoints >= 20000 {
		return models.MemberLevelPlatinum
	} else if lifetimePoints >= 5000 {
		return models.MemberLevelGold
	} else if lifetimePoints >= 1000 {
		return models.MemberLevelSilver
	}
	return models.MemberLevelBronze
}

// GetMaxUsablePoints 获取订单最大可用积分
func (s *PointsService) GetMaxUsablePoints(orderTotal float64, memberLevel models.MemberLevel, availablePoints int) (int, error) {
	db := database.GetDB()
	var levelConfig models.MemberLevelConfig

	if err := db.Where("level_name = ? AND is_active = ?", memberLevel, true).First(&levelConfig).Error; err != nil {
		return 0, errors.New("会员等级配置不存在")
	}

	// 计算最大抵扣金额
	maxDiscountAmount := orderTotal * (levelConfig.MaxDiscountPercentage / 100.0)

	// 转换为积分（1元=100积分）
	maxUsablePoints := int(maxDiscountAmount * 100)

	// 取用户可用积分和最大可用积分的较小值
	if availablePoints < maxUsablePoints {
		return availablePoints, nil
	}

	return maxUsablePoints, nil
}

// NewPointsService 创建积分服务实例
func NewPointsService() *PointsService {
	return &PointsService{}
}
