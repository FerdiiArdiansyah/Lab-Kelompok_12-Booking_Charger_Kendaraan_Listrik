package main

import (
	"log"

	"booking-service/config"
	deliveryHTTP "booking-service/internal/delivery/http"
	repoPostgres "booking-service/internal/repository/postgres"
	"booking-service/internal/usecase"

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
			log.Printf("Warning: booking_db connection failed, running with in-memory fallback: %v", err)
			gormDB = nil
		} else {
			log.Println("Connected to booking_db successfully via GORM")
		}
	}

	repo := repoPostgres.NewBookingRepository(gormDB)
	uc := usecase.NewBookingUsecase(repo)
	handler := deliveryHTTP.NewBookingHandler(uc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "booking-service", "status": "UP"})
	})

	// === ENDPOINTS BOOKING-SERVICE ===
	// Booking Routes
	e.GET("/bookings", handler.GetAllBookings)
	e.POST("/bookings", handler.CreateBooking)
	e.GET("/bookings/:id", handler.GetBookingByID)
	e.GET("/bookings/user/:user_id", handler.GetBookingsByUserID)

	// Station Availability & Waitlist Routes
	e.GET("/stations/:id/availability", handler.GetAvailability)
	e.GET("/stations/:id/waitlist", handler.GetWaitlist)

	// Operational Routes
	e.POST("/bookings/:id/check-in", handler.CheckIn)
	e.POST("/bookings/:id/cancel", handler.CancelBooking)
	e.POST("/bookings/:id/complete", handler.CompleteBooking)
	e.POST("/bookings/auto-release", handler.AutoRelease)

	log.Printf("booking-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
