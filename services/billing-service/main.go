package main

import (
	"database/sql"
	"log"

	"billing-service/config"
	deliveryHTTP "billing-service/internal/delivery/http"
	"billing-service/internal/repository/postgres"
	"billing-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Printf("Warning: billing_db connection string error: %v", err)
	} else {
		defer db.Close()
		if err := db.Ping(); err != nil {
			log.Printf("Warning: billing_db ping failed, running with in-memory fallback: %v", err)
		} else {
			log.Println("Connected to billing_db successfully")
		}
	}

	repo := postgres.NewBillingRepository(db)
	uc := usecase.NewBillingUsecase(repo)
	handler := deliveryHTTP.NewBillingHandler(uc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "billing-service", "status": "UP"})
	})

	// === ENDPOINTS BILLING-SERVICE ===
	// Invoice Routes
	e.POST("/invoices", handler.GenerateInvoice)
	e.POST("/invoices/generate", handler.GenerateInvoice)
	e.GET("/invoices/:id", handler.GetInvoiceByID)
	e.GET("/invoices/session/:session_id", handler.GetInvoiceBySessionID)
	e.GET("/invoices/user/:user_id", handler.GetInvoicesByUserID)

	// Payment Routes
	e.POST("/payments", handler.ProcessPayment)
	e.GET("/payments/:id", handler.GetPaymentByID)
	e.POST("/payments/:id/confirm", handler.ConfirmPayment)

	log.Printf("billing-service starting on port :%s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
