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

	log.Println("数据库迁移成功")
}
