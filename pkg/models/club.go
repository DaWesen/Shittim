package models

import "time"

//学生社团
type Club struct {
	ClubName      string    `gorm:"primaryKey;size:50;comment:社团名称"`
	StudentCounts int       `gorm:"default:0;comment:学生数量"`
	SchoolName    string    `gorm:"index;size:50;comment:所属学校名称"`
	School        School    `gorm:"foreignKey:SchoolName;comment:所属学校"`
	Students      []Student `gorm:"foreignKey:ClubName;comment:所属学生"`
	CreatedAt     time.Time `gorm:"comment:创建时间"`
	UpdatedAt     time.Time `gorm:"comment:更新时间"`
	DeletedAt     time.Time `gorm:"index;comment:软删除时间"`
}
