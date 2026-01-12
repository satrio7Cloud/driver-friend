package app

import (
	"be/internal/db"

	authControllerPkg "be/internal/modules/auth/controller"
	authRepository "be/internal/modules/auth/repository"
	authRoutesPkg "be/internal/modules/auth/routes"
	authServicePkg "be/internal/modules/auth/service"

	roleRepository "be/internal/modules/role/repository"

	driverControllerPkg "be/internal/modules/driver/controller"
	driverRepository "be/internal/modules/driver/repository"
	driverRoutesPkg "be/internal/modules/driver/routes"
	driverServicePkg "be/internal/modules/driver/service"
	otpRepository "be/internal/modules/otp/repository"

	vehicleControllerPkg "be/internal/modules/vehicle/controller"
	vehicleRepository "be/internal/modules/vehicle/repository"
	vehicleRoutesPkg "be/internal/modules/vehicle/routes"
	vehicleServicePkg "be/internal/modules/vehicle/service"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
	Port   string
}

func NewApp(port string) *App {
	engine := gin.Default()

	app := &App{
		Engine: engine,
		Port:   port,
	}

	app.RegisterRoutes()
	return app
}

func (a *App) RegisterRoutes() {
	// ================= AUTH =================
	userRepo := authRepository.NewUserRepository(db.DB)
	roleRepo := roleRepository.NewRoleRepository(db.DB)
	otpRepo := otpRepository.NewOTPRepository()
	driverRepo := driverRepository.NewDriverRepository(db.DB)
	authSvc := authServicePkg.NewAuthService(userRepo, roleRepo, driverRepo, otpRepo, "SUPERSECRETKEY")
	authCtrl := authControllerPkg.NewAuthController(authSvc)

	// ================= DRIVER =================
	// driverRepo := driverRepository.NewDriverRepository(db.DB)
	driverSvc := driverServicePkg.NewDriverService(driverRepo)

	driverCtrl := driverControllerPkg.NewDriverController(driverSvc)
	adminDriverCtrl := driverControllerPkg.NewAdminDriverController(driverSvc)

	// ================= VEHICLE =================
	vehicleRepo := vehicleRepository.NewVehicleRepository(db.DB)
	vehicleSvc := vehicleServicePkg.NewVehicleService(vehicleRepo, driverRepo)
	vehicleCtrl := vehicleControllerPkg.NewVehicleController(vehicleSvc)

	// ================= ROUTER =================
	api := a.Engine.Group("/api")

	// Auth
	authRoutes := authRoutesPkg.NewAuthRoutes(authCtrl)
	authRoutes.RegisterRoutes(api)

	// Driver (User)
	driverRoutes := driverRoutesPkg.NewDriverRoutes(driverCtrl)
	driverRoutes.RegisterRoutes(api)

	// Driver (Admin)
	adminDriverRoutes := driverRoutesPkg.NewAdminDriverRoutes(adminDriverCtrl)
	adminDriverRoutes.RegisterRoutes(api)

	// Vehicle
	vehicleRoutes := vehicleRoutesPkg.NewVehicleRoutes(vehicleCtrl)
	vehicleRoutes.RegisterRoutes(api)

	// DEBUG
	fmt.Println("== Registered Routes ==")
	for _, r := range a.Engine.Routes() {
		fmt.Println(r.Method, r.Path)
	}
}

func (a *App) Run() {
	addr := fmt.Sprintf(":%s", a.Port)
	log.Println("Server running on", addr)
	a.Engine.Run(addr)
}
