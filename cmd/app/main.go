package main

import (
	"app/internal/app"
	"app/internal/config"
	"fmt"
	"log"
)

func main() {
	// 1️⃣ Configuration
	// // обрабатываю аргументы командной строки и сохраняю их значения в соответствующих переменных
	// config.ParseFlags()

	cfg, err := config.NewConfig()
	fmt.Println("cfg in main-", *cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
