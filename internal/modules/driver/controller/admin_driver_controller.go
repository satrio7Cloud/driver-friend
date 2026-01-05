package controller

import (
	appErr "be/internal/errors"
	"be/internal/modules/driver/service"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminDriverController struct {
	driverService service.DriverService
}

func NewAdminDriverController(driverService service.DriverService) *AdminDriverController {
	return &AdminDriverController{
		driverService: driverService,
	}
}

func (c *AdminDriverController) ApproveDriver(ctx *gin.Context) {
	driverIDParam := ctx.Param("id")

	driverID, err := uuid.Parse(driverIDParam)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid driver",
		})
		return
	}

	if err := c.driverService.ApproveDriver(driverID); err != nil {
		appErr.NewBadRequest(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (c *AdminDriverController) RejectDriver(ctx *gin.Context) {
	driverIDParam := ctx.Param("id")

	driverID, err := uuid.Parse(driverIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error invalid",
		})
		return
	}

	if err := c.driverService.RejectDriver(driverID, "reject by admin"); err != nil {
		appErr.NewInternalServerError(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})

}
