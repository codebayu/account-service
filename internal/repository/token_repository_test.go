package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenRepositoryTestSuite struct {
	suite.Suite
	mr   *miniredis.Miniredis
	repo TokenRepository
}

func (s *TokenRepositoryTestSuite) SetupTest() {
	mr, err := miniredis.Run()
	s.Require().NoError(err)
	s.mr = mr

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	s.repo = NewTokenRepository(rdb)
}

func (s *TokenRepositoryTestSuite) TearDownTest() {
	s.mr.Close()
}

func (s *TokenRepositoryTestSuite) TestSaveRefreshToken() {
	ctx := context.Background()
	err := s.repo.SaveRefreshToken(ctx, "user-1", "jti-1", 5*time.Minute)
	s.NoError(err)

	// Verify key exists
	val, err := s.repo.GetRefreshToken(ctx, "user-1", "jti-1")
	s.NoError(err)
	s.Equal("valid", val)
}

func (s *TokenRepositoryTestSuite) TestGetRefreshToken_NotFound() {
	ctx := context.Background()
	_, err := s.repo.GetRefreshToken(ctx, "user-missing", "jti-missing")
	s.Error(err)
}

func (s *TokenRepositoryTestSuite) TestDeleteRefreshToken() {
	ctx := context.Background()
	_ = s.repo.SaveRefreshToken(ctx, "user-2", "jti-2", 5*time.Minute)

	err := s.repo.DeleteRefreshToken(ctx, "user-2", "jti-2")
	s.NoError(err)

	// After delete, should not be found
	_, err = s.repo.GetRefreshToken(ctx, "user-2", "jti-2")
	s.Error(err)
}

func (s *TokenRepositoryTestSuite) TestBlacklistAccessToken() {
	ctx := context.Background()
	err := s.repo.BlacklistAccessToken(ctx, "access-jti-1", 5*time.Minute)
	s.NoError(err)
}

func (s *TokenRepositoryTestSuite) TestIsAccessTokenBlacklisted_True() {
	ctx := context.Background()
	_ = s.repo.BlacklistAccessToken(ctx, "access-jti-2", 5*time.Minute)

	blacklisted, err := s.repo.IsAccessTokenBlacklisted(ctx, "access-jti-2")
	s.NoError(err)
	s.True(blacklisted)
}

func (s *TokenRepositoryTestSuite) TestIsAccessTokenBlacklisted_False() {
	ctx := context.Background()

	blacklisted, err := s.repo.IsAccessTokenBlacklisted(ctx, "nonexistent-jti")
	s.NoError(err)
	s.False(blacklisted)
}

func (s *TokenRepositoryTestSuite) TestRefreshTokenExpiry() {
	ctx := context.Background()
	_ = s.repo.SaveRefreshToken(ctx, "user-3", "jti-3", 1*time.Second)

	// Fast-forward time in miniredis
	s.mr.FastForward(2 * time.Second)

	_, err := s.repo.GetRefreshToken(ctx, "user-3", "jti-3")
	s.Error(err, "token should be expired")
}

func TestTokenRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TokenRepositoryTestSuite))
}

func TestNewTokenRepository(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := NewTokenRepository(rdb)
	assert.NotNil(t, repo)
}
