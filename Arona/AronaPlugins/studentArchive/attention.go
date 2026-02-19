package studentArchive

import (
	"fmt"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 更新学生关注状态
func UpdateAttentionStatus(name, status string) (*models.Student, error) {
	var student models.Student
	if err := database.GetDB().First(&student, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("没有找到这名学生：%v", err)
	}

	// 更新关注状态
	student.UnderEye = status
	if err := database.GetDB().Save(&student).Error; err != nil {
		return nil, fmt.Errorf("更新关注状态失败：%v", err)
	}

	return &student, nil
}

// 获取学生信息
func GetStudent(name string) (*models.Student, error) {
	var student models.Student
	if err := database.GetDB().First(&student, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("没有找到这名学生：%v", err)
	}

	return &student, nil
}
