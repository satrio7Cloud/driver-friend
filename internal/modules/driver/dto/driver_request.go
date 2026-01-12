package dto

import "github.com/google/uuid"

type ApplyDriverRequest struct {
	FullName string `json:"full_name" binding:"required"`
	NIK      string `json:"nik" binding:"required"`
	Address  string `json:"addressS" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

type DriverLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type DriverRequestOTP struct {
	Phone string `json:"phone" binding:"required"`
}

type DriverVerifyOTP struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type DriverLoginResponse struct {
	Token  string        `json:"token"`
	Driver DriverProfile `json:"driver"`
}

type DriverProfile struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Status   string    `json:"status"`
	IsOnline bool      `json:"is_online"`
}
