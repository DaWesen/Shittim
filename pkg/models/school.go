package models

import "time"

//学院名称
type School struct {
	SchoolName    string    `gorm:"primaryKey;size:50;comment:学校名称"`
	StudentCounts int       `gorm:"default:0;comment:学生数量"`
	Students      []Student `gorm:"foreignKey:SchoolName;comment:所属学生"`
	Clubs         []Club    `gorm:"foreignKey:SchoolName;comment:所属社团"`
	CreatedAt     time.Time `gorm:"comment:创建时间"`
	UpdatedAt     time.Time `gorm:"comment:更新时间"`
}
