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
	RejectDriver(id uuid.UUID, message string) error
	UpdateDriver(driver *model.Driver) error
	DeleteDriver(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status string) error
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
	existing, err := s.repo.FindByPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, appErr.NewBadRequest("phone alreadey registered as driver")
	}

	req.Status = "pending"
	req.IsOnline = false

	if err = s.repo.Create(req); err != nil {
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
	return s.repo.UpdateStatus(id, "Approved")
}

func (s *driverService) RejectDriver(id uuid.UUID, message string) error {
	return s.repo.UpdateStatus(id, "Rejected")
}

func (s *driverService) UpdateDriver(driver *model.Driver) error {
	return s.repo.Update(driver)
}

func (s *driverService) UpdateStatus(id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *driverService) DeleteDriver(id uuid.UUID) error {
	return s.repo.Delete(id)
}
