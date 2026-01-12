package model

import (
	"time"

	"be/internal/modules/driver/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DriverID uuid.UUID `gorm:"type:uuid;not null;index" json:"driver_id"`

	Driver model.Driver `gorm:"foreignKey:DriverID" json:"-"`

	Type      string `gorm:"type:varchar(20);not null;check:type IN ('motor','mobil')" json:"type"`
	Brand     string `gorm:"type:varchar(50);not null" json:"brand"`
	Model     string `gorm:"type:varchar(50)" json:"model"`
	Year      int    `gorm:"not null" json:"year"`
	Plate     string `gorm:"type:varchar(20);index"`
	STNKPhoto string `gorm:"type:text" json:"stnk_photo"`

	Status   string `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','approved','rejected')"`
	IsActive bool   `gorm:"default:true"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
