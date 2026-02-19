package studentArchive

import (
	"fmt"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 关注状态常量
const (
	Default   = "默认"
	Attention = "重点关注"
)

// 创建新学生
func CreateStudent(name, level string, age uint, clubName, schoolName, height, love, affection string) (*models.Student, error) {
	// 检查学校是否存在
	var school models.School
	if err := database.GetDB().First(&school, "school_name = ?", schoolName).Error; err != nil {
		return nil, fmt.Errorf("学校不存在: %v", err)
	}

	// 检查社团是否存在
	var club models.Club
	if err := database.GetDB().First(&club, "club_name = ?", clubName).Error; err != nil {
		return nil, fmt.Errorf("社团不存在: %v", err)
	}

	// 创建学生记录
	student := models.Student{
		Name:       name,
		Level:      level,
		Age:        age,
		ClubName:   clubName,
		SchoolName: schoolName,
		Height:     height,
		Love:       love,
		Affection:  affection,
		UnderEye:   Default,
	}

	// 保存到数据库
	if err := database.GetDB().Create(&student).Error; err != nil {
		return nil, fmt.Errorf("创建学生失败: %v", err)
	}

	// 更新学校和社团的学生数量
	if err := database.GetDB().Model(&school).Update("student_counts", school.StudentCounts+1).Error; err != nil {
		return nil, fmt.Errorf("更新学校学生数量失败: %v", err)
	}

	if err := database.GetDB().Model(&club).Update("student_counts", club.StudentCounts+1).Error; err != nil {
		return nil, fmt.Errorf("更新社团学生数量失败: %v", err)
	}

	return &student, nil
}

// 获取所有学校
func GetAllSchools() ([]models.School, error) {
	var schools []models.School
	if err := database.GetDB().Find(&schools).Error; err != nil {
		return nil, fmt.Errorf("获取学校列表失败: %v", err)
	}
	return schools, nil
}

// 获取所有社团
func GetAllClubs() ([]models.Club, error) {
	var clubs []models.Club
	if err := database.GetDB().Find(&clubs).Error; err != nil {
		return nil, fmt.Errorf("获取社团列表失败: %v", err)
	}
	return clubs, nil
}
