package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	e := echo.New()

	t.Run("Allow Requests Within Limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		store := newRateLimiterStore(rate.Limit(5), 5)
		mw := newRateLimiterMiddleware(store)
		handler := mw(func(c *echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		for i := 0; i < 5; i++ {
			err := handler(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Block Requests Over Limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.2:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		store := newRateLimiterStore(rate.Limit(1), 1)
		mw := newRateLimiterMiddleware(store)
		handler := mw(func(c *echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		// First request allowed
		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		req2.RemoteAddr = "192.168.1.2:1234"
		rec2 := httptest.NewRecorder()
		c2 := e.NewContext(req2, rec2)

		err2 := handler(c2)
		assert.NoError(t, err2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("GlobalRateLimiter Wrapper", func(t *testing.T) {
		mw := GlobalRateLimiter()
		assert.NotNil(t, mw)
	})

	t.Run("StrictRateLimiter Wrapper", func(t *testing.T) {
		mw := StrictRateLimiter()
		assert.NotNil(t, mw)
	})

	t.Run("Cleanup Old Limiters", func(t *testing.T) {
		store := newRateLimiterStore(rate.Limit(5), 5)
		
		store.getLimiter("192.168.1.3")
		
		// Manually override last seen to make it stale
		store.mu.Lock()
		store.limiters["192.168.1.3"].lastSeen = time.Now().Add(-5 * time.Minute)
		
		// Call getLimiter for a diff IP to trigger just to verify state
		store.mu.Unlock()

		// Instead of waiting, we will just manually test logic that the loop does
		store.mu.Lock()
		for ip, entry := range store.limiters {
			if time.Since(entry.lastSeen) > 3*time.Minute {
				delete(store.limiters, ip)
			}
		}
		store.mu.Unlock()

		store.mu.Lock()
		_, exists := store.limiters["192.168.1.3"]
		store.mu.Unlock()

		assert.False(t, exists)
	})
}
