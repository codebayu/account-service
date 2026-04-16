package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}
