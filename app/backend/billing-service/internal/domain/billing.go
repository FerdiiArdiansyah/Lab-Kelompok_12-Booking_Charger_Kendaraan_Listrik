package domain

import (
	"context"
	"time"
)

type Invoice struct {
	ID          string    `gorm:"primaryKey;size:64" json:"id"`
	SessionID   string    `gorm:"index;size:64" json:"session_id"`
	UserID      string    `gorm:"index;size:64" json:"user_id"`
	TariffID    string    `gorm:"index;size:64" json:"tariff_id"`
	ConsumedKwh float64   `json:"consumed_kwh"`
	PricePerKwh float64   `json:"price_per_kwh"`
	ServiceFee  float64   `json:"service_fee"`
	Subtotal    float64   `json:"subtotal"`
	Tax         float64   `json:"tax"`
	Total       float64   `json:"total"`
	Status      string    `gorm:"size:32;default:'UNPAID'" json:"status"` // UNPAID, PAID, VOID, REFUNDED
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Payments    []Payment `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}

type Payment struct {
	ID             string     `gorm:"primaryKey;size:64" json:"id"`
	InvoiceID      string     `gorm:"index;size:64" json:"invoice_id"`
	PaymentMethod  string     `gorm:"size:64" json:"payment_method"` // QRIS, VA_BCA, E_WALLET_GOPAY, CREDIT_CARD
	Amount         float64    `json:"amount"`
	Status         string     `gorm:"size:32;default:'PENDING'" json:"status"` // PENDING, SUCCESS, FAILED, EXPIRED
	TransactionRef string     `gorm:"index;size:128" json:"transaction_ref,omitempty"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}


type BillingRepository interface {
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*Invoice, error)
	GetInvoiceBySessionID(ctx context.Context, sessionID string) (*Invoice, error)
	GetInvoicesByUserID(ctx context.Context, userID string) ([]Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, id string, status string) error
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPaymentByID(ctx context.Context, id string) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, id string, status string, paidAt time.Time) error
	SaveAuditLog(ctx context.Context, entityName, entityID, action string, oldVal, newVal interface{}) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error
}

type BillingUsecase interface {
	GenerateInvoice(ctx context.Context, sessionID, userID, tariffID string, consumedKwh, pricePerKwh, serviceFee float64) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, id string) (*Invoice, error)
	GetInvoiceBySessionID(ctx context.Context, sessionID string) (*Invoice, error)
	GetInvoicesByUserID(ctx context.Context, userID string) ([]Invoice, error)
	ProcessPayment(ctx context.Context, invoiceID, paymentMethod string, amount float64) (*Payment, error)
	GetPaymentByID(ctx context.Context, paymentID string) (*Payment, error)
	ConfirmPayment(ctx context.Context, paymentID, transactionRef string) (*Payment, error)
}
