package database

import (
	"log"

	"Shittim/pkg/models"
)

// 自动迁移数据库模型
func AutoMigrate() {
	log.Println("正在自动迁移数据库模型...")

	db := GetDB()
	err := db.AutoMigrate(
		&models.User{},
		&models.Signin{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移成功")
}
