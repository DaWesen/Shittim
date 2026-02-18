package models

import (
	"time"
)

// 签到记录模型
type Signin struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;comment:关联用户ID"`
	User      User      `gorm:"foreignKey:UserID;comment:关联用户"`
	QQ        int64     `gorm:"index;comment:用户QQ号（冗余，方便查询）"`
	Date      string    `gorm:"size:10;index;comment:签到日期（YYYY-MM-DD）"`
	Reward    int       `gorm:"default:0;comment:签到奖励"`
	Streak    int       `gorm:"default:1;comment:连续签到天数"`
	CreatedAt time.Time `gorm:"comment:签到时间"`
}
