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

func getOTPLimiterByPhone(key string) *rate.Limiter {
	otpMu.Lock()
	defer otpMu.Unlock()

	limiter, exists := otpLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(1.0/60.0, 1) // 1 request pe 60sec
		otpLimiters[key] = limiter
	}
	return limiter
}

func OTPRateLimiterByPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Phone string `json:"phone"`
		}

		if err := c.ShouldBindJSON(&body); err != nil || body.Phone == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Phone is required",
			})
			return
		}

		limiter := getOTPLimiterByPhone(body.Phone)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "OTP request to frequent, please wait",
			})
			return
		}

		c.Set("phone", body.Phone)
		c.Next()
	}
}
