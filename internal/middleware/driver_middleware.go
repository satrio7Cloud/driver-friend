package middleware

import (
	"github.com/gin-gonic/gin"
)

func OnlyDriver() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rolesVal, exists := ctx.Get("roles")
		if !exists {
			ctx.AbortWithStatusJSON(403, gin.H{
				"error": "roles not found",
			})
			return
		}

		rolesInterface, ok := rolesVal.([]interface{})
		if !ok {
			ctx.AbortWithStatusJSON(403, gin.H{
				"error": "Invalid roles format",
			})
			return
		}

		isDriver := false
		for _, r := range rolesInterface {
			if roleStr, ok := r.(string); ok && roleStr == "driver" {
				isDriver = true
				break
			}
		}

		if !isDriver {
			ctx.AbortWithStatusJSON(403, gin.H{
				"error": "Driver access only",
			})
			return
		}

		driverID, exists := ctx.Get("driver_id")
		if !exists {
			ctx.AbortWithStatusJSON(403, gin.H{
				"error": "driver_id not found",
			})
			return
		}

		ctx.Set("driver_id", driverID)
		ctx.Next()

	}
}
