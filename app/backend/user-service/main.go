package main

import (
	"log"

	"user-service/config"
	deliveryHTTP "user-service/internal/delivery/http"
	repoPostgres "user-service/internal/repository/postgres"
	"user-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	var gormDB *gorm.DB
	var err error

	if cfg.DSN() != "" {
		gormDB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
		if err != nil {
			log.Printf("Warning: user_db connection failed, running with in-memory fallback: %v", err)
			gormDB = nil
		} else {
			log.Println("Connected to user_db successfully via GORM")
		}
	}

	repo := repoPostgres.NewUserRepository(gormDB)
	uc := usecase.NewUserUsecase(repo, cfg.JWTSecret)
	handler := deliveryHTTP.NewUserHandler(uc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "user-service", "status": "UP"})
	})

	// === ENDPOINTS USER-SERVICE ===
	// Auth Routes
	e.POST("/auth/register", handler.Register)
	e.POST("/auth/login", handler.Login)

	// User Profile Routes
	e.GET("/users/me", handler.GetProfile)
	e.PUT("/users/me", handler.UpdateProfile)
	e.GET("/users", handler.GetAllUsers)
	e.GET("/users/:id", handler.GetUserByID)

	// Vehicle Management Routes
	e.GET("/users/me/vehicles", handler.GetVehicles)
	e.POST("/users/me/vehicles", handler.AddVehicle)
	e.DELETE("/users/me/vehicles/:vehicle_id", handler.DeleteVehicle)

	log.Printf("user-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
