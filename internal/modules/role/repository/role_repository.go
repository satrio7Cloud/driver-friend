package repository

import (
	"be/internal/modules/role/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(name string) (*model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("name = ?", name).
		First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
