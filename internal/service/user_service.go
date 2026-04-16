package service

import (
	"errors"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/internal/models"
	"github.com/codebayu/account-service/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(uuid string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(uuid string) (*models.User, error) {
	user, err := s.repo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.ErrInternalServer
	}
	return user, nil
}
