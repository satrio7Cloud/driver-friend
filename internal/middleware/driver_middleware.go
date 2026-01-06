package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyDriver() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rolesAny, exists := ctx.Get("roles")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: role not found",
			})
			ctx.Abort()
			return
		}

		roles, ok := rolesAny.([]string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid roles format",
			})
			ctx.Abort()
			return
		}

		isDriver := false
		for _, r := range roles {
			if r == "driver" {
				isDriver = true
				break
			}
		}

		if !isDriver {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Forbiden: driver access only",
			})
			return
		}

		driverID, exists := ctx.Get("driver_id")
		if !exists || driverID == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "driver_id is missing in token",
			})
			ctx.Abort()
			return
		}

		ctx.Next()

	}
}
