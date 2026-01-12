package service

import (
	appErr "be/internal/errors"
	"be/internal/modules/auth/dto"

	authModel "be/internal/modules/auth/model"
	driverRepository "be/internal/modules/driver/repository"
	otpRepo "be/internal/modules/otp/repository"
	roleModel "be/internal/modules/role/model"
	roleRepository "be/internal/modules/role/repository"

	"be/internal/modules/auth/repository"
	"be/internal/utils"

	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(input *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RequestDriverOTP(req *dto.DriverLoginRequest) error
	VerifyDriverOTP(req *dto.VerifyOTPRequest) (*dto.DriverAuthResponse, error)
	LoginDriver(req *dto.DriverLoginRequest) error
	Login(input *dto.LoginRequest) (*dto.AuthResponse, error)
	LoginByPhone(input *dto.LoginPhoneRequest) (*dto.AuthResponse, error)
	GetProfile(userID string) (*dto.ProfileResponse, error)

	VerifyEmail(userId string) error
	VerifyPhone(userId string) error

	TopUp(userID string, amount int64) (*dto.WalletResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	roleRepo   roleRepository.RoleRepository
	driverRepo driverRepository.DriverRepository
	otpRepo    otpRepo.OTPRepository
	jwtSecret  []byte
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo roleRepository.RoleRepository,
	driverRepo driverRepository.DriverRepository,
	otpRepo otpRepo.OTPRepository,
	jwtSecret string) AuthService {
	return &authService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		driverRepo: driverRepo,
		otpRepo:    otpRepo,
		jwtSecret:  []byte(jwtSecret),
	}
}

// REGISTER (Customer)
func (s *authService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, appErr.NewBadRequest("email already in use")
	}

	existingPhone, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	if existingPhone != nil {
		return nil, appErr.NewBadRequest("phone number already in use")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &authModel.User{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     hashed,
		Language:     "id",
		NotifEnabled: true,
		DarkMode:     false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	customerRole, err := s.roleRepo.FindByName("customer")
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.AsignRole(user.ID, customerRole.ID); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}, nil
}

// Verifikasi Email
func (s *authService) VerifyEmail(userID string) error {
	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return appErr.NewBadRequest("user not found")
	}

	if user.IsEmailVerified {
		return nil
	}

	return s.userRepo.VerifyEmail(userID)

}

// Verifikasi Phone
func (s *authService) VerifyPhone(userID string) error {
	user, err := s.userRepo.FindById(userID)

	if err != nil {
		return nil
	}

	if user == nil {
		return appErr.NewAuthorized("User not found")
	}

	if user.IsPhoneVerified {
		return nil
	}

	return s.userRepo.VerifyPhone(userID)
}

func extractRoleNames(roles []roleModel.Role) []string {
	var result []string
	for _, r := range roles {
		result = append(result, r.Name)
	}
	return result
}

// LOGIN by email (Customer)
func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {

	identifier := strings.TrimSpace(req.Email)
	var user *authModel.User
	var err error

	if strings.Contains(identifier, "@") {
		user, _ = s.userRepo.FindByEmailWithRoles(identifier)
	} else {
		user, _ = s.userRepo.FindPhoneWithRoles(identifier)
	}

	if user == nil {
		return nil, appErr.NewBadRequest("account not found")
	}

	if user.IsBlocked {
		return nil, appErr.NewAuthorized("akun anda diblokir sementara, hubungi support")
	}

	if strings.Contains(identifier, "@") && !user.IsEmailVerified {
		return nil, appErr.NewVerificationRequired(
			"email not verified",
		)
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		user.FailedLoginAttempt += 1

		if user.FailedLoginAttempt >= 5 {
			user.IsBlocked = true
		}

		_ = s.userRepo.UpdateLoginStatus(user)
		return nil, appErr.NewBadRequest("Password Wrong")
	}

	// user failed login
	user.FailedLoginAttempt = 0
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.UpdateLoginStatus(user)

	roleNames := extractRoleNames(user.Role)
	if len(roleNames) == 0 {
		return nil, appErr.NewBadRequest("user has no role assigned")
	}

	token, err := utils.GenerateToken(
		user.ID.String(),
		roleNames,
		nil,
		false,
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:      user.ID.String(),
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Token:   token,
		Balance: user.Balance,
		Roles:   roleNames,
	}, nil

}

// Login by phone (Customer)
func (s *authService) LoginByPhone(req *dto.LoginPhoneRequest) (*dto.AuthResponse, error) {
	identifier := strings.TrimSpace(req.Phone)
	var user *authModel.User
	var err error

	if strings.Contains(identifier, "@") {
		user, _ = s.userRepo.FindByEmailWithRoles(identifier)
	} else {
		user, _ = s.userRepo.FindPhoneWithRoles(identifier)
	}

	if user == nil {
		return nil, appErr.NewBadRequest("Account Not Found")
	}

	if user.IsBlocked {
		return nil, appErr.NewAuthorized("Akun anda di blokir")
	}

	if !user.IsPhoneVerified {
		return nil, appErr.NewVerificationRequired(
			"nomor handphone belum diverifikasi",
		)
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		user.FailedLoginAttempt += 1

		if user.FailedLoginAttempt >= 5 {
			user.IsBlocked = true
		}
		_ = s.userRepo.UpdateLoginStatus(user)
		return nil, appErr.NewBadRequest("Password Wrong")
	}

	user.FailedLoginAttempt = 0
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.UpdateLoginStatus(user)

	roleNames := extractRoleNames(user.Role)
	if len(roleNames) == 0 {
		return nil, appErr.NewBadRequest("user has no role assigned")
	}

	token, err := utils.GenerateToken(
		user.ID.String(),
		roleNames,
		nil,
		false,
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:      user.ID.String(),
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Token:   token,
		Balance: user.Balance,
		Roles:   roleNames,
	}, nil
}

// GetProfile (Customer)
func (s *authService) GetProfile(userID string) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.FindProfileById(userID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, appErr.NewNotFound("user not found")
	}

	return &dto.ProfileResponse{
		ID:              user.ID.String(),
		Name:            user.Name,
		Email:           user.Email,
		Phone:           user.Phone,
		Balance:         user.Balance,
		IsPhoneVerified: user.IsPhoneVerified,
		IsEmailVerified: user.IsEmailVerified,
		Language:        user.Language,
		NotifEnabled:    user.NotifEnabled,
		DarkMode:        user.DarkMode,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *authService) RequestDriverOTP(req *dto.DriverLoginRequest) error {
	user, _ := s.userRepo.FindPhoneWithRoles(req.Phone)

	if user == nil {
		return appErr.NewBadRequest("Account Not Found")
	}

	driver, _ := s.driverRepo.FindByUserID(user.ID)
	if driver == nil {
		return appErr.NewForbidden("Anda belum terdaftar sebagai driver")
	}

	if driver.Status != "approved" {
		return appErr.NewForbidden("Driver belum di setujui")
	}

	otp := utils.GenerateOTP()

	if err := s.otpRepo.Save(user.ID, otp); err != nil {
		return err
	}

	utils.SendOTP(req.Phone, otp)
	return nil

}

// Login by phone (Driver)
func (s *authService) LoginDriver(req *dto.DriverLoginRequest) error {
	driver, err := s.driverRepo.FindByPhone(req.Phone)
	if err != nil {
		return appErr.NewBadRequest("Nomor belum terdaftar sebagai driver")
	}

	if driver.Status != "approved" {
		return appErr.NewForbidden("Driver belum di setujui")
	}

	user, err := s.userRepo.FindByPhone(req.Phone)

	if user == nil {
		user = &authModel.User{
			ID:    uuid.New(),
			Name:  driver.FullName,
			Phone: driver.Phone,
		}

		if err := s.userRepo.Create(user); err != nil {
			return err
		}

		role, err := s.roleRepo.FindByName("driver")
		if err != nil {
			return err
		}

		if err := s.userRepo.AsignRole(user.ID, role.ID); err != nil {
			return err
		}

		if err := s.driverRepo.AttachUser(driver.ID, user.ID); err != nil {
			return err
		}
	}

	otp := utils.GenerateOTP()
	_ = s.otpRepo.Save(user.ID, otp)
	utils.SendOTP(req.Phone, otp)

	return nil
}

// Verify Driver OTP
func (s *authService) VerifyDriverOTP(req *dto.VerifyOTPRequest) (*dto.DriverAuthResponse, error) {
	user, _ := s.userRepo.FindPhoneWithRoles(req.Phone)
	if user == nil {
		return nil, appErr.NewBadRequest("Account not found")
	}

	valid, err := s.otpRepo.Verify(user.ID, req.OTP)
	if err != nil || !valid {
		return nil, appErr.NewBadRequest("OTP invalid or expired")
	}

	_ = s.otpRepo.Delete(user.ID)

	driver, _ := s.driverRepo.FindByUserID(user.ID)

	roles := extractRoleNames(user.Role)

	token, err := utils.GenerateToken(
		user.ID.String(),
		roles,
		&driver.ID,
		true,
	)

	if err != nil {
		return nil, err
	}

	return &dto.DriverAuthResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Phone: user.Phone,
		Token: token,
		Roles: roles,
	}, nil

}

// TopUp
func (s *authService) TopUp(userID string, amount int64) (*dto.WalletResponse, error) {
	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, err
	}

	user.Balance += amount
	user.BalanceUpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.WalletResponse{
		Balance:          user.Balance,
		BalanceUpdatedAt: user.BalanceUpdatedAt.Format(time.RFC3339),
	}, nil
}
