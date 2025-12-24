package service

import (
	appErr "be/internal/errors"

	"be/internal/modules/driver/dto"
	"be/internal/modules/driver/model"
	"be/internal/modules/driver/repository"

	"github.com/google/uuid"
)

type VehicleService interface {
	RegisterVehicle(driverID uuid.UUID, req dto.RegisterVehicle) (*model.Vehicle, error)
	ApproveVehicle(driverID uuid.UUID) error
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

func (s *vehicleService) RegisterVehicle(driverID uuid.UUID, req dto.RegisterVehicle) (*model.Vehicle, error) {
	if req.Type != "motor" && req.Type != "mobil" {
		return nil, appErr.NewNotFound("Invalid vehicle type")
	}

	vehicle := &model.Vehicle{
		DriverID:  driverID,
		Type:      req.Type,
		Brand:     req.Brand,
		Model:     req.Model,
		Year:      req.Year,
		Plate:     req.Plate,
		STNKPhoto: req.STNKPhoto,
		Status:    "pending",
	}

	err := s.vehicleRepo.Create(vehicle)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *vehicleService) ApproveVehicle(vehicleID uuid.UUID) error {
	return s.vehicleRepo.ApproveVehicle(vehicleID)
}

func (s *vehicleService) GetDriverVehicle(vehicleID uuid.UUID) ([]model.Vehicle, error) {
	return s.vehicleRepo.FindByDriverID(vehicleID)
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
