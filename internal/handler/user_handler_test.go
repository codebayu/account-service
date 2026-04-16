package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/models"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetCurrentUser(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandler(mockSvc)
	e := echo.New()

	uuid := "user-uuid"
	user := &models.User{UUID: uuid, Name: "Test User", Email: "test@example.com"}

	t.Run("Success Get Current User", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/current", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_uuid", uuid)

		mockSvc.On("GetProfile", uuid).Return(user, nil).Once()

		err := h.GetCurrentUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp response.Response
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, 200000, resp.Result.Code)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Unauthorized - Missing Context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/current", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetCurrentUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
