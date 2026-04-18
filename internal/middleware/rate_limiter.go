package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rate     rate.Limit
	burst    int
}

func newRateLimiterStore(r rate.Limit, burst int) *rateLimiterStore {
	store := &rateLimiterStore{
		limiters: make(map[string]*ipLimiter),
		rate:     r,
		burst:    burst,
	}
	// Cleanup expired limiters setiap 1 menit
	go store.cleanup()
	return store
}

func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, exists := s.limiters[ip]; exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(s.rate, s.burst)
	s.limiters[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func (s *rateLimiterStore) cleanup() {
	for {
		time.Sleep(1 * time.Minute)
		s.mu.Lock()
		for ip, entry := range s.limiters {
			if time.Since(entry.lastSeen) > 3*time.Minute {
				delete(s.limiters, ip)
			}
		}
		s.mu.Unlock()
	}
}

func newRateLimiterMiddleware(store *rateLimiterStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()
			limiter := store.getLimiter(ip)
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "too many requests, please slow down",
				})
			}
			return next(c)
		}
	}
}

// GlobalRateLimiter — 20 req/s, burst 50
func GlobalRateLimiter() echo.MiddlewareFunc {
	store := newRateLimiterStore(20, 50)
	return newRateLimiterMiddleware(store)
}

// StrictRateLimiter — 5 req/s, burst 10 (untuk auth routes)
func StrictRateLimiter() echo.MiddlewareFunc {
	store := newRateLimiterStore(5, 10)
	return newRateLimiterMiddleware(store)
}
