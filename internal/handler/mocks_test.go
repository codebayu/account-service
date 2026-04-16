package handler

import (
	"github.com/codebayu/account-service/internal/dto"
	"github.com/codebayu/account-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req dto.RegisterRequest) (*dto.AuthResponseData, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponseData), args.Error(1)
}

func (m *MockAuthService) Login(req dto.LoginRequest) (*dto.AuthResponseData, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponseData), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetProfile(uuid string) (*models.User, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
