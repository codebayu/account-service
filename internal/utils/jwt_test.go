package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	uuid := "test-uuid"
	token, exp, err := GenerateToken(uuid, 1*time.Hour, false)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, exp.After(time.Now()))
}

func TestParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	uuid := "test-uuid"
	token, _, _ := GenerateToken(uuid, 1*time.Hour, false)

	claims, err := ParseToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, uuid, claims.UUID)
	assert.False(t, claims.IsRefresh)
}

func TestParseToken_Invalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	_, err := ParseToken("invalid.token.string")
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	uuid := "test-uuid"
	// Generate token with negative duration
	token, _, _ := GenerateToken(uuid, -1*time.Hour, false)

	_, err := ParseToken(token)
	assert.Error(t, err)
}

func TestGenerateToken_EmptySecret(t *testing.T) {
	// Temporarily clear JWT_SECRET to test fallback
	oldSecret := os.Getenv("JWT_SECRET")
	os.Unsetenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", oldSecret)

	uuid := "test-uuid"
	token, _, err := GenerateToken(uuid, 1*time.Hour, false)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
