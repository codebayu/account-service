package services

import (
	"errors"

	"github.com/codebayu/account-service/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) GetProfile(uuid string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
