package main

import (
	"log"

	"billing-service/config"
	"billing-service/internal/consumer"
	deliveryHTTP "billing-service/internal/delivery/http"
	repoPostgres "billing-service/internal/repository/postgres"
	"billing-service/internal/usecase"
	"billing-service/pkg/kafka"

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
			log.Printf("Warning: billing_db connection failed, running with in-memory fallback: %v", err)
			gormDB = nil
		} else {
			log.Println("Connected to billing_db successfully via GORM")
		}
	}

	repo := repoPostgres.NewBillingRepository(gormDB)
	uc := usecase.NewBillingUsecase(repo)
	handler := deliveryHTTP.NewBillingHandler(uc)

	// Initialize Kafka Message Broker & Consumer Pipeline
	broker := kafka.GetBroker()
	consumer.NewKafkaBillingConsumer(uc)

	// Publish initial startup system event
	broker.Publish(kafka.TopicBookingConfirmed, "Booking DB", "bkg-sys-001", map[string]interface{}{
		"user_id":    "usr-driver",
		"station_id": "stn-009",
		"slot_id":    "s1",
		"message":    "System initial booking reservation event initialized",
	})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": "billing-service", "status": "UP", "kafka_status": "ONLINE"})
	})

	// === KAFKA MESSAGE BROKER ENDPOINTS ===
	e.GET("/kafka/events", func(c echo.Context) error {
		events := broker.GetEvents()
		return c.JSON(200, map[string]interface{}{
			"broker": "Kafka Message Broker (Port 9092)",
			"topics": []string{
				kafka.TopicBookingConfirmed,
				kafka.TopicBookingExpired,
				kafka.TopicBookingCancelled,
				kafka.TopicChargingStarted,
				kafka.TopicChargingCompleted,
				kafka.TopicPaymentCreated,
				kafka.TopicPaymentCompleted,
				kafka.TopicPaymentFailed,
			},
			"total_events": len(events),
			"events":       events,
		})
	})

	e.POST("/kafka/publish", func(c echo.Context) error {
		var req struct {
			Topic       string                 `json:"topic"`
			Source      string                 `json:"source"`
			AggregateID string                 `json:"aggregate_id"`
			Payload     map[string]interface{} `json:"payload"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, map[string]string{"error": "Invalid Kafka event payload"})
		}
		if req.Topic == "" || req.AggregateID == "" {
			return c.JSON(400, map[string]string{"error": "topic and aggregate_id are required"})
		}
		if req.Source == "" {
			req.Source = "System Event Publisher"
		}
		evt := broker.Publish(req.Topic, req.Source, req.AggregateID, req.Payload)
		return c.JSON(200, map[string]interface{}{
			"message": "Kafka Event successfully published to broker",
			"event":   evt,
		})
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

	log.Printf("billing-service starting on port :%s with Kafka Message Broker...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
