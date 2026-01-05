package repository

import (
	"be/internal/modules/driver/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverRepository interface {
	Create(driver *model.Driver) error
	FindByID(id uuid.UUID) (*model.Driver, error)
	FindByUserID(userID uuid.UUID) (*model.Driver, error)
	FindByPhone(phone string) (*model.Driver, error)
	Update(driver *model.Driver) error
	UpdateStatus(id uuid.UUID, status string) error
	Delete(id uuid.UUID) error
}

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{
		db: db,
	}
}

func (r *driverRepository) Create(driver *model.Driver) error {
	return r.db.Create(driver).Error
}

func (r *driverRepository) FindByID(id uuid.UUID) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.First(&driver, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (r *driverRepository) FindByUserID(userID uuid.UUID) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.Where("user_id = ?", userID).First(&driver).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) FindByPhone(phone string) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.Where("phone = ?", phone).First(&driver).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &driver, nil

}

func (r *driverRepository) Update(driver *model.Driver) error {
	return r.db.Save(driver).Error
}

func (r *driverRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Driver{}).
		Where("id = ?").
		Update("status", status).Error
}

func (r *driverRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Driver{}, "id = ?").Error
}
