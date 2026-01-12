package service

import (
	appErr "be/internal/errors"

	driverRepo "be/internal/modules/driver/repository"
	"be/internal/modules/vehicle/model"
	vehicleRepo "be/internal/modules/vehicle/repository"

	"github.com/google/uuid"
)

type VehicleService interface {
	CreateVehicle(driverID uuid.UUID, vehicle *model.Vehicle) (*model.Vehicle, error)
	ApproveVehicle(vehicleID uuid.UUID) error
	GetDriverVehicle(driverID uuid.UUID) ([]model.Vehicle, error)

	DeleteVehicle(vehicleID uuid.UUID, driverID uuid.UUID) error
}

type vehicleService struct {
	vehicleRepo vehicleRepo.VehicleRepository
	driverRepo  driverRepo.DriverRepository
}

func NewVehicleService(
	vehileRepo vehicleRepo.VehicleRepository,
	driverRepo driverRepo.DriverRepository,
) VehicleService {
	return &vehicleService{
		vehicleRepo: vehileRepo,
		driverRepo:  driverRepo,
	}
}

func (s *vehicleService) CreateVehicle(
	driverID uuid.UUID,
	vehicle *model.Vehicle,
) (*model.Vehicle, error) {
	driver, err := s.driverRepo.FindByID(driverID)
	if err != nil {
		return nil, appErr.NewNotFound("Driver Not Found")
	}

	if driver.Status != "approved" {
		return nil, appErr.NewForbidden("Driver not approved")
	}

	existint, err := s.vehicleRepo.FindByDriverID(driverID)
	if len(existint) > 0 {
		return nil, appErr.NewBadRequest("Pengemudi sudah memiliki kendaraan.")
	}

	vehicle.DriverID = driverID
	vehicle.Status = "pending"

	if err := s.vehicleRepo.Create(vehicle); err != nil {
		return nil, appErr.NewInternalServerError("Failed to register vehicle")
	}

	return vehicle, nil
}

func (s *vehicleService) ApproveVehicle(vehicleID uuid.UUID) error {
	vehicle, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return appErr.NewNotFound("Vehicle not found")
	}

	if err := s.vehicleRepo.ApproveVehicle(vehicleID); err != nil {
		return err
	}

	_ = s.driverRepo.SetOnline(vehicle.DriverID, true)
	return nil
}

func (s *vehicleService) GetDriverVehicle(driverID uuid.UUID) ([]model.Vehicle, error) {
	return s.vehicleRepo.FindByDriverID(driverID)
}

func (s *vehicleService) DeleteVehicle(vehicleID uuid.UUID, driverID uuid.UUID) error {
	vehicle, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return appErr.NewNotFound("Vehicle not found")
	}

	if vehicle.DriverID != driverID {
		return appErr.NewForbidden("Not Authorized")
	}

	if err := s.vehicleRepo.Delete(vehicleID); err != nil {
		return err
	}

	_ = s.driverRepo.SetOnline(driverID, false)

	return nil

}
