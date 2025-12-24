package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var loginLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLoginLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := loginLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(5.0/60.0, 5)
		loginLimiters[ip] = limiter
	}
	return limiter
}

func LoginRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLoginLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too Many Request login attempsts, please try again letter",
			})
			return
		}
		c.Next()
	}
}
