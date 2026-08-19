package main

import (
	"billing-service/config"
	deliveryHTTP "billing-service/internal/delivery/http"
	"billing-service/internal/repository/postgres"
	"billing-service/internal/usecase"
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
		log.Fatalf("Failed to connect to billing_db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: billing_db ping failed: %v", err)
	} else {
		log.Println("Connected to billing_db successfully")
	}

	repo := postgres.NewBillingRepository(db)
	uc := usecase.NewBillingUsecase(repo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "billing-service", "status": "UP"})
	})

	deliveryHTTP.NewBillingHandler(e, uc)

	log.Printf("billing-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
