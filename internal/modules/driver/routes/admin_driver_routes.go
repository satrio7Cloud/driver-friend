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
	admin := router.Group("/admin/drivers")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.OnlyAdmin(),
	)
	{
		admin.GET("/pending", r.adminController.GetPendingDrivers)
		admin.PUT("/:id/approve", r.adminController.ApproveDriver)
		admin.PUT("/drivers/:id/reject", r.adminController.RejectDriver)
	}
}
