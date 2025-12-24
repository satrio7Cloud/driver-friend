package routes

import (
	"be/internal/middleware"
	"be/internal/modules/auth/controller"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authController *controller.AuthController
}

func NewAuthRoutes(authController *controller.AuthController) *AuthRoutes {
	return &AuthRoutes{
		authController: authController,
	}
}

func (r *AuthRoutes) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		// ========================
		// PUBLIC
		// ========================
		auth.POST("/register", r.authController.Register)
		auth.POST(
			"/login",
			middleware.LoginRateLimiter(),
			r.authController.Login,
		)
		auth.POST("/login-phone", middleware.LoginRateLimiter(), r.authController.LoginPhone)

		auth.POST("/verify/email", r.authController.VerifyEmail)
		auth.POST("/verify/phone", r.authController.VerifyPhone)

		// ========================
		// PROTECTED
		// ========================
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/profile", r.authController.Profile)
			auth.POST("/topup", r.authController.TopUp)

		}
	}
}
