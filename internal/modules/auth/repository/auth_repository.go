package repository

import (
	"be/internal/modules/auth/model"
	// "be/internal/modules/role/model"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error

	AsignRole(userID, roleID uuid.UUID) error

	FindByEmail(email string) (*model.User, error)
	FindByEmailWithRoles(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	FindPhoneWithRoles(phone string) (*model.User, error)
	FindById(id string) (*model.User, error)

	FindByGoogleID(googleID string) (*model.User, error)
	FindByAppleID(appleID string) (*model.User, error)
	FindByFacebookID(facebookID string) (*model.User, error)

	FindProfileById(id string) (*model.User, error)

	UpdateDeviceInfo(id string, deviceID string, fcmToken string) error

	UpdateLoginStatus(user *model.User) error

	UpdateLocation(id string, lat float64, long float64) error

	VerifyEmail(userId string) error
	VerifyPhone(userId string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByIdentifier(identifier string) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		Where("email = ? OR phone = ?", identifier, identifier).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// for test insert
func (r *userRepository) AsignRole(userID, roleID uuid.UUID) error {
	return r.db.Exec(`
	INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`,
		userID, roleID,
	).Error
}

// func (r *userRepository) AsignRole(user *model.User, role *model.Role) error {
// 	return r.db.Model(user).Association("Role").Append(role)
// }

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.
		Where("email = ?", email).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindByEmailWithRoles(email string) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		Where("email = ?", email).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.
		Preload("Role").
		Where("phone = ?", phone).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindPhoneWithRoles(phone string) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		Where("phone = ?", phone).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindById(id string) (*model.User, error) {
	var user model.User
	err := r.db.
		Where("id = ?", id).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByAppleID(appleID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("apple_id = ?", appleID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindByFacebookID(facebookID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("facebook_id = ?", facebookID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindProfileById(id string) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("id = ?", id).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) UpdateDeviceInfo(id string, deviceID string, fcmToken string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"device":    deviceID,
			"fcm_token": fcmToken,
		}).Error
}

func (r *userRepository) UpdateLoginStatus(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateLocation(id string, lat float64, lot float64) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"latitude":  lat,
			"longitude": lot,
		}).Error
}

func (r *userRepository) VerifyEmail(userID string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_email_verified": true,
			"updated_at":        gorm.Expr("Now()"),
		}).Error
}

func (r *userRepository) VerifyPhone(userID string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_phone_verified": true,
			"updated_at":        gorm.Expr("Now()"),
		}).Error

}
