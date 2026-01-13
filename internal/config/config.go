package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	EnableAdminSeed bool

	AdminEmail    string
	AdminPassword string
	AdminPhone    string
	AdminName     string

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
		AppEnv:          os.Getenv("APP_ENV"),
		EnableAdminSeed: os.Getenv("ENABLE_ADMIN_SEED") == "true",

		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		AdminPhone:    os.Getenv("ADMIN_PHONE"),
		AdminName:     os.Getenv("ADMIN_NAME"),

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
