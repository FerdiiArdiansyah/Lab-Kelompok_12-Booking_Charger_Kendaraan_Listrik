package main

import (
	"database/sql"
	"log"
	"session-service/config"
	deliveryHTTP "session-service/internal/delivery/http"
	"session-service/internal/repository/postgres"
	"session-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to session_db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: session_db ping failed: %v", err)
	} else {
		log.Println("Connected to session_db successfully")
	}

	repo := postgres.NewSessionRepository(db)
	uc := usecase.NewSessionUsecase(repo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "session-service", "status": "UP"})
	})

	deliveryHTTP.NewSessionHandler(e, uc)

	log.Printf("session-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
