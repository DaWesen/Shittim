package vector

import (
	"Shittim/pkg/database"
	"Shittim/pkg/models"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// SearchService 向量搜索服务
type SearchService struct {
	db *gorm.DB
}

// NewSearchService 创建搜索服务
func NewSearchService() *SearchService {
	return &SearchService{
		db: database.GetDB(),
	}
}

// SearchSimilarConversations 搜索相似的对话
func (s *SearchService) SearchSimilarConversations(userID int64, queryVector pgvector.Vector, limit int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	result := s.db.Model(&models.Conversation{}).
		Where("user_id = ?", userID).
		Select("*, 1 - (topic_vector <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&conversations)
	return conversations, result.Error
}

// SearchSimilarMessages 搜索相似的消息
func (s *SearchService) SearchSimilarMessages(conversationID uint, queryVector pgvector.Vector, limit int) ([]models.Message, error) {
	var messages []models.Message
	result := s.db.Model(&models.Message{}).
		Where("conversation_id = ?", conversationID).
		Select("*, 1 - (content_vector <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&messages)
	return messages, result.Error
}

// SearchSimilarMemories 搜索相似的记忆
func (s *SearchService) SearchSimilarMemories(studentName string, queryVector pgvector.Vector, limit int) ([]models.ExclusiveMemory, error) {
	var memories []models.ExclusiveMemory
	result := s.db.Model(&models.ExclusiveMemory{}).
		Where("student_name = ?", studentName).
		Select("*, 1 - (content_vector <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&memories)
	return memories, result.Error
}

// SearchSimilarStudents 搜索相似的学生
func (s *SearchService) SearchSimilarStudents(queryVector pgvector.Vector, limit int) ([]models.Student, error) {
	var students []models.Student
	result := s.db.Model(&models.Student{}).
		Select("*, 1 - (embedding <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&students)
	return students, result.Error
}

// SearchSimilarDailyStories 搜索相似的日常故事
func (s *SearchService) SearchSimilarDailyStories(queryVector pgvector.Vector, limit int) ([]models.DailyStory, error) {
	var stories []models.DailyStory
	result := s.db.Model(&models.DailyStory{}).
		Select("*, 1 - (content_vector <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&stories)
	return stories, result.Error
}

// SearchSimilarEventStories 搜索相似的活动故事
func (s *SearchService) SearchSimilarEventStories(queryVector pgvector.Vector, limit int) ([]models.EventStory, error) {
	var stories []models.EventStory
	result := s.db.Model(&models.EventStory{}).
		Select("*, 1 - (content_vector <=> ?) as similarity", queryVector).
		Order("similarity DESC").
		Limit(limit).
		Find(&stories)
	return stories, result.Error
}
