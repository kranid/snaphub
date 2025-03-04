package config

import (
	"os"
)

type Config struct {
	ApiKey    string
	AccessKey string
	DBName    string
	DBUser    string
	DBPass    string
}

// InitConfig загружает конфигурацию из переменных окружения
func InitConfig() Config {
	return Config{
		ApiKey:    os.Getenv("API_ACCESS_KEY"),      // Загружает API_ACCESS_KEY из переменной окружения
		AccessKey: os.Getenv("API_SECRET_KEY"),      // Загружает API_SECRET_KEY
		DBName:    os.Getenv("MYSQL_DATABASE"),      // Загружает MYSQL_DATABASE
		DBUser:    os.Getenv("MYSQL_USER"),          // Загружает MYSQL_USER
		DBPass:    os.Getenv("MYSQL_ROOT_PASSWORD"), // Загружает MYSQL_ROOT_PASSWORD
	}
}
