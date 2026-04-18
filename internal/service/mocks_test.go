package service

import (
	"context"
	"time"

	"github.com/codebayu/account-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUUID(uuid string) (*models.User, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	args := m.Called(ctx, userID, tokenID, expiresIn)
	return args.Error(0)
}

func (m *MockTokenRepository) GetRefreshToken(ctx context.Context, userID string, tokenID string) (string, error) {
	args := m.Called(ctx, userID, tokenID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	args := m.Called(ctx, userID, tokenID)
	return args.Error(0)
}

func (m *MockTokenRepository) BlacklistAccessToken(ctx context.Context, tokenID string, expiresIn time.Duration) error {
	args := m.Called(ctx, tokenID, expiresIn)
	return args.Error(0)
}

func (m *MockTokenRepository) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	args := m.Called(ctx, tokenID)
	return args.Bool(0), args.Error(1)
}
