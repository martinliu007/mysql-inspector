package config

import (
	"encoding/json"
	"os"
)

// Config 存储应用程序配置
type Config struct {
	ServerPort int    `json:"server_port"`
	MySQL      MySQL  `json:"mysql"`
}

// MySQL 存储MySQL连接配置
type MySQL struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Load 从配置文件加载配置
func Load() (*Config, error) {
	// 默认配置
	cfg := &Config{
		ServerPort: 8080,
		MySQL: MySQL{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "root",
			Database: "mysql",
		},
	}

	// 尝试从文件加载配置
	if _, err := os.Stat("config.json"); err == nil {
		file, err := os.ReadFile("config.json")
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(file, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}