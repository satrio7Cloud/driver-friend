package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	otpLimiters = make(map[string]*rate.Limiter)
	otpMu       sync.Mutex
)

func getOTPLimiter(key string) *rate.Limiter {
	otpMu.Lock()
	defer otpMu.Unlock()

	limiter, exists := otpLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(1.0/60.0, 1)
		otpLimiters[key] = limiter
	}
	return limiter
}

func OTPRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		key := userID.(string)
		limiter := getOTPLimiter(key)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "OTP request to fequent, please wait",
			})
			return
		}
		c.Next()
	}
}
