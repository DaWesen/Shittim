package models

import "time"

// 独家记忆模型，存储用户与学生一对一的剧场故事
type ExclusiveMemory struct {
	StoryBase
	StudentName   string    `gorm:"size:50;index;comment:相关学生姓名"`
	IntimacyLevel int       `gorm:"default:1;comment:亲密程度等级"`
	MemoryType    string    `gorm:"size:50;comment:记忆类型"`
	SpecialDate   string    `gorm:"size:50;comment:特殊日期"`
	EmotionTag    string    `gorm:"size:100;comment:情感标签，逗号分隔"`
	LastRecallAt  time.Time `gorm:"comment:最后回忆时间"`
}
