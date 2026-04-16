package handlers

import (
	"github.com/codebayu/account-service/cmd/api/services"
	"gorm.io/gorm"
)

type Handler struct {
	DB          *gorm.DB
	AuthService services.AuthService
}
