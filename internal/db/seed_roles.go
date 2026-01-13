package db

import (
	roleModel "be/internal/modules/role/model"

	"log"

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
		if err := db.Model(&roleModel.Role{}).
			Where("name = ?", roleName).
			Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			if err := db.Create(&roleModel.Role{
				ID:   uuid.New(),
				Name: roleName,
			}).Error; err != nil {
				return err
			}
			log.Println("role created:", roleName)
		}
	}

	return nil
}
