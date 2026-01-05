package dto

type ApplyDriverRequest struct {
	FullName string `json:"full_name" binding:"required"`
	NIK      string `json:"nik" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}
