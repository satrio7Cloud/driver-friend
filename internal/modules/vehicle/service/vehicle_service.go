package service

import (
	appErr "be/internal/errors"

	// "be/internal/modules/vehicle/dto"
	"be/internal/modules/vehicle/model"
	"be/internal/modules/vehicle/repository"

	"github.com/google/uuid"
)

type VehicleService interface {
	RegisterVehicle(vehicle *model.Vehicle) (*model.Vehicle, error)
	ApproveVehicle(vehicleID uuid.UUID) error
	GetDriverVehicle(driverID uuid.UUID) ([]model.Vehicle, error)

	DeleteVehicle(vehicleID uuid.UUID, driverID uuid.UUID) error
}

type vehicleService struct {
	vehicleRepo repository.VehicleRepository
}

func NewVehicleService(vehileRepo repository.VehicleRepository) VehicleService {
	return &vehicleService{
		vehicleRepo: vehileRepo,
	}
}

func (s *vehicleService) RegisterVehicle(vehicle *model.Vehicle) (*model.Vehicle, error) {
	vehicle.Status = "pending"

	if err := s.vehicleRepo.Create(vehicle); err != nil {
		return nil, appErr.NewAuthorized("failed to register vehicle")
	}
	return vehicle, nil
}

func (s *vehicleService) ApproveVehicle(vehicleID uuid.UUID) error {
	return s.vehicleRepo.ApproveVehicle(vehicleID)
}

func (s *vehicleService) GetDriverVehicle(driverID uuid.UUID) ([]model.Vehicle, error) {
	return s.vehicleRepo.FindByDriverID(driverID)
}

func (s *vehicleService) DeleteVehicle(vehicleID uuid.UUID, driverID uuid.UUID) error {
	vehicle, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return nil
	}

	if vehicle.DriverID != driverID {
		return appErr.NewBadRequest("not authorized to delete this vehicle")
	}

	return s.vehicleRepo.Delete(vehicleID)
}
