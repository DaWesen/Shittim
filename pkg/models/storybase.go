package models

import "time"

// StoryBase 故事基础模型，作为所有故事类型的父模型
type StoryBase struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:255;not null;comment:故事标题"`
	Content     string    `gorm:"type:text;not null;comment:故事内容"`
	Students    []Student `gorm:"many2many:story_students;"` // 与学生的多对多关联
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间"`
	IsUserSaved bool      `gorm:"default:false;comment:是否用户确认保存"`
	StoryType   string    `gorm:"size:50;comment:故事类型"` // 日常故事、活动故事、独家记忆
}
