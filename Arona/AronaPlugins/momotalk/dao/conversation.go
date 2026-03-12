package dao

import (
	"Shittim/pkg/database"
	"Shittim/pkg/models"
	"time"
)

// ConversationDAO 对话数据访问对象
type ConversationDAO struct{}

// NewConversationDAO 创建对话数据访问对象
func NewConversationDAO() *ConversationDAO {
	return &ConversationDAO{}
}

// SaveConversation 保存对话
func (dao *ConversationDAO) SaveConversation(userID int64, userMessage, botResponse string) error {
	// 创建对话
	conversation := models.Conversation{
		UserID:     userID,
		StartAt:    time.Now(),
		LastActive: time.Now(),
		Topic:      "日常对话",
		Status:     "active",
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 创建对话
	if err := tx.Create(&conversation).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建用户消息
	userMessageModel := models.Message{
		ConversationID: conversation.ID,
		Content:        userMessage,
		IsUser:         true,
		Timestamp:      time.Now(),
		MessageType:    "text",
	}
	if err := tx.Create(&userMessageModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建机器人消息
	botMessageModel := models.Message{
		ConversationID: conversation.ID,
		Content:        botResponse,
		IsUser:         false,
		Timestamp:      time.Now(),
		MessageType:    "text",
	}
	if err := tx.Create(&botMessageModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
