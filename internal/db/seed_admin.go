package db

import (
	authModel "be/internal/modules/auth/model"
	roleModel "be/internal/modules/role/model"

	"be/internal/utils"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	if os.Getenv("ENABLE_ADMIN_SEED") != "true" {
		log.Println("admin seed disabled")
		return nil
	}

	app := os.Getenv("APP_ENV")
	if app == "production" {
		log.Println("skipping admin seed in production")
		return nil
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	phone := os.Getenv("ADMIN_PHONE")
	name := os.Getenv("ADMIN_NAME")

	if email == "" || password == "" {
		log.Println("admin env not set, skipping admin seed")
		return nil
	}

	var adminRole roleModel.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var count int64
	if err := db.Model(&authModel.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("admin already exists")
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	admin := authModel.User{
		ID:              uuid.New(),
		Name:            name,
		Email:           email,
		Phone:           phone,
		Password:        hashedPassword,
		IsEmailVerified: true,
		IsPhoneVerified: true,
		IsActive:        true,
		Role:            []roleModel.Role{adminRole},
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("admin user seeded success")
	return nil

}
