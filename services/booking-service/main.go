package main

import (
	"booking-service/config"
	deliveryHTTP "booking-service/internal/delivery/http"
	"booking-service/internal/repository/postgres"
	"booking-service/internal/usecase"
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to booking_db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: booking_db ping failed: %v", err)
	} else {
		log.Println("Connected to booking_db successfully")
	}

	repo := postgres.NewBookingRepository(db)
	uc := usecase.NewBookingUsecase(repo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "booking-service", "status": "UP"})
	})

	deliveryHTTP.NewBookingHandler(e, uc)

	log.Printf("booking-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
