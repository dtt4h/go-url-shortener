package middleware

import (
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
	return &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
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

var limiter = newRateLimiter(100, time.Minute)

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
