package main

import (
	"log"

	"github.com/codebayu/account-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
