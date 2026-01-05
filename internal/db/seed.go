package db

import (
	authModel "be/internal/modules/auth/model"
	roleModel "be/internal/modules/role/model"
	"be/internal/utils"

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

func SeedAdmin(db *gorm.DB) error {
	var adminRole roleModel.Role

	if err := db.Where("name =  ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var count int64
	db.Model(&authModel.User{}).
		Where("email = ?", "admin@ojol.com").
		Count(&count)

	if count > 0 {
		log.Println("admin user already exists")
		return nil
	}

	hashedPasswod, err := utils.HashPassword("admin123")
	if err != nil {
		return nil
	}

	admin := authModel.User{
		ID:              uuid.New(),
		Name:            "System Admin",
		Email:           "admin@ojol.com",
		Phone:           "085156419000",
		Password:        hashedPasswod,
		IsEmailVerified: true,
		IsPhoneVerified: true,
		IsActive:        true,
		Role:            []roleModel.Role{adminRole},
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("admin user seeded successfully")
	return nil

}
