package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyDriver() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: role not found",
			})
			ctx.Abort()
			return
		}

		if role != "driver" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Forbiden: only driver can access this routes",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
