package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: role not found",
			})
			ctx.Abort()
			return
		}

		if role != "admin" && role != "super admin" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Forbiden: only admin can access this routes",
			})
			ctx.Abort()
			return
		}
		ctx.Next()

	}
}
