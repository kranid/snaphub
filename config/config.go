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
  ApiKey:    os.Getenv("API_ACCESS_KEY"),        // Загружает API_ACCESS_KEY из переменной окружения
  AccessKey: os.Getenv("API_SECRET_KEY"),        // Загружает API_SECRET_KEY
  DBName:    os.Getenv("DB_NAME"),                // Загружает DB_NAME
  DBUser:    os.Getenv("DB_USER"),                // Загружает DB_USER
  DBPass:    os.Getenv("DB_PASS"),                // Загружает DB_PASS
 }
}