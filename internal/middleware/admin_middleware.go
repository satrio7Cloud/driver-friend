package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rolesValue, exists := ctx.Get("roles")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: roles not found",
			})
			return
		}

		roles, ok := rolesValue.([]interface{})
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid roles format",
			})
			return
		}

		isAdmin := false
		for _, r := range roles {
			if r == "admin" || r == "super admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Forbiden: only admin can access",
			})
			return
		}
		ctx.Next()
	}
}
