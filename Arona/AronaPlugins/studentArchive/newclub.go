package studentArchive

import (
	"fmt"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 创建新社团（需要挂靠到学院）
func CreateClub(clubName, schoolName string) (*models.Club, error) {
	// 检查学院是否存在
	var school models.School
	if err := database.GetDB().First(&school, "school_name = ?", schoolName).Error; err != nil {
		return nil, fmt.Errorf("学院不存在：%v", err)
	}

	// 检查社团是否已存在
	var existingClub models.Club
	if err := database.GetDB().First(&existingClub, "club_name = ?", clubName).Error; err == nil {
		return nil, fmt.Errorf("社团已存在：%s", clubName)
	}

	// 创建社团记录
	club := models.Club{
		ClubName:      clubName,
		SchoolName:    schoolName,
		StudentCounts: 0,
	}

	// 保存到数据库
	if err := database.GetDB().Create(&club).Error; err != nil {
		return nil, fmt.Errorf("创建社团失败：%v", err)
	}

	return &club, nil
}
