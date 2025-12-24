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
			"error": "Invalid request body",
		})
		return
	}

	userID, ok := ctx.Get("user_id")

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	req.UserID = uid
	driver, err := dc.driverService.RegisterDriver(&req)
	if err != nil {
		if app, ok := err.(*appErr.AppError); ok {
			ctx.JSON(app.Status, gin.H{
				"message": app.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Driver registered successfully",
		"driver":  driver,
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

func (dc *DriverController) ApprovedDriver(ctx *gin.Context) {
	id := ctx.Param("driver_id")
	driverID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid driver ID",
		})
		return
	}

	err = dc.driverService.ApproveDriver(driverID)
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

func (dc *DriverController) RejectDriver(ctx *gin.Context) {
	id := ctx.Param("driver_id")
	driverID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid driver ID",
		})
		return
	}

	err = dc.driverService.RejectDriver(driverID)
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
