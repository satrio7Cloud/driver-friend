package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	ipLimiters = make(map[string]*rate.Limiter)
	ipMu       sync.Mutex
)

func getIpLimiter(ip string) *rate.Limiter {
	ipMu.Lock()
	defer ipMu.Unlock()

	limiter, exists := ipLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(3.0/300.0, 3)
		ipLimiters[ip] = limiter
	}
	return limiter
}

func IPRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getIpLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too Many Request From This IP",
			})
			return
		}
		c.Next()
	}
}
