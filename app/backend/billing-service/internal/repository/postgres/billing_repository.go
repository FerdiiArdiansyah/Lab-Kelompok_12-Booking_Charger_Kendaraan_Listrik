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

	if r.gormDB != nil {
		r.gormDB.Exec("DELETE FROM payments")
		r.gormDB.Exec("DELETE FROM invoices")
	}
	r.invoices = make(map[string]*domain.Invoice)
	r.payments = make(map[string]*domain.Payment)
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
