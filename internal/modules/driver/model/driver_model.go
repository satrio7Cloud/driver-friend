package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Driver struct {
	ID     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID *uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"user_id,omitempty"`

	// Personal Information
	FullName         string     `gorm:"type:varchar(100);not null" json:"full_name"`
	NIK              string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"nik"`
	Address          string     `gorm:"type:varchar(255);not null" json:"address"`
	DOB              *time.Time `json:"dob"`
	Gender           string     `gorm:"type:varchar(10)" json:"gender"`
	Phone            string     `gorm:"type:varchar(20);not null" json:"phone"`
	EmergencyContact string     `gorm:"type:varchar(20)" json:"emergency_contact"`

	// Documents
	SelfiePhoto string `gorm:"type:text" json:"selfie_photo"`
	KTPPhoto    string `gorm:"type:text" json:"ktp_photo"`
	SIMPhoto    string `gorm:"type:text" json:"sim_photo"`
	SKCKPhoto   string `gorm:"type:text" json:"skck_photo"`

	// Status string `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Status string `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','approved','rejected')" json:"status"`

	// Driver metadata
	IsOnline bool `gorm:"default:false" json:"is_online"`

	CurrentLat *float64 `json:"current_lat"`
	CurrentLng *float64 `json:"current_lng"`

	Rating     float64 `gorm:"default:5" json:"rating"`
	TotalTrips int     `gorm:"default:0" json:"total_trips"`
	CancelRate float64 `gorm:"default:0" json:"cancel_rate"`

	// timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Driver) BeforeCreate(tsx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return
}
