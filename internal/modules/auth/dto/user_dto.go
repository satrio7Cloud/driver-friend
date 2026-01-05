package dto

type LoginRequestAdmin struct {
	Identifier string `json:"identifier" binding:"required"` // email / phone
	Password   string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"required"`
}

type VerifyRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginPhoneRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Balance int64  `json:"balance"`
}

type AuthResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Phone   string   `json:"phone"`
	Token   string   `json:"token"`
	Balance int64    `json:"balance"`
	Roles   []string `json:"role"`
}

type ProfileResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Balance int64  `json:"balance"`

	IsPhoneVerified bool   `json:"is_phone_verified"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Language        string `json:"language"`
	NotifEnabled    bool   `json:"notif_enabled"`
	DarkMode        bool   `json:"dark_mode"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TopUp struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type WalletResponse struct {
	Balance          int64  `json:"balance"`
	BalanceUpdatedAt string `json:"balance_updated_at"`
}
