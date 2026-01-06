package dto

type ApplyDriverRequest struct {
	FullName string `json:"full_name" binding:"required"`
	NIK      string `json:"nik" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

type RequestDriverLogin struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DriverLoginResponse struct { // nanti di hapus
	UserID   string   `json:"user_id"`
	DriverID string   `json:"driver_id"`
	Token    string   `json:"token"`
	Roles    []string `json:"roles"`
}
