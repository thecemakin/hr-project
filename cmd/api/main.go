package main

// @title HR Backend API
// @version 1.0
// @description This is a HR management system API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /

import (
	"log"
	"time"

	"github.com/thecemakin/hr-project/internal/platform/config"
	"github.com/thecemakin/hr-project/internal/platform/db"
	server "github.com/thecemakin/hr-project/internal/platform/http"
	"github.com/thecemakin/hr-project/internal/platform/auth"

	authHandler "github.com/thecemakin/hr-project/internal/modules/auth/handler"
	authRepo "github.com/thecemakin/hr-project/internal/modules/auth/repository"
	authSvc "github.com/thecemakin/hr-project/internal/modules/auth/service"

	corehrHandler "github.com/thecemakin/hr-project/internal/modules/corehr/handler"
	corehrRepo "github.com/thecemakin/hr-project/internal/modules/corehr/repository"
	corehrSvc "github.com/thecemakin/hr-project/internal/modules/corehr/service"

	leaveHandler "github.com/thecemakin/hr-project/internal/modules/leave/handler"
	leaveRepo "github.com/thecemakin/hr-project/internal/modules/leave/repository"
	leaveSvc "github.com/thecemakin/hr-project/internal/modules/leave/service"
)

func main() {
	log.Println("HR Backend API starting...")

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Setup database connection (ignoring currently to just test boot)
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v. Cannot boot up modules.", err)
	}
	log.Println("Database connection established:", database.Name())

	// 3. Setup Auth Platform
	ttl, err := time.ParseDuration(cfg.JWTAccessTokenTTL)
	if err != nil {
		log.Printf("Invalid JWT TTL: %v, falling back to 15m", err)
		ttl = 15 * time.Minute
	}
	tp := auth.NewTokenProvider(cfg.JWTSecret, ttl)

	// 4. Setup HTTP server and routing
	srv := server.NewServer()

	// 5. Initialize Auth Module
	userRepository := authRepo.NewUserRepository(database)
	authService := authSvc.NewAuthService(userRepository, tp)
	authHdl := authHandler.NewAuthHandler(authService)

	// Mount Auth Routes
	authHandler.SetupRoutesFiber(srv.App, authHdl, tp)

	// 6. Initialize CoreHR Module
	employeeRepo := corehrRepo.NewEmployeeRepository(database)
	departmentRepo := corehrRepo.NewDepartmentRepository(database)
	positionRepo := corehrRepo.NewPositionRepository(database)
	assetRepo := corehrRepo.NewAssetRepository(database)
	assetAssignmentRepo := corehrRepo.NewAssetAssignmentRepository(database)

	employeeSvc := corehrSvc.NewEmployeeService(employeeRepo)
	departmentSvc := corehrSvc.NewDepartmentService(departmentRepo)
	positionSvc := corehrSvc.NewPositionService(positionRepo)
	assetSvc := corehrSvc.NewAssetService(assetRepo)
	assetAssignmentSvc := corehrSvc.NewAssetAssignmentService(assetAssignmentRepo, assetRepo, employeeRepo)

	employeeHdl := corehrHandler.NewEmployeeHandler(employeeSvc)
	departmentHdl := corehrHandler.NewDepartmentHandler(departmentSvc)
	positionHdl := corehrHandler.NewPositionHandler(positionSvc)
	assetHdl := corehrHandler.NewAssetHandler(assetSvc)
	assetAssignmentHdl := corehrHandler.NewAssetAssignmentHandler(assetAssignmentSvc)

	// Mount Core HR Routes
	corehrHandler.SetupRoutesFiber(srv.App, employeeHdl, departmentHdl, positionHdl, assetHdl, assetAssignmentHdl, tp)

	// 5. Initialize Leave Module
	leaveRepository := leaveRepo.NewSQLRepository(database)
	leaveService := leaveSvc.NewLeaveService(leaveRepository, employeeRepo)
	
	leaveHandler.SetupRoutesFiber(srv.App, leaveService, tp)

	log.Printf("Listening and serving HTTP on :%s", cfg.HTTPPort)
	if err := srv.App.Listen(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
