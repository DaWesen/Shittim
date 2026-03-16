package database

import (
	"log"

	"Shittim/pkg/models"
)

// 自动迁移数据库模型
func AutoMigrate() {
	log.Println("正在自动迁移数据库模型...")

	db := GetDB()

	// 创建向量扩展
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		log.Printf("创建向量扩展失败: %v", err)
	}

	// 确保向量扩展版本支持压缩
	if err := db.Exec("ALTER EXTENSION vector UPDATE").Error; err != nil {
		log.Printf("更新向量扩展失败: %v", err)
	}

	// 按照依赖顺序创建表
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	err = db.AutoMigrate(&models.Signin{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	err = db.AutoMigrate(&models.School{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	err = db.AutoMigrate(&models.Club{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	err = db.AutoMigrate(
		&models.StoryBase{},
		&models.DailyStory{},
		&models.EventStory{},
		&models.ExclusiveMemory{},
		&models.Conversation{},
		&models.Message{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建向量索引
	log.Println("正在创建向量索引...")

	// 为 Conversation 表创建向量索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_conversation_message_vector ON conversation USING ivfflat (message_vector vector_cosine_ops) WITH (lists = 100);").Error; err != nil {
		log.Printf("创建 conversation.message_vector 索引失败: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_conversation_response_vector ON conversation USING ivfflat (response_vector vector_cosine_ops) WITH (lists = 100);").Error; err != nil {
		log.Printf("创建 conversation.response_vector 索引失败: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_conversation_topic_vector ON conversation USING ivfflat (topic_vector vector_cosine_ops) WITH (lists = 100);").Error; err != nil {
		log.Printf("创建 conversation.topic_vector 索引失败: %v", err)
	}

	// 为 ExclusiveMemory 表创建向量索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_exclusivememory_content_vector ON exclusive_memory USING ivfflat (content_vector vector_cosine_ops) WITH (lists = 100);").Error; err != nil {
		log.Printf("创建 exclusive_memory.content_vector 索引失败: %v", err)
	}

	// 为 DailyStory 表创建向量索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_dailystory_content_vector ON daily_story USING ivfflat (content_vector vector_cosine_ops) WITH (lists = 100);").Error; err != nil {
		log.Printf("创建 daily_story.content_vector 索引失败: %v", err)
	}

	// 为 EventStory 表创建向量索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_eventstory_content_vector ON event_story USING ivfflat (content_vector vector_cosine_ops) WITH (lists = 100);").Error; err != nil {
		log.Printf("创建 event_story.content_vector 索引失败: %v", err)
	}

	log.Println("数据库迁移成功，向量索引创建完成")
}
