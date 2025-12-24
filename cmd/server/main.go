package main

import (
	"be/internal/app"
	"be/internal/config"
	"be/internal/db"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db.ConnectDatabase(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	server := app.NewApp(cfg.ServerPort)

	log.Println("Server is running on port", cfg.ServerPort)
	server.Run()
}
