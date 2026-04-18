package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	GetRefreshToken(ctx context.Context, userID string, tokenID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error

	// Access Token Blocklist
	BlacklistAccessToken(ctx context.Context, tokenID string, expiresIn time.Duration) error
	IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

type tokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) TokenRepository {
	return &tokenRepository{rdb: rdb}
}

func (r *tokenRepository) getKey(userID, tokenID string) string {
	return fmt.Sprintf("auth:refresh_token:%s:%s", userID, tokenID)
}

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := r.getKey(userID, tokenID)
	return r.rdb.Set(ctx, key, "valid", expiresIn).Err()
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, userID string, tokenID string) (string, error) {
	key := r.getKey(userID, tokenID)
	return r.rdb.Get(ctx, key).Result()
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("auth:refresh_token:%s:%s", userID, tokenID)
	return r.rdb.Del(ctx, key).Err()
}

func (r *tokenRepository) BlacklistAccessToken(ctx context.Context, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("auth:blacklist:%s", tokenID)
	return r.rdb.Set(ctx, key, "revoked", expiresIn).Err()
}

func (r *tokenRepository) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("auth:blacklist:%s", tokenID)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "revoked", nil
}
