package repository

import (
	"fmt"
	"time"

	"be/internal/db"

	"github.com/google/uuid"
)

type otpRepository struct{}

func NewOTPRepository() OTPRepository {
	return &otpRepository{}
}

func redisKey(userID uuid.UUID) string {
	return fmt.Sprintf("otp:driver:%s", userID.String())
}

func (r *otpRepository) Save(userID uuid.UUID, otp string) error {
	return db.Redis.Set(
		db.Ctx,
		redisKey(userID),
		otp,
		5*time.Minute,
	).Err()
}

func (r *otpRepository) Verify(userID uuid.UUID, otp string) (bool, error) {
	val, err := db.Redis.Get(db.Ctx, redisKey(userID)).Result()
	if err != nil {
		return false, err
	}

	return val == otp, nil
}

func (r *otpRepository) Delete(userID uuid.UUID) error {
	return db.Redis.Del(db.Ctx, redisKey(userID)).Err()
}
