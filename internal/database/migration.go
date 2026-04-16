package database

import (
	"log"

	"github.com/codebayu/account-service/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("🔄 running migration...")
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("❌ migration failed: %v\n", err)
		return err
	}
	log.Println("✅ migration success")
	return nil
}
