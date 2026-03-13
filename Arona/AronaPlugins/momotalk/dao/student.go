package dao

import (
	"fmt"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
	"Shittim/pkg/vector"

	"github.com/pgvector/pgvector-go"
)

// StudentDAO 学生数据访问对象
type StudentDAO struct{}

// NewStudentDAO 创建学生数据访问对象
func NewStudentDAO() *StudentDAO {
	return &StudentDAO{}
}

// GetStudentsByQuery 根据查询获取学生信息（使用向量搜索）
func (dao *StudentDAO) GetStudentsByQuery(query string) []models.Student {
	// 生成查询向量
	queryVector, err := vector.GenerateEmbedding(query)
	if err == nil && len(queryVector.Slice()) > 0 {
		// 使用向量搜索
		return dao.GetStudentsByVector(queryVector, 10)
	}

	// 回退到传统的LIKE查询
	var students []models.Student
	database.GetDB().Where("name LIKE ? OR club_name LIKE ? OR school_name LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&students)
	return students
}

// GetStudentsByVector 根据向量相似度获取学生信息
func (dao *StudentDAO) GetStudentsByVector(queryVector pgvector.Vector, limit int) []models.Student {
	var students []models.Student
	result := database.GetDB().Model(&models.Student{}).
		Select("*, 1 - (embedding <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&students)

	if result.Error != nil {
		// 处理错误，回退到传统查询
		return []models.Student{}
	}
	return students
}

// UpdateStudentEmbedding 更新学生信息的向量嵌入
func (dao *StudentDAO) UpdateStudentEmbedding(studentName string) error {
	var student models.Student
	if err := database.GetDB().Where("name = ?", studentName).First(&student).Error; err != nil {
		return err
	}

	// 构建学生信息文本
	studentInfo := fmt.Sprintf("%s %s %s %s %s",
		student.Name, student.SchoolName, student.ClubName, student.Love, student.Level)

	// 生成向量嵌入
	embedding, err := vector.GenerateEmbedding(studentInfo)
	if err != nil {
		return err
	}

	// 更新向量
	return database.GetDB().Model(&student).Update("embedding", embedding).Error
}

// UpdateAllStudentEmbeddings 更新所有学生的向量嵌入
func (dao *StudentDAO) UpdateAllStudentEmbeddings() error {
	var students []models.Student
	if err := database.GetDB().Find(&students).Error; err != nil {
		return err
	}

	for _, student := range students {
		if err := dao.UpdateStudentEmbedding(student.Name); err != nil {
			// 记录错误但继续处理其他学生
			fmt.Printf("更新学生 %s 的向量嵌入失败: %v\n", student.Name, err)
		}
	}

	return nil
}
