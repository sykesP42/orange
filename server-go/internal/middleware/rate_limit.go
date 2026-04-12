package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

type visitor struct {
	tokens    int
	lastSeen  time.Time
}

const (
	maxTokens     = 100
	refillRate  = 10
	cleanupInterval = time.Minute * 5
)

var limiter = &RateLimiter{
	visitors: make(map[string]*visitor),
}

func init() {
	go limiter.cleanup()
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(cleanupInterval)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > time.Minute*10 {
			delete(rl.visitors, ip)
		}
	}
	rl.mu.Unlock()
}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:   maxTokens - 1,
			lastSeen: now,
		}
		return true
	}

	elapsed := now.Sub(v.lastSeen)
	tokensToAdd := int(elapsed.Seconds()) * refillRate
	if tokensToAdd > 0 {
		v.tokens = min(v.tokens+tokensToAdd, maxTokens)
	}

	v.lastSeen = now

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
