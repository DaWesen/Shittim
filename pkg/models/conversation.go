package models

import "time"

// Conversation 对话存储模型，存储用户与bot的对话
type Conversation struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null;index;comment:用户ID"`
	Messages  []Message `gorm:"foreignKey:ConversationID;comment:消息列表"`
	StartAt   time.Time `gorm:"autoCreateTime;comment:对话开始时间"`
	EndAt     time.Time `gorm:"comment:对话结束时间"`
	LastActive time.Time `gorm:"autoUpdateTime;comment:最后活跃时间"`
	Topic     string    `gorm:"size:100;comment:对话主题"`
	Status    string    `gorm:"size:20;default:active;comment:对话状态"` // active, ended
}

// Message 消息模型，存储单条对话消息
type Message struct {
	ID             uint      `gorm:"primaryKey"`
	ConversationID uint      `gorm:"index;comment:所属对话ID"`
	Content        string    `gorm:"type:text;not null;comment:消息内容"`
	IsUser         bool      `gorm:"default:true;comment:是否用户消息"`
	Timestamp      time.Time `gorm:"autoCreateTime;comment:消息时间戳"`
	MessageType    string    `gorm:"size:20;default:text;comment:消息类型"` // text, image, voice
	Emotion        string    `gorm:"size:20;comment:情绪标签"`
}
