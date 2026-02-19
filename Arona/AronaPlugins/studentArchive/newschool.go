package studentArchive

import (
	"fmt"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 创建新学院
func CreateSchool(schoolName string) (*models.School, error) {
	// 检查学院是否已存在
	var existingSchool models.School
	if err := database.GetDB().First(&existingSchool, "school_name = ?", schoolName).Error; err == nil {
		return nil, fmt.Errorf("学院已存在：%s", schoolName)
	}

	// 创建学院记录
	school := models.School{
		SchoolName:    schoolName,
		StudentCounts: 0,
	}

	// 保存到数据库
	if err := database.GetDB().Create(&school).Error; err != nil {
		return nil, fmt.Errorf("创建学院失败：%v", err)
	}

	return &school, nil
}
