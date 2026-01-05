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
	vehicle := router.Group("/driver")
	{
		vehicle.Use(middleware.OnlyDriver())
		{
			vehicle.POST("/vehicles", v.vehicleController.RegisterVehicle)
			// vehicle.PUT("/:vehicle_id/approve", v.vehicleController.ApproveVehicle)
			vehicle.GET("/vehicles", v.vehicleController.GetDriverVehicle)
			vehicle.DELETE("/vehicles/:vehicle_id", v.vehicleController.DeleteVehicle)
		}

		admin := router.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.OnlyAdmin())
		{
			admin.PUT("/vehicles/:id/approve", v.vehicleController.ApproveVehicle)
		}

	}
}
