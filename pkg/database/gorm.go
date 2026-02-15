package database

import (
	"fmt"
	"log"

	"Shittim/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() error {
	dsn := config.GetPostgresDSN()
	fmt.Printf("正在连接数据库: %s\n", dsn)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}
