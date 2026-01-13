package db

import (
	authModel "be/internal/modules/auth/model"
	driverModel "be/internal/modules/driver/model"
	roleModel "be/internal/modules/role/model"
	vehicleModel "be/internal/modules/vehicle/model"

	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(host, port, user, password, dbname string) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection established")

	err = db.AutoMigrate(

		&authModel.User{},
		&roleModel.Role{},
		&vehicleModel.Vehicle{},
		&driverModel.Driver{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully")

	if err := SeedRoles(db); err != nil {
		log.Fatal("Failed to seed roles:", err)
	}

	if err := SeedAdmin(db); err != nil {
		log.Fatal("Failed to seed admin:", err)
	}

	DB = db
}
