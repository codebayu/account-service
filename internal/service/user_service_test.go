package service

import (
	"testing"

	"github.com/codebayu/account-service/common/apperror"
	"github.com/codebayu/account-service/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserService_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewUserService(mockRepo)

	uuid := "user-uuid"
	user := &models.User{
		UUID:  uuid,
		Name:  "Test User",
		Email: "test@example.com",
	}

	t.Run("Success Get Profile", func(t *testing.T) {
		mockRepo.On("FindByUUID", uuid).Return(user, nil).Once()

		res, err := svc.GetProfile(uuid)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uuid, res.UUID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindByUUID", "notfound").Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := svc.GetProfile("notfound")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, apperror.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
