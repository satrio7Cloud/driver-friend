package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"be/internal/modules/vehicle/model"
)

type VehicleRepository interface {
	Create(vehicle *model.Vehicle) error
	FindByID(id uuid.UUID) (*model.Vehicle, error)
	FindByDriverID(driverID uuid.UUID) ([]model.Vehicle, error)
	Update(vehicle *model.Vehicle) error
	Delete(id uuid.UUID) error
	ApproveVehicle(id uuid.UUID) error
}

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(vehicle *model.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *vehicleRepository) FindByID(id uuid.UUID) (*model.Vehicle, error) {
	var vehicle model.Vehicle
	err := r.db.First(&vehicle, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *vehicleRepository) FindByDriverID(driverID uuid.UUID) ([]model.Vehicle, error) {
	var vehicle []model.Vehicle
	err := r.db.Where("driver_id = ?", driverID).Find(&vehicle).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (r *vehicleRepository) Update(vehicle *model.Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *vehicleRepository) ApproveVehicle(id uuid.UUID) error {
	return r.db.Model(&model.Vehicle{}).
		Where("id = ? ", id).
		Updates(map[string]interface{}{
			"status":     "approved",
			"updated_at": gorm.Expr("Now()"),
		}).Error
}

func (r *vehicleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Vehicle{}, "id = ?").Error
}
