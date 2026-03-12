package dao

import (
	"Shittim/pkg/database"
	"Shittim/pkg/models"
	"time"
)

// MemoryDAO 记忆数据访问对象
type MemoryDAO struct{}

// NewMemoryDAO 创建记忆数据访问对象
func NewMemoryDAO() *MemoryDAO {
	return &MemoryDAO{}
}

// SaveExclusiveMemory 保存独家记忆
func (dao *MemoryDAO) SaveExclusiveMemory(studentName, content, emotionTag string) error {
	memory := models.ExclusiveMemory{
		StoryBase: models.StoryBase{
			Title:     studentName + "的独家记忆",
			Content:   content,
			StoryType: "exclusive",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		StudentName:   studentName,
		IntimacyLevel: 1,
		MemoryType:    "chat",
		EmotionTag:    emotionTag,
		LastRecallAt:  time.Now(),
	}
	return database.GetDB().Create(&memory).Error
}

// LoadExclusiveMemories 加载独家记忆
func (dao *MemoryDAO) LoadExclusiveMemories(studentName string) []models.ExclusiveMemory {
	var memories []models.ExclusiveMemory
	database.GetDB().Where("student_name = ?", studentName).Find(&memories)
	return memories
}

// SaveDailyStory 保存日常故事
func (dao *MemoryDAO) SaveDailyStory(students []string, content string) error {
	story := models.DailyStory{
		StoryBase: models.StoryBase{
			Title:     "日常故事",
			Content:   content,
			StoryType: "daily",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Theme:   "日常",
		Setting: "校园",
		Mood:    "轻松",
	}
	return database.GetDB().Create(&story).Error
}

// SaveEventStory 保存活动故事
func (dao *MemoryDAO) SaveEventStory(students []string, content string) error {
	story := models.EventStory{
		StoryBase: models.StoryBase{
			Title:     "活动故事",
			Content:   content,
			StoryType: "event",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		EventName:    "活动",
		StudentCount: len(students),
		Duration:     "一天",
		EventType:    "校园活动",
		Location:     "学校",
	}
	return database.GetDB().Create(&story).Error
}
