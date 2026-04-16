package main

import (
	"log"

	"github.com/codebayu/account-service/internal/config"
	"github.com/codebayu/account-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
