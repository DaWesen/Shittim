package signin

import (
	"fmt"
	"time"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 检查用户今日是否已签到
func CheckSignin(qq int64) bool {
	var signin models.Signin
	today := time.Now().Format("2006-01-02")

	result := database.GetDB().Where("qq = ? AND date = ?", qq, today).First(&signin)
	return result.RowsAffected > 0
}

// 执行签到操作
func DoSignin(qq int64, nickname string) (int, int, error) {
	//检查是否已签到
	signed := CheckSignin(qq)
	if signed {
		return 0, 0, fmt.Errorf("今天已经签到过了")
	}

	//开始事务
	tx := database.GetDB().Begin()

	//查找或创建用户
	var user models.User
	if err := tx.Where("qq = ?", qq).FirstOrCreate(&user, models.User{
		QQ:       qq,
		Nickname: nickname,
	}).Error; err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	//计算连续签到天数
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var lastSignin models.Signin
	streak := 1

	if tx.Where("qq = ? AND date = ?", qq, yesterday).First(&lastSignin).Error == nil {
		streak = lastSignin.Streak + 1
	}

	//计算奖励（根据连续签到天数）
	baseReward := 10
	bonus := streak * 2
	totalReward := baseReward + bonus

	//创建签到记录
	signin := models.Signin{
		UserID: user.ID,
		QQ:     qq,
		Date:   time.Now().Format("2006-01-02"),
		Reward: totalReward,
		Streak: streak,
	}
	if err := tx.Create(&signin).Error; err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	//更新用户经验
	user.Exp += totalReward
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	return totalReward, streak, nil
}
