package database

import (
	"fmt"
	"log"

	"Shittim/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 初始化数据库连接
func InitDatabase() {
	dsn := config.GetPostgresDSN()
	fmt.Printf("正在连接数据库: %s\n", dsn)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功")
}

// 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
