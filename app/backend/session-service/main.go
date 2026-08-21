package main

import (
	"log"

	"session-service/config"
	deliveryHTTP "session-service/internal/delivery/http"
	repoPostgres "session-service/internal/repository/postgres"
	"session-service/internal/usecase"

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
			log.Printf("Warning: session_db connection failed, running with in-memory fallback: %v", err)
			gormDB = nil
		} else {
			log.Println("Connected to session_db successfully via GORM")
		}
	}

	repo := repoPostgres.NewSessionRepository(gormDB)
	uc := usecase.NewSessionUsecase(repo)
	handler := deliveryHTTP.NewSessionHandler(uc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "session-service", "status": "UP"})
	})

	// === ENDPOINTS SESSION-SERVICE ===
	// Charging Session Routes
	e.POST("/sessions/start", handler.StartSession)
	e.GET("/sessions/:id", handler.GetSessionByID)
	e.GET("/sessions/booking/:booking_id", handler.GetSessionByBookingID)
	e.GET("/sessions/user/:user_id", handler.GetSessionsByUserID)

	// Meter & Telemetry Routes
	e.POST("/sessions/:id/meter", handler.RecordMeter)

	// Finish Session Route
	e.POST("/sessions/:id/finish", handler.FinishSession)

	log.Printf("session-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
