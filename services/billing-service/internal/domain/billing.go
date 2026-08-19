package domain

import (
	"context"
	"time"
)

type Invoice struct {
	ID          string     `json:"id"`
	SessionID   string     `json:"session_id"`
	UserID      string     `json:"user_id"`
	TariffID    string     `json:"tariff_id"`
	ConsumedKwh float64    `json:"consumed_kwh"`
	PricePerKwh float64    `json:"price_per_kwh"`
	Subtotal    float64    `json:"subtotal"`
	Tax         float64    `json:"tax"`
	Total       float64    `json:"total"`
	Status      string     `json:"status"` // UNPAID, PAID, VOID, REFUNDED
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Payments    []Payment  `json:"payments,omitempty"`
}

type Payment struct {
	ID             string     `json:"id"`
	InvoiceID      string     `json:"invoice_id"`
	PaymentMethod  string     `json:"payment_method"` // QRIS, VA_BCA, E_WALLET_GOPAY, CREDIT_CARD
	Amount         float64    `json:"amount"`
	Status         string     `json:"status"` // PENDING, SUCCESS, FAILED, EXPIRED
	TransactionRef string     `json:"transaction_ref,omitempty"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type BillingRepository interface {
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*Invoice, error)
	GetInvoiceBySessionID(ctx context.Context, sessionID string) (*Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, id string, status string) error
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPaymentByID(ctx context.Context, id string) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, id string, status string, paidAt time.Time) error
	SaveAuditLog(ctx context.Context, entityName, entityID, action string, oldVal, newVal interface{}) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error
}

type BillingUsecase interface {
	GenerateInvoice(ctx context.Context, sessionID, userID, tariffID string, consumedKwh, pricePerKwh float64) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, id string) (*Invoice, error)
	ProcessPayment(ctx context.Context, invoiceID, paymentMethod string, amount float64) (*Payment, error)
	ConfirmPayment(ctx context.Context, paymentID, transactionRef string) (*Payment, error)
}
