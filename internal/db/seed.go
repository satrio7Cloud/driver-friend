package db

import (
	roleModel "be/internal/modules/role/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []string{
		"super admin",
		"admin",
		"customer",
		"driver",
	}

	for _, roleName := range roles {
		var count int64
		db.Model(&roleModel.Role{}).
			Where("name = ?", roleName).
			Count(&count)

		if count == 0 {
			db.Create(&roleModel.Role{
				ID:   uuid.New(),
				Name: roleName,
			})
		}
	}
	return nil
}
