package routes

import (
	// "be/internal/middleware"
	"be/internal/modules/driver/controller"

	"github.com/gin-gonic/gin"
)

type DriverRoutes struct {
	driverController *controller.DriverController
}

func NewDriverRoutes(driverController *controller.DriverController) *DriverRoutes {
	return &DriverRoutes{
		driverController: driverController,
	}
}

func (d *DriverRoutes) RegisterRoutes(router *gin.RouterGroup) {
	driver := router.Group("/driver")
	{
		driver.Use()
		{
			driver.POST("/", d.driverController.RegisterDriver)
			driver.GET("/profile", d.driverController.GetDriverProfile)
			driver.PUT("/:driver_id/approve", d.driverController.ApprovedDriver)
			driver.PUT("/:driver_id/reject", d.driverController.RejectDriver)
			driver.DELETE("/:driver_id", d.driverController.DeleteDriver)
		}
	}
}
