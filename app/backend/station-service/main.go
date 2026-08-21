package main

import (
	"log"

	"station-service/config"
	deliveryHTTP "station-service/internal/delivery/http"
	repoPostgres "station-service/internal/repository/postgres"
	"station-service/internal/usecase"

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
			log.Printf("Warning: station_db connection failed, running with in-memory fallback: %v", err)
			gormDB = nil
		} else {
			log.Println("Connected to station_db successfully via GORM")
		}
	}

	repo := repoPostgres.NewStationRepository(gormDB)
	uc := usecase.NewStationUsecase(repo)
	handler := deliveryHTTP.NewStationHandler(uc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "station-service", "status": "UP"})
	})

	// === ENDPOINTS STATION-SERVICE ===
	// Station Endpoints
	e.GET("/stations", handler.GetAll)
	e.POST("/stations", handler.Create)
	e.GET("/stations/:id", handler.GetByID)
	e.PUT("/stations/:id", handler.Update)
	e.DELETE("/stations/:id", handler.Delete)

	// Slots Endpoints
	e.GET("/stations/:id/slots", handler.GetSlots)
	e.POST("/stations/:id/slots", handler.AddSlot)
	e.PUT("/stations/:id/slots/:slot_id", handler.UpdateSlot)

	// Tariffs Endpoints
	e.GET("/stations/:id/tariff", handler.GetTariff)
	e.POST("/stations/:id/tariffs", handler.AddTariff)
	e.GET("/tariffs", handler.GetAllTariffs)

	log.Printf("station-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
