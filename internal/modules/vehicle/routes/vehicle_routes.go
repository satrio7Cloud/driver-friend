package routes

import (
	"be/internal/middleware"
	"be/internal/modules/vehicle/controller"

	"github.com/gin-gonic/gin"
)

type VehicleRoutes struct {
	vehicleController *controller.VehicleController
}

func NewVehicleRoutes(vehicleController *controller.VehicleController) *VehicleRoutes {
	return &VehicleRoutes{
		vehicleController: vehicleController,
	}
}

func (v *VehicleRoutes) RegisterRoutes(router *gin.RouterGroup) {
	driver := router.Group("/driver")
	{
		driver.Use(
			middleware.AuthMiddleware(),
			middleware.OnlyDriver())
		{
			driver.POST("/vehicles", v.vehicleController.RegisterVehicle)
			driver.GET("/vehicles", v.vehicleController.GetDriverVehicle)
			driver.DELETE("/vehicles/:vehicle_id", v.vehicleController.DeleteVehicle)
		}

		admin := router.Group("/admin")
		admin.Use(
			middleware.AuthMiddleware(),
			middleware.OnlyAdmin(),
		)
		{
			admin.PUT("/vehicles/:id/approve", v.vehicleController.ApproveVehicle)
		}

	}
}
