package dao

import (
	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// StudentDAO 学生数据访问对象
type StudentDAO struct{}

// NewStudentDAO 创建学生数据访问对象
func NewStudentDAO() *StudentDAO {
	return &StudentDAO{}
}

// GetStudentsByQuery 根据查询获取学生信息
func (dao *StudentDAO) GetStudentsByQuery(query string) []models.Student {
	var students []models.Student
	database.GetDB().Where("name LIKE ? OR club_name LIKE ? OR school_name LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&students)
	return students
}
