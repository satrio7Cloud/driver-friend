package utils

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// FIX (DUMMY dulu, production nanti pakai SMS gateway)
func SendOTP(phone, otp string) {
	log.Printf("Send OTP %s to Phone %s\n", otp, phone)
}
