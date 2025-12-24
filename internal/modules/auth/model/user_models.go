package model

import (
	roleModel "be/internal/modules/role/model"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Email    string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone    string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Password string    `gorm:"type:varchar(255)" json:"-"`

	// Verification
	IsPhoneVerified bool `gorm:"default:false" json:"is_phone_verified"`
	IsEmailVerified bool `gorm:"default:false" json:"is_email_verified"`

	// Social Login
	GoogleID   *string `gorm:"type:varchar(100)" json:"google_id,omitempty"`
	AppleID    *string `gorm:"type:varchar(100)" json:"apple_id,omitempty"`
	FacebookID *string `gorm:"type:varchar(100)" json:"facebook_id,omitempty"`

	// Login Info
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	FailedLoginAttempt int        `gorm:"default:0" json:"failed_login_attempt"`
	IsBlocked          bool       `gorm:"default:false" json:"is_blocked"`

	// Wallet Info
	Balance          int64      `gorm:"type:bigint;default:0" json:"balance"`
	BalanceUpdatedAt time.Time  `json:"balance_updated_at"`
	Pin              *string    `gorm:"type:varchar(255)" json:"-"` // hashed PIN
	PinUpdatedAt     *time.Time `json:"pin_updated_at,omitempty"`

	// Preferences
	Language     string `gorm:"type:varchar(10);default:'id'" json:"language"`
	NotifEnabled bool   `gorm:"default:true" json:"notif_enabled"`
	DarkMode     bool   `gorm:"default:false" json:"dark_mode"`

	// Device Data
	DeviceID   *string `gorm:"type:varchar(100)" json:"device_id,omitempty"`
	FCMToken   *string `gorm:"type:text" json:"fcm_token,omitempty"`
	AppVersion *string `gorm:"type:varchar(20)" json:"app_version,omitempty"`
	OSPlatform *string `gorm:"type:varchar(20)" json:"os_platform,omitempty"` // android/ios
	OSVersion  *string `gorm:"type:varchar(20)" json:"os_version,omitempty"`

	// Location
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Role     string `gorm:"type:varchar(20);default:'customer'" json:"role"`
	Role     []roleModel.Role `gorm:"many2many:user_roles;" json:"role"`
	IsActive bool             `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
