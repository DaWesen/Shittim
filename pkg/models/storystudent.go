package models

// 故事与学生的多对多关联模型
type StoryStudent struct {
	StoryID     uint    `gorm:"primaryKey"`
	StudentName string  `gorm:"primaryKey;size:50"`
	Student     Student `gorm:"foreignKey:StudentName;references:Name"`
	Role        string  `gorm:"size:50;comment:在故事中的角色"`
	Importance  int     `gorm:"default:1;comment:重要性等级"`
}
