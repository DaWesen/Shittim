package models

import "time"

//学生档案
type Student struct {
	Name       string    `gorm:"primaryKey;size:50;comment:学生姓名"`
	Level      string    `gorm:"size:20;comment:学生年级"`
	Age        uint      `gorm:"comment:学生年龄"`
	ClubName   string    `gorm:"index;size:50;comment:所属社团名称"`
	Club       Club      `gorm:"foreignKey:ClubName;comment:所属社团"`
	SchoolName string    `gorm:"index;size:50;comment:所属学校名称"`
	School     School    `gorm:"foreignKey:SchoolName;comment:所属学校"`
	Height     string    `gorm:"size:20;comment:学生身高"`
	Love       string    `gorm:"size:50;comment:爱好"`
	UnderEye   string    `gorm:"size:20;comment:是否需要重点关注"`
	Affection  string    `gorm:"size:20;comment:好感度"`
	CreatedAt  time.Time `gorm:"comment:创建时间"`
	UpdatedAt  time.Time `gorm:"comment:更新时间"`
	Stories    []StoryBase `gorm:"many2many:story_students;"` // 与故事的多对多关联
}
