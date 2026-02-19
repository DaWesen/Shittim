package models

// EventStory 活动故事模型，存储多个学生和用户长时间交互的故事
type EventStory struct {
	StoryBase
	EventName      string `gorm:"size:100;comment:活动名称"`
	StudentCount   int    `gorm:"comment:参与学生数量"`
	Duration       string `gorm:"size:50;comment:活动持续时间"`
	EventType      string `gorm:"size:50;comment:活动类型"`
	Location       string `gorm:"size:100;comment:活动地点"`
	ChapterCount   int    `gorm:"default:1;comment:章节数量"`
	CurrentChapter int    `gorm:"default:1;comment:当前章节"`
}
