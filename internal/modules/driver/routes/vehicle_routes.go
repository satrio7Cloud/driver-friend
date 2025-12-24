package routes

import (
	"be/internal/middleware"
	"be/internal/modules/driver/controller"

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
	vehicle := router.Group("/vehicle")
	{
		vehicle.Use(middleware.AuthMiddleware(), middleware.OnlyDriver())
		{
			vehicle.POST("/", v.vehicleController.RegisterVehicle)
			vehicle.PUT("/:vehicle_id/approve", v.vehicleController.ApproveVehicle)
			vehicle.GET("/", v.vehicleController.GerDriverVehicle)
			vehicle.DELETE("/:vehicle_id", v.vehicleController.DeleteVehcile)
		}
	}
}
