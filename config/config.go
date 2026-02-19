package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Arona    AronaConfig    `mapstructure:"arona"`
	Momotalk MomotalkConfig `mapstructure:"momotalk"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type AronaConfig struct {
	NickName      []string        `mapstructure:"nickname"`
	CommandPrefix string          `mapstructure:"command_prefix"`
	SuperUsers    []int64         `mapstructure:"super_users"`
	WebSocket     WebSocketConfig `mapstructure:"websocket"`
}

type WebSocketConfig struct {
	Port int    `mapstructure:"port"`
	URL  string `mapstructure:"url"`
}

type MomotalkConfig struct {
	AI AIConfig `mapstructure:"ai"`
}

type AIConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

var AppConfig Config

func LoadConfig() error {
	// 创建viper实例
	v := viper.New()

	// 配置文件名称
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 尝试从多个可能的路径加载配置文件
	possiblePaths := []string{
		"./config",             // 从项目根目录运行
		"../config",            // 从Arona目录运行
		"../../config",         // 从Arona/Arona目录运行
		filepath.Dir("./"),     // 当前目录
		filepath.Dir("../"),    // 上级目录
		filepath.Dir("../../"), // 上上级目录
	}

	// 添加所有可能的配置路径
	for _, path := range possiblePaths {
		v.AddConfigPath(path)
	}

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置到结构体
	err = v.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

func InitConfig() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
}

func GetPostgresDSN() string {
	pg := AppConfig.Database.Postgres
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode)
}
