package main

import (
	"database/sql"
	"log"

	"user-service/config"
	deliveryHTTP "user-service/internal/delivery/http"
	"user-service/internal/repository/postgres"
	"user-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Printf("Warning: user_db connection string error: %v", err)
	} else {
		defer db.Close()
		if err := db.Ping(); err != nil {
			log.Printf("Warning: user_db ping failed, running with in-memory fallback: %v", err)
		} else {
			log.Println("Connected to user_db successfully")
		}
	}

	repo := postgres.NewUserRepository(db)
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
