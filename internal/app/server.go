package app

import (
	"be/internal/db"

	authController "be/internal/modules/auth/controller"
	authRepository "be/internal/modules/auth/repository"
	authRoutes "be/internal/modules/auth/routes"
	authService "be/internal/modules/auth/service"

	roleRepository "be/internal/modules/role/repository"

	driverController "be/internal/modules/driver/controller"
	driverRepository "be/internal/modules/driver/repository"
	driverRoutes "be/internal/modules/driver/routes"
	driverService "be/internal/modules/driver/service"

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
	// === Init Dependencies ===
	userRepo := authRepository.NewUserRepository(db.DB)
	roleRepo := roleRepository.NewRoleRepository(db.DB)
	authService := authService.NewAuthService(userRepo, roleRepo, "SUPERSECRETKEY")
	authController := authController.NewAuthController(authService)

	// === Init Driver Vehicle Dependencies ===
	vehicleRepo := driverRepository.NewVehicleRepository(db.DB)
	vehicleService := driverService.NewVehicleService(vehicleRepo)
	vehicleController := driverController.NewVehicleController(vehicleService)

	// === Init Driver Routes Driver ===
	driverRepo := driverRepository.NewDriverRepository(db.DB)
	driverService := driverService.NewDriverService(driverRepo)
	driverController := driverController.NewDriverController(driverService)

	// === Main API Group ===
	api := a.Engine.Group("/api")

	// === Auth Routes ===
	authRoutes := authRoutes.NewAuthRoutes(authController)
	authRoutes.RegisterRoutes(api)

	// === Vehicle Routes ===
	vehicleRoutes := driverRoutes.NewVehicleRoutes(vehicleController)
	vehicleRoutes.RegisterRoutes(api)

	// === Driver Routes ===
	driverRoutes := driverRoutes.NewDriverRoutes(driverController)
	driverRoutes.RegisterRoutes(api)

	// === DEBUG: LIST ROUTES ===
	fmt.Println("== Registered Routes ==")
	for _, r := range a.Engine.Routes() {
		fmt.Println(r.Method, r.Path)
	}
}

func (a *App) Run() {
	addr := fmt.Sprintf(":%s", a.Port)
	log.Println("Starting server on", addr)
	a.Engine.Run(addr)
}
