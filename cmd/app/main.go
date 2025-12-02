package main

import (
	"app/internal/config"
	"app/internal/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 1️⃣ Configuration
	// загружает переменные окружения из файла .env (в корне)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app := server.NewApp()
	app.Run(cfg)
}
