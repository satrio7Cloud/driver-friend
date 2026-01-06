package service

import (
	appErr "be/internal/errors"
	"be/internal/modules/driver/dto"

	userRepo "be/internal/modules/auth/repository"
	driverRepo "be/internal/modules/driver/repository"
	"be/internal/utils"
)

type DriverAuthService interface {
	LoginDriver(req dto.RequestDriverLogin) (*dto.RequestDriverLogin, error)
}

type driverAuthService struct {
	userRepo   userRepo.UserRepository
	driverRepo driverRepo.DriverRepository
}

func NewDriverAuthService(
	userRepo userRepo.UserRepository,
	driverRepo driverRepo.DriverRepository,
) DriverAuthService {
	return &driverAuthService{
		userRepo:   userRepo,
		driverRepo: driverRepo,
	}
}

func (s *driverAuthService) LoginDriver(req dto.RequestDriverLogin) (*dto.RequestDriverLogin, error) {
	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil || user == nil {
		return nil, appErr.NewAuthorized("Invalid phone or password")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, appErr.NewAuthorized("Invalid phone or passowrd")
	}

	driver, err := s.driverRepo.FindByUserID(user.ID)
	if err != nil || driver == nil {
		return nil, appErr.NewForbiden("user is not registered as driver")
	}

	if driver.Status != "approved" {
		return nil, appErr.NewForbiden("Driver not approved")
	}

	token, err := utils.GenerateToken(
		user.ID.String(),
		driver.ID.String(),
	)
	if err != nil {
		return nil, appErr.NewInternalServerError("Failed to generate token")
	}

	return &dto.DriverLoginResponse{
		UserID:   user.ID.String(),
		DriverID: driver.ID.String(),
		Token:    token,
		Roles:    []string{"driver"},
	}, nil

}
