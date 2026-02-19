package studentArchive

import (
	"fmt"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 更新学生档案信息（名字不可修改）
func UpdateStudentArchive(name, level string, age uint, clubName, schoolName, height, love, underEye, affection string) (*models.Student, error) {
	// 查找学生
	var student models.Student
	if err := database.GetDB().First(&student, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("没有找到这名学生：%v", err)
	}

	// 更新学生信息（名字保持不变）
	student.Level = level
	student.Age = age
	student.ClubName = clubName
	student.SchoolName = schoolName
	student.Height = height
	student.Love = love
	student.UnderEye = underEye
	student.Affection = affection

	// 保存更新
	if err := database.GetDB().Save(&student).Error; err != nil {
		return nil, fmt.Errorf("更新学生档案失败：%v", err)
	}

	return &student, nil
}
