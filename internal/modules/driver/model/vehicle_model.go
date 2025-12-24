package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DriverID uuid.UUID `gorm:"type:uuid;not null;index" json:"driver_id"`

	Type      string `gorm:"type:varchar(20);not null" json:"type"` // motor/mobil
	Brand     string `gorm:"type:varchar(50);not null" json:"brand"`
	Model     string `gorm:"type:varchar(50)" json:"model"`
	Year      int    `gorm:"not null" json:"year"`
	Plate     string `gorm:"type:varchar(20);uniqueIndex" json:"plate"`
	STNKPhoto string `gorm:"type:text" json:"stnk_photo"`

	Status string `gorm:"type:varchar(20);default:'pending'" json:"status"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
