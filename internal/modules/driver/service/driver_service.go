package service

import (
	appErr "be/internal/errors"
	"be/internal/modules/driver/model"
	"be/internal/modules/driver/repository"

	"github.com/google/uuid"
)

type DriverService interface {
	RegisterDriver(req *model.Driver) (*model.Driver, error)
	GetDriverByID(id uuid.UUID) (*model.Driver, error)
	GetDriverByUserID(useID uuid.UUID) (*model.Driver, error)
	ApproveDriver(id uuid.UUID) error
	RejectDriver(id uuid.UUID) error
	UpdateDriver(driver *model.Driver) error
	DeleteDriver(id uuid.UUID) error
}

type driverService struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverService{
		repo: repo,
	}
}

func (s *driverService) RegisterDriver(req *model.Driver) (*model.Driver, error) {
	exists, _ := s.repo.FindByUserID(req.UserID)

	if exists != nil {
		return nil, appErr.NewBadRequest("Driver already registered")
	}

	req.Status = "pending"

	err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return req, nil

}

func (s *driverService) GetDriverByID(id uuid.UUID) (*model.Driver, error) {
	driver, err := s.repo.FindByID(id)
	if err != nil {
		return nil, appErr.NewNotFound("Driver not found")
	}
	return driver, nil
}

func (s *driverService) GetDriverByUserID(userID uuid.UUID) (*model.Driver, error) {
	driver, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, appErr.NewNotFound("Driver not found")
	}
	return driver, nil
}

func (s *driverService) ApproveDriver(id uuid.UUID) error {
	return s.repo.UpdateStatus(id, "approved")
}

func (s *driverService) RejectDriver(id uuid.UUID) error {
	return s.repo.UpdateStatus(id, "Rejected")
}

func (s *driverService) UpdateDriver(driver *model.Driver) error {
	return s.repo.Update(driver)
}

func (s *driverService) DeleteDriver(id uuid.UUID) error {
	return s.repo.Delete(id)
}
