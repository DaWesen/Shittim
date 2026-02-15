package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Arona    AronaConfig    `yaml:"arona"`
	Momotalk MomotalkConfig `yaml:"momotalk"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type AronaConfig struct {
	NickName      []string        `yaml:"nickname"`
	CommandPrefix string          `yaml:"command_prefix"`
	SuperUsers    []int64         `yaml:"super_users"`
	WebSocket     WebSocketConfig `yaml:"websocket"`
}

type WebSocketConfig struct {
	Port int    `yaml:"port"`
	URL  string `yaml:"url"`
}

type MomotalkConfig struct {
	AI AIConfig `yaml:"ai"`
}

type AIConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

var AppConfig Config

func LoadConfig() error {
	//尝试从多个可能的路径加载配置文件
	possiblePaths := []string{
		"config/config.yaml",       // 从项目根目录运行
		"../config/config.yaml",    // 从Arona目录运行
		"../../config/config.yaml", // 从Arona/Arona目录运行
	}

	var file []byte
	var err error

	//尝试所有可能的路径
	for _, path := range possiblePaths {
		file, err = os.ReadFile(path)
		if err == nil {
			//找到配置文件
			break
		}
	}

	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	err = yaml.Unmarshal(file, &AppConfig)
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
