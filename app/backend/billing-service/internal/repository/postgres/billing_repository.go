package postgres

import (
	"billing-service/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutboxEventModel struct {
	ID            string    `gorm:"primaryKey;size:64"`
	AggregateType string    `gorm:"index;size:64"`
	AggregateID   string    `gorm:"index;size:64"`
	EventType     string    `gorm:"index;size:64"`
	Payload       string    `gorm:"type:text"`
	Status        string    `gorm:"index;size:32;default:'PENDING'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (OutboxEventModel) TableName() string {
	return "outbox_events"
}

type AuditLogModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	EntityName  string    `gorm:"size:64"`
	EntityID    string    `gorm:"size:64"`
	Action      string    `gorm:"size:64"`
	OldValue    string    `gorm:"type:text"`
	NewValue    string    `gorm:"type:text"`
	PerformedBy string    `gorm:"size:64"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (AuditLogModel) TableName() string {
	return "audit_logs"
}

type billingRepository struct {
	gormDB   *gorm.DB
	invoices map[string]*domain.Invoice
	payments map[string]*domain.Payment
	mu       sync.RWMutex
}

func NewBillingRepository(gormDB *gorm.DB) domain.BillingRepository {
	repo := &billingRepository{
		gormDB:   gormDB,
		invoices: make(map[string]*domain.Invoice),
		payments: make(map[string]*domain.Payment),
	}

	if gormDB != nil {
		// AutoMigrate database tables directly from domain models
		_ = gormDB.AutoMigrate(&domain.Invoice{}, &domain.Payment{}, &OutboxEventModel{}, &AuditLogModel{})
	}

	repo.seedInitialData()
	return repo
}

func (r *billingRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	paidAt := now.Add(-30 * time.Minute)

	sampleInvoices := []domain.Invoice{
		{ID: "inv-001", SessionID: "ses-001", UserID: "usr-001", TariffID: "trf-001", ConsumedKwh: 25.5, PricePerKwh: 2467.0, ServiceFee: 25000.0, Subtotal: 87908.5, Tax: 9669.93, Total: 97578.43, Status: "PAID", CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-002", SessionID: "ses-002", UserID: "usr-002", TariffID: "trf-002", ConsumedKwh: 18.2, PricePerKwh: 2467.0, ServiceFee: 25000.0, Subtotal: 69899.4, Tax: 7688.93, Total: 77588.33, Status: "UNPAID", CreatedAt: now.Add(-1 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)},
		{ID: "inv-003", SessionID: "ses-003", UserID: "usr-005", TariffID: "trf-005", ConsumedKwh: 42.0, PricePerKwh: 2467.0, ServiceFee: 57000.0, Subtotal: 160614.0, Tax: 17667.54, Total: 178281.54, Status: "PAID", CreatedAt: now.Add(-4 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-004", SessionID: "ses-004", UserID: "usr-003", TariffID: "trf-003", ConsumedKwh: 55.0, PricePerKwh: 2467.0, ServiceFee: 57000.0, Subtotal: 192685.0, Tax: 21195.35, Total: 213880.35, Status: "PAID", CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-005", SessionID: "ses-005", UserID: "usr-004", TariffID: "trf-004", ConsumedKwh: 30.1, PricePerKwh: 2467.0, ServiceFee: 25000.0, Subtotal: 99256.7, Tax: 10918.23, Total: 110174.93, Status: "PAID", CreatedAt: now.Add(-6 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-006", SessionID: "ses-006", UserID: "usr-006", TariffID: "trf-006", ConsumedKwh: 22.8, PricePerKwh: 2466.0, ServiceFee: 4000.0, Subtotal: 60224.8, Tax: 6624.72, Total: 66849.52, Status: "PAID", CreatedAt: now.Add(-7 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-007", SessionID: "ses-007", UserID: "usr-007", TariffID: "trf-007", ConsumedKwh: 15.4, PricePerKwh: 3000.0, ServiceFee: 0.0, Subtotal: 46200.0, Tax: 5082.0, Total: 51282.0, Status: "PAID", CreatedAt: now.Add(-8 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-008", SessionID: "ses-008", UserID: "usr-008", TariffID: "trf-008", ConsumedKwh: 60.0, PricePerKwh: 2467.0, ServiceFee: 57000.0, Subtotal: 205020.0, Tax: 22552.2, Total: 227572.2, Status: "PAID", CreatedAt: now.Add(-9 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-009", SessionID: "ses-009", UserID: "usr-009", TariffID: "trf-009", ConsumedKwh: 38.9, PricePerKwh: 2467.0, ServiceFee: 25000.0, Subtotal: 120966.3, Tax: 13306.29, Total: 134272.59, Status: "PAID", CreatedAt: now.Add(-10 * time.Hour), UpdatedAt: paidAt},
		{ID: "inv-010", SessionID: "ses-010", UserID: "usr-010", TariffID: "trf-010", ConsumedKwh: 19.5, PricePerKwh: 2467.0, ServiceFee: 25000.0, Subtotal: 73106.5, Tax: 8041.71, Total: 81148.21, Status: "UNPAID", CreatedAt: now.Add(-11 * time.Hour), UpdatedAt: now.Add(-11 * time.Hour)},
		{ID: "inv-011", SessionID: "ses-011", UserID: "usr-011", TariffID: "trf-011", ConsumedKwh: 48.0, PricePerKwh: 2500.0, ServiceFee: 0.0, Subtotal: 120000.0, Tax: 13200.0, Total: 133200.0, Status: "PAID", CreatedAt: now.Add(-12 * time.Hour), UpdatedAt: paidAt},
	}

	for i := range sampleInvoices {
		inv := sampleInvoices[i]
		r.invoices[inv.ID] = &inv
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&inv, domain.Invoice{ID: inv.ID})
		}
	}

	samplePayments := []domain.Payment{
		{ID: "pay-001", InvoiceID: "inv-001", PaymentMethod: "QRIS", Amount: 97578.43, Status: "SUCCESS", TransactionRef: "TX-QRIS-001", PaidAt: &paidAt, CreatedAt: now.Add(-2 * time.Hour)},
		{ID: "pay-002", InvoiceID: "inv-002", PaymentMethod: "VA_BCA", Amount: 77588.33, Status: "PENDING", TransactionRef: "", PaidAt: nil, CreatedAt: now.Add(-1 * time.Hour)},
		{ID: "pay-003", InvoiceID: "inv-003", PaymentMethod: "E_WALLET_GOPAY", Amount: 178281.54, Status: "SUCCESS", TransactionRef: "TX-GOPAY-003", PaidAt: &paidAt, CreatedAt: now.Add(-4 * time.Hour)},
		{ID: "pay-004", InvoiceID: "inv-004", PaymentMethod: "VA_MANDIRI", Amount: 213880.35, Status: "SUCCESS", TransactionRef: "TX-VA-004", PaidAt: &paidAt, CreatedAt: now.Add(-5 * time.Hour)},
		{ID: "pay-005", InvoiceID: "inv-005", PaymentMethod: "CREDIT_CARD", Amount: 110174.93, Status: "SUCCESS", TransactionRef: "TX-CC-005", PaidAt: &paidAt, CreatedAt: now.Add(-6 * time.Hour)},
		{ID: "pay-006", InvoiceID: "inv-006", PaymentMethod: "QRIS", Amount: 66849.52, Status: "SUCCESS", TransactionRef: "TX-QRIS-006", PaidAt: &paidAt, CreatedAt: now.Add(-7 * time.Hour)},
		{ID: "pay-007", InvoiceID: "inv-007", PaymentMethod: "E_WALLET_SHOPEEPAY", Amount: 51282.0, Status: "SUCCESS", TransactionRef: "TX-SPAY-007", PaidAt: &paidAt, CreatedAt: now.Add(-8 * time.Hour)},
		{ID: "pay-008", InvoiceID: "inv-008", PaymentMethod: "VA_BRI", Amount: 227572.2, Status: "SUCCESS", TransactionRef: "TX-BRI-008", PaidAt: &paidAt, CreatedAt: now.Add(-9 * time.Hour)},
		{ID: "pay-009", InvoiceID: "inv-009", PaymentMethod: "QRIS", Amount: 134272.59, Status: "SUCCESS", TransactionRef: "TX-QRIS-009", PaidAt: &paidAt, CreatedAt: now.Add(-10 * time.Hour)},
		{ID: "pay-010", InvoiceID: "inv-010", PaymentMethod: "VA_BCA", Amount: 81148.21, Status: "PENDING", TransactionRef: "", PaidAt: nil, CreatedAt: now.Add(-11 * time.Hour)},
		{ID: "pay-011", InvoiceID: "inv-011", PaymentMethod: "E_WALLET_OVO", Amount: 133200.0, Status: "SUCCESS", TransactionRef: "TX-OVO-011", PaidAt: &paidAt, CreatedAt: now.Add(-12 * time.Hour)},
	}

	for i := range samplePayments {
		p := samplePayments[i]
		r.payments[p.ID] = &p
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&p, domain.Payment{ID: p.ID})
		}
	}
}

func (r *billingRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(invoice).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoices[invoice.ID] = invoice
	return nil
}

func (r *billingRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	if r.gormDB != nil {
		var inv domain.Invoice
		if err := r.gormDB.WithContext(ctx).Preload("Payments").First(&inv, "id = ?", id).Error; err == nil {
			return &inv, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if inv, ok := r.invoices[id]; ok {
		var pmts []domain.Payment
		for _, p := range r.payments {
			if p.InvoiceID == id {
				pmts = append(pmts, *p)
			}
		}
		inv.Payments = pmts
		return inv, nil
	}
	return nil, errors.New("invoice not found")
}

func (r *billingRepository) GetInvoiceBySessionID(ctx context.Context, sessionID string) (*domain.Invoice, error) {
	if r.gormDB != nil {
		var inv domain.Invoice
		if err := r.gormDB.WithContext(ctx).Preload("Payments").First(&inv, "session_id = ?", sessionID).Error; err == nil {
			return &inv, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, inv := range r.invoices {
		if inv.SessionID == sessionID {
			return inv, nil
		}
	}
	return nil, errors.New("invoice not found for session")
}

func (r *billingRepository) GetInvoicesByUserID(ctx context.Context, userID string) ([]domain.Invoice, error) {
	if r.gormDB != nil {
		var list []domain.Invoice
		if err := r.gormDB.WithContext(ctx).Preload("Payments").Where("user_id = ?", userID).Find(&list).Error; err == nil {
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Invoice
	for _, inv := range r.invoices {
		if inv.UserID == userID {
			list = append(list, *inv)
		}
	}
	return list, nil
}

func (r *billingRepository) UpdateInvoiceStatus(ctx context.Context, id string, status string) error {
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.Invoice{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if inv, ok := r.invoices[id]; ok {
		inv.Status = status
		inv.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("invoice not found")
}

func (r *billingRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	payment.CreatedAt = time.Now()

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(payment).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments[payment.ID] = payment
	return nil
}

func (r *billingRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	if r.gormDB != nil {
		var p domain.Payment
		if err := r.gormDB.WithContext(ctx).First(&p, "id = ?", id).Error; err == nil {
			return &p, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.payments[id]; ok {
		return p, nil
	}
	return nil, errors.New("payment not found")
}

func (r *billingRepository) UpdatePaymentStatus(ctx context.Context, id string, status string, paidAt time.Time) error {
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.Payment{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":  status,
			"paid_at": paidAt,
		})
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.payments[id]; ok {
		p.Status = status
		p.PaidAt = &paidAt
		return nil
	}
	return errors.New("payment not found")
}

func (r *billingRepository) SaveAuditLog(ctx context.Context, entityName, entityID, action string, oldVal, newVal interface{}) error {
	oldBytes, _ := json.Marshal(oldVal)
	newBytes, _ := json.Marshal(newVal)

	if r.gormDB != nil {
		audit := &AuditLogModel{
			EntityName:  entityName,
			EntityID:    entityID,
			Action:      action,
			OldValue:    string(oldBytes),
			NewValue:    string(newBytes),
			PerformedBy: "billing-service",
			CreatedAt:   time.Now(),
		}
		return r.gormDB.WithContext(ctx).Create(audit).Error
	}
	return nil
}

func (r *billingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if r.gormDB != nil {
		model := &OutboxEventModel{
			ID:            uuid.New().String(),
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			EventType:     eventType,
			Payload:       string(payloadBytes),
			Status:        "PENDING",
			CreatedAt:     time.Now(),
		}
		return r.gormDB.WithContext(ctx).Create(model).Error
	}
	return nil
}
