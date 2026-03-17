package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Запускаем cleanup каждые 5 минут
	go rl.cleanup(5 * time.Minute)

	return rl
}

func (r *rateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-r.window)

		for key, times := range r.requests {
			var valid []time.Time
			for _, t := range times {
				if t.After(windowStart) {
					valid = append(valid, t)
				}
			}

			if len(valid) == 0 {
				delete(r.requests, key)
			} else {
				r.requests[key] = valid
			}
		}
		r.mu.Unlock()
	}
}

func (r *rateLimiter) allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	requests := r.requests[key]
	var valid []time.Time

	for _, t := range requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= r.limit {
		r.requests[key] = valid
		return false
	}

	r.requests[key] = append(valid, now)
	return true
}

// Конфигурация rate limiter через переменные окружения
var (
	rateLimitLimit  = 100
	rateLimitWindow = time.Minute
)

func init() {
	if v := getEnv("RATE_LIMIT_REQUESTS", "100"); v != "" {
		if parsed, err := parseInt(v); err == nil && parsed > 0 {
			rateLimitLimit = parsed
		}
	}
	if v := getEnv("RATE_LIMIT_WINDOW_SECONDS", "60"); v != "" {
		if parsed, err := parseInt(v); err == nil && parsed > 0 {
			rateLimitWindow = time.Duration(parsed) * time.Second
		}
	}
}

var limiter = newRateLimiter(rateLimitLimit, rateLimitWindow)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if !limiter.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}

func getEnv(key, defaultValue string) string {
	if v := getEnvFunc(key); v != "" {
		return v
	}
	return defaultValue
}

var getEnvFunc = func(key string) string {
	// Переопределяется в тестах
	return ""
}

var parseInt = func(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
