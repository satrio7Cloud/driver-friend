package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appErr "be/internal/errors"
	"be/internal/modules/driver/dto"
	"be/internal/modules/driver/service"

	"github.com/google/uuid"
)

type VehicleController struct {
	vehicleService service.VehicleService
}

func NewVehicleController(vehicleService service.VehicleService) *VehicleController {
	return &VehicleController{
		vehicleService: vehicleService,
	}
}

func (vc *VehicleController) RegisterVehicle(ctx *gin.Context) {
	driverIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	driverID, err := uuid.Parse(driverIDVal.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Driver ID",
		})
		return
	}

	var req dto.RegisterVehicle
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	vehicle, err := vc.vehicleService.RegisterVehicle(driverID, req)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{"message": "Invalid"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return

	}
	ctx.JSON(http.StatusOK, vehicle)
}

func (vc *VehicleController) ApproveVehicle(ctx *gin.Context) {
	vehicleId := ctx.Param("vehicle_id")

	vehicleUUID, err := uuid.Parse(vehicleId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid vehicle id",
		})
		return
	}

	err = vc.vehicleService.ApproveVehicle(vehicleUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Vehicle Approved",
	})
}

func (vc *VehicleController) GerDriverVehicle(ctx *gin.Context) {
	driverIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unautorized",
		})
		return
	}

	driverID, err := uuid.Parse(driverIDVal.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Driver Id",
		})
		return
	}

	vehicles, err := vc.vehicleService.GetDriverVehicle(driverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    vehicles,
		"message": "success",
	})
}

func (vc *VehicleController) DeleteVehcile(ctx *gin.Context) {
	vehicleID := ctx.Param("vehicle_id")

	vehicleUUID, err := uuid.Parse(vehicleID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Vehicle ID",
		})
		return
	}

	driverIDVal, exists := ctx.Get("user_id")
	if exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorize",
		})
		return
	}

	driverID, err := uuid.Parse(driverIDVal.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Driver ID",
		})
		return
	}

	err = vc.vehicleService.DeleteVehicle(vehicleUUID, driverID)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{"error": app.Message})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Vehicle has delete",
	})

}
