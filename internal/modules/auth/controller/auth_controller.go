package controller

import (
	appErr "be/internal/errors"
	"be/internal/modules/auth/dto"
	"be/internal/modules/auth/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register
func (ac *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := ac.authService.Register(&req)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"error":   app.Code,
				"message": app.Message,
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Success", "data": res})
}

// Verify Email
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err := ac.authService.VerifyEmail(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfull",
	})
}

// Verify Phone
func (ac *AuthController) VerifyPhone(ctx *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request",
		})
		return
	}

	err := ac.authService.VerifyPhone(req.UserID)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"error":   app.Code,
				"message": app.Message,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Phone verified successfully",
	})
}

// Login Admin
func (ac *AuthController) LoginAdmin(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request",
		})
		return
	}

	res, err := ac.authService.Login(&req)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"error":   app.Code,
				"message": app.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	isAdmin := false
	for _, role := range res.Roles {
		if role == "admin" || role == "super admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		ctx.JSON(http.StatusOK, gin.H{
			"error":   "FORBIDEN",
			"message": "Admin access only",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    res,
	})

}

// Login by email
func (ac *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := ac.authService.Login(&req)
	if err != nil {

		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"error":   app.Code,
				"message": app.Message,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successfull", "data": res})
}

// Login by Phone
func (ac *AuthController) LoginPhone(ctx *gin.Context) {
	var req dto.LoginPhoneRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request",
		})
		return
	}

	res, err := ac.authService.LoginByPhone(&req)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"error":   app.Code,
				"message": app.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login Successfull",
		"data":    res,
	})

}

// Profile
func (ac *AuthController) Profile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := ac.authService.GetProfile(userID.(string))
	if err != nil {

		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"error":   app.Code,
				"message": app.Message,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// Top Up
func (ac *AuthController) TopUp(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.TopUp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	res, err := ac.authService.TopUp(userID.(string), req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Top Up Success",
		"data":    res,
	})
}
