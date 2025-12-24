package model

import (
	"github.com/google/uuid"
)

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID uuid.UUID `gorm:"type:uuid;primaryKey"`
}
