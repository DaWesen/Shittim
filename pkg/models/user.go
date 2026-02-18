package models

import (
	"time"
)

// 用户基础信息模型
type User struct {
	ID        uint      `gorm:"primaryKey"`
	QQ        int64     `gorm:"uniqueIndex;comment:用户QQ号"`
	Nickname  string    `gorm:"size:50;comment:用户昵称"`
	Level     int       `gorm:"default:1;comment:用户等级"`
	Exp       int       `gorm:"default:0;comment:用户经验"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}
