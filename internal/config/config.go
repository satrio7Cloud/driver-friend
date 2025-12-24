package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBPort     string
	DBUser     string
	DBHost     string
	DBPassword string
	DBName     string
	ServerPort string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("APP_PORT"),
		Port:       os.Getenv("PORT"),
	}
	log.Println("Configuration loaded:", cfg)
	return *cfg
}
