package models

import "time"

// DailyStory 日常故事模型，存储日常化的青春故事
type DailyStory struct {
	StoryBase
	Theme       string    `gorm:"size:100;comment:故事主题"`
	Setting     string    `gorm:"size:100;comment:故事场景"`
	Mood        string    `gorm:"size:50;comment:故事氛围"`
	LastReadAt  time.Time `gorm:"comment:最后阅读时间"`
}
