package main

import (
	"database/sql"
	"log"

	"station-service/config"
	deliveryHTTP "station-service/internal/delivery/http"
	"station-service/internal/repository/postgres"
	"station-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to station_db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: station_db ping failed: %v", err)
	} else {
		log.Println("Connected to station_db successfully")
	}

	repo := postgres.NewStationRepository(db)
	uc := usecase.NewStationUsecase(repo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "station-service", "status": "UP"})
	})

	deliveryHTTP.NewStationHandler(e, uc)

	log.Printf("station-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
