package routes

import (
	"be/internal/middleware"
	"be/internal/modules/driver/controller"

	"github.com/gin-gonic/gin"
)

type AdminDriverRoutes struct {
	adminController *controller.AdminDriverController
}

func NewAdminDriverRoutes(adminController *controller.AdminDriverController) *AdminDriverRoutes {
	return &AdminDriverRoutes{adminController: adminController}
}

func (r *AdminDriverRoutes) RegisterRoutes(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.OnlyAdmin(),
	)
	{
		admin.PUT("/drivers/:id/approve", r.adminController.ApproveDriver)
		admin.PUT("/drivers/:id/reject", r.adminController.RejectDriver)
	}
}
