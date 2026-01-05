package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appErr "be/internal/errors"
	"be/internal/modules/vehicle/model"
	"be/internal/modules/vehicle/service"

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
	driverIDStr, existing := ctx.Get("driver_id")

	if !existing {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	driverID, _ := uuid.Parse(driverIDStr.(string))
	var req model.Vehicle
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	req.DriverID = driverID

	vehicle, err := vc.vehicleService.RegisterVehicle(&req)
	if err != nil {
		appErr.NewInternalServerError(err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success",
		"data":    vehicle,
	})
}

func (vc *VehicleController) GetMyVehicle(ctx *gin.Context) {
	driverIDStr, exists := ctx.Get("driver_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	driverID, _ := uuid.Parse(driverIDStr.(string))

	vehicle, err := vc.vehicleService.GetDriverVehicle(driverID)
	if err != nil {
		appErr.NewInternalServerError(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, vehicle)
}

func (vc *VehicleController) GetDriverVehicle(ctx *gin.Context) {
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

func (vc *VehicleController) DeleteVehicle(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid vehicle",
		})
		return
	}

	driverIDStr, exists := ctx.Get("driver_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	driverID, _ := uuid.Parse(driverIDStr.(string))
	if err := vc.vehicleService.DeleteVehicle(vehicleID, driverID); err != nil {
		appErr.NewInternalServerError(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (vc *VehicleController) ApproveVehicle(ctx *gin.Context) {
	vehicleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Vehicle",
		})
		return
	}

	if err := vc.vehicleService.ApproveVehicle(vehicleID); err != nil {
		appErr.NewInternalServerError(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}
