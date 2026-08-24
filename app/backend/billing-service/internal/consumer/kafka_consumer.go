package consumer

import (
	"context"
	"log"

	"billing-service/internal/domain"
	"billing-service/pkg/kafka"
)

type KafkaBillingConsumer struct {
	billingUsecase domain.BillingUsecase
	broker         *kafka.EventBroker
}

func NewKafkaBillingConsumer(uc domain.BillingUsecase) *KafkaBillingConsumer {
	c := &KafkaBillingConsumer{
		billingUsecase: uc,
		broker:         kafka.GetBroker(),
	}
	c.registerSubscriptions()
	return c
}

func (c *KafkaBillingConsumer) registerSubscriptions() {
	// 1. Listen for BookingConfirmed Event (Booking DB -> Kafka -> Billing Service)
	c.broker.Subscribe(kafka.TopicBookingConfirmed, func(evt kafka.Event) {
		log.Printf("[KAFKA CONSUMER] Processing Topic [%s] for Booking ID: %s", evt.Topic, evt.AggregateID)
		userID, _ := evt.Payload["user_id"].(string)
		if userID == "" {
			userID = "usr-driver"
		}
		// Pre-calculate estimated booking reservation fee
		log.Printf("[KAFKA BILLING] Step 1: Calculate Charges for Booking Reservation ID %s", evt.AggregateID)
	})

	// 2. Listen for ChargingCompleted Event (Session DB -> Kafka -> Billing Service)
	c.broker.Subscribe(kafka.TopicChargingCompleted, func(evt kafka.Event) {
		log.Printf("[KAFKA CONSUMER] Processing Topic [%s] for Session ID: %s", evt.Topic, evt.AggregateID)

		sessionID := evt.AggregateID
		userID, _ := evt.Payload["user_id"].(string)
		if userID == "" {
			userID = "usr-driver"
		}
		consumedKwh, _ := evt.Payload["consumed_kwh"].(float64)
		if consumedKwh <= 0 {
			consumedKwh = 15.5
		}

		// Step 1: Calculate Charges & Step 2: Generate Invoice
		log.Printf("[KAFKA BILLING] Executing Pipeline: 1. Calculate Charges -> 2. Generate Invoice for Session %s", sessionID)
		inv, err := c.billingUsecase.GenerateInvoice(
			context.Background(),
			sessionID,
			userID,
			"trf-pln-standard",
			consumedKwh,
			2467.0, // PLN rate per kWh
			5000.0, // Admin/Service Fee
		)

		if err != nil {
			log.Printf("[KAFKA BILLING ERROR] Failed to generate invoice from ChargingCompleted event: %v", err)
			c.broker.Publish(kafka.TopicPaymentFailed, "Billing DB", sessionID, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Step 3: Process Payment (Create Pending Payment QRIS/VA)
		log.Printf("[KAFKA BILLING] Step 3: Process Payment for Invoice ID %s (Total: Rp %.2f)", inv.ID, inv.Total)
		pmt, err := c.billingUsecase.ProcessPayment(context.Background(), inv.ID, "QRIS", inv.Total)
		if err != nil {
			log.Printf("[KAFKA BILLING ERROR] Failed to process payment: %v", err)
			return
		}

		// Step 4: Publish PaymentCreated Event to Kafka
		c.broker.Publish(kafka.TopicPaymentCreated, "Billing DB", inv.ID, map[string]interface{}{
			"invoice_id": inv.ID,
			"payment_id": pmt.ID,
			"amount":     pmt.Amount,
			"status":     pmt.Status,
		})
	})

	// 3. Listen for PaymentCompleted Event (Billing DB -> Kafka)
	c.broker.Subscribe(kafka.TopicPaymentCompleted, func(evt kafka.Event) {
		log.Printf("[KAFKA CONSUMER] Step 4: Payment Status -> COMPLETED for Invoice ID: %s", evt.AggregateID)
	})

	// 4. Listen for PaymentFailed Event
	c.broker.Subscribe(kafka.TopicPaymentFailed, func(evt kafka.Event) {
		log.Printf("[KAFKA CONSUMER] Step 4: Payment Status -> FAILED for Invoice ID: %s", evt.AggregateID)
	})
}
