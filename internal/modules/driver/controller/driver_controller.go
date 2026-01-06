package controller

import (
	appErr "be/internal/errors"
	"be/internal/modules/driver/model"
	"be/internal/modules/driver/service"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DriverController struct {
	driverService service.DriverService
}

func NewDriverController(driverService service.DriverService) *DriverController {
	return &DriverController{
		driverService: driverService,
	}
}

func (dc *DriverController) RegisterDriver(ctx *gin.Context) {
	var req model.Driver
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request Payload",
		})
		return
	}

	driver := model.Driver{
		FullName:         req.FullName,
		NIK:              req.NIK,
		Phone:            req.Phone,
		Address:          req.Address,
		Gender:           req.Gender,
		EmergencyContact: req.EmergencyContact,
	}

	result, err := dc.driverService.RegisterDriver(&driver)
	if err != nil {
		ctx.JSON(appErr.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return

	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success",
		"data":    result,
	})
}

func (dc *DriverController) GetDriverProfile(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	uid, _ := uuid.Parse(userID.(string))

	driver, err := dc.driverService.GetDriverByID(uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Driver not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    driver,
	})
}

func (dc *DriverController) DeleteDriver(ctx *gin.Context) {
	id := ctx.Param("driver_id")
	driverID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid driver ID",
		})
		return
	}

	err = dc.driverService.DeleteDriver(driverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})

}
