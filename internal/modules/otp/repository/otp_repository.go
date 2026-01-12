package repository

import "github.com/google/uuid"

type OTPRepository interface {
	Save(userID uuid.UUID, otp string) error
	Verify(userID uuid.UUID, otp string) (bool, error)
	Delete(userID uuid.UUID) error
}
