package usecase

import (
	"billing-service/internal/domain"
	"context"
	"errors"
	"testing"
	"time"
)

type MockBillingRepository struct {
	invoices     map[string]*domain.Invoice
	payments     map[string]*domain.Payment
	auditLogs    []map[string]interface{}
	outboxEvents []map[string]interface{}
}

func NewMockBillingRepository() *MockBillingRepository {
	return &MockBillingRepository{
		invoices:     make(map[string]*domain.Invoice),
		payments:     make(map[string]*domain.Payment),
		auditLogs:    make([]map[string]interface{}, 0),
		outboxEvents: make([]map[string]interface{}, 0),
	}
}

func (m *MockBillingRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	m.invoices[invoice.ID] = invoice
	return nil
}

func (m *MockBillingRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	if inv, ok := m.invoices[id]; ok {
		return inv, nil
	}
	return nil, errors.New("invoice not found")
}

func (m *MockBillingRepository) GetInvoiceBySessionID(ctx context.Context, sessionID string) (*domain.Invoice, error) {
	for _, inv := range m.invoices {
		if inv.SessionID == sessionID {
			return inv, nil
		}
	}
	return nil, errors.New("invoice not found for session")
}

func (m *MockBillingRepository) GetInvoicesByUserID(ctx context.Context, userID string) ([]domain.Invoice, error) {
	var list []domain.Invoice
	for _, inv := range m.invoices {
		if inv.UserID == userID {
			list = append(list, *inv)
		}
	}
	return list, nil
}

func (m *MockBillingRepository) UpdateInvoiceStatus(ctx context.Context, id string, status string) error {
	if inv, ok := m.invoices[id]; ok {
		inv.Status = status
		return nil
	}
	return errors.New("invoice not found")
}

func (m *MockBillingRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	m.payments[payment.ID] = payment
	return nil
}

func (m *MockBillingRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	if p, ok := m.payments[id]; ok {
		return p, nil
	}
	return nil, errors.New("payment not found")
}

func (m *MockBillingRepository) UpdatePaymentStatus(ctx context.Context, id string, status string, paidAt time.Time) error {
	if p, ok := m.payments[id]; ok {
		p.Status = status
		p.PaidAt = &paidAt
		return nil
	}
	return errors.New("payment not found")
}

func (m *MockBillingRepository) SaveAuditLog(ctx context.Context, entityName, entityID, action string, oldVal, newVal interface{}) error {
	m.auditLogs = append(m.auditLogs, map[string]interface{}{
		"entity_name": entityName,
		"entity_id":   entityID,
		"action":      action,
		"old_val":     oldVal,
		"new_val":     newVal,
	})
	return nil
}

func (m *MockBillingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	m.outboxEvents = append(m.outboxEvents, map[string]interface{}{
		"aggregate_type": aggregateType,
		"aggregate_id":   aggregateID,
		"event_type":     eventType,
		"payload":        payload,
	})
	return nil
}

// Unit Tests for TICKET-05: ESDM Tariff & PPN 11% Calculation Engine

func TestGenerateInvoice_FastChargingESDMTariff_CalculationAccuracy(t *testing.T) {
	repo := NewMockBillingRepository()
	uc := NewBillingUsecase(repo)
	ctx := context.Background()

	// Fast charging: 20 kWh, Rp 2467 per kWh, Service fee Rp 25.000
	consumedKwh := 20.0
	pricePerKwh := 2467.0
	serviceFee := 25000.0

	inv, err := uc.GenerateInvoice(ctx, "ses-fast-1", "usr-1", "trf-fast", consumedKwh, pricePerKwh, serviceFee)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedSubtotal := 20.0*2467.0 + 25000.0 // 74,340
	expectedTax := 74340.0 * 0.11             // 8,177.4
	expectedTotal := 74340.0 + 8177.4         // 82,517.4

	if inv.Subtotal != expectedSubtotal {
		t.Errorf("Expected Subtotal %f, got %f", expectedSubtotal, inv.Subtotal)
	}

	if inv.Tax != expectedTax {
		t.Errorf("Expected Tax %f, got %f", expectedTax, inv.Tax)
	}

	if inv.Total != expectedTotal {
		t.Errorf("Expected Total %f, got %f", expectedTotal, inv.Total)
	}

	if inv.Status != "UNPAID" {
		t.Errorf("Expected Status UNPAID, got %s", inv.Status)
	}
}

func TestGenerateInvoice_UltraFastChargingESDMTariff_CalculationAccuracy(t *testing.T) {
	repo := NewMockBillingRepository()
	uc := NewBillingUsecase(repo)
	ctx := context.Background()

	// Ultra-Fast charging: 50 kWh, Rp 2467 per kWh, Service fee Rp 57.000
	consumedKwh := 50.0
	pricePerKwh := 2467.0
	serviceFee := 57000.0

	inv, err := uc.GenerateInvoice(ctx, "ses-ultrafast-1", "usr-2", "trf-ultrafast", consumedKwh, pricePerKwh, serviceFee)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedSubtotal := 50.0*2467.0 + 57000.0 // 180,350
	expectedTax := 180350.0 * 0.11            // 19,838.5
	expectedTotal := 180350.0 + 19838.5       // 200,188.5

	if inv.Subtotal != expectedSubtotal {
		t.Errorf("Expected Subtotal %f, got %f", expectedSubtotal, inv.Subtotal)
	}

	if inv.Tax != expectedTax {
		t.Errorf("Expected Tax %f, got %f", expectedTax, inv.Tax)
	}

	if inv.Total != expectedTotal {
		t.Errorf("Expected Total %f, got %f", expectedTotal, inv.Total)
	}
}

func TestProcessAndConfirmPayment_UpdatesInvoiceToPaid(t *testing.T) {
	repo := NewMockBillingRepository()
	uc := NewBillingUsecase(repo)
	ctx := context.Background()

	inv, _ := uc.GenerateInvoice(ctx, "ses-pay-1", "usr-3", "trf-1", 10.0, 2467.0, 25000.0)

	payment, err := uc.ProcessPayment(ctx, inv.ID, "QRIS", inv.Total)
	if err != nil {
		t.Fatalf("Expected no error processing payment, got: %v", err)
	}

	if payment.Status != "PENDING" {
		t.Errorf("Expected payment status PENDING, got %s", payment.Status)
	}

	confirmedPay, err := uc.ConfirmPayment(ctx, payment.ID, "TX-REF-9999")
	if err != nil {
		t.Fatalf("Expected no error confirming payment, got: %v", err)
	}

	if confirmedPay.Status != "SUCCESS" {
		t.Errorf("Expected confirmed payment status SUCCESS, got %s", confirmedPay.Status)
	}

	updatedInv, _ := repo.GetInvoiceByID(ctx, inv.ID)
	if updatedInv.Status != "PAID" {
		t.Errorf("Expected invoice status to be updated to PAID, got %s", updatedInv.Status)
	}

	if len(repo.auditLogs) != 1 {
		t.Errorf("Expected 1 audit log created, got %d", len(repo.auditLogs))
	}
}

// Tests for TICKET-06: Idempotent Payment Processing & Webhook Handler

func TestConfirmPayment_IdempotentDuplicateWebhook_SuccessWithoutDoublePosting(t *testing.T) {
	repo := NewMockBillingRepository()
	uc := NewBillingUsecase(repo)
	ctx := context.Background()

	inv, _ := uc.GenerateInvoice(ctx, "ses-idem-1", "usr-idem", "trf-1", 15.0, 2467.0, 25000.0)
	payment, _ := uc.ProcessPayment(ctx, inv.ID, "VA_BCA", inv.Total)

	// First webhook confirmation
	pay1, err1 := uc.ConfirmPayment(ctx, payment.ID, "TX-IDEM-001")
	if err1 != nil {
		t.Fatalf("First payment confirmation failed: %v", err1)
	}

	initialEventsCount := len(repo.outboxEvents)

	// Second duplicate webhook confirmation (Idempotency test)
	pay2, err2 := uc.ConfirmPayment(ctx, payment.ID, "TX-IDEM-001")
	if err2 != nil {
		t.Fatalf("Second (duplicate) payment confirmation failed: %v", err2)
	}

	if pay1.Status != pay2.Status || pay2.Status != "SUCCESS" {
		t.Errorf("Expected both payments to have SUCCESS status")
	}

	// Ensure no duplicate outbox event was generated for the second confirmation
	if len(repo.outboxEvents) != initialEventsCount {
		t.Errorf("Expected outbox event count to stay %d, but got %d (duplicate outbox event detected)", initialEventsCount, len(repo.outboxEvents))
	}
}

func TestProcessPayment_AlreadyPaidInvoice_ReturnsError(t *testing.T) {
	repo := NewMockBillingRepository()
	uc := NewBillingUsecase(repo)
	ctx := context.Background()

	inv, _ := uc.GenerateInvoice(ctx, "ses-paid-1", "usr-paid", "trf-1", 10.0, 2467.0, 25000.0)
	payment, _ := uc.ProcessPayment(ctx, inv.ID, "QRIS", inv.Total)
	_, _ = uc.ConfirmPayment(ctx, payment.ID, "TX-CONFIRMED")

	// Attempt to create a second payment for already paid invoice
	_, err := uc.ProcessPayment(ctx, inv.ID, "E_WALLET_GOPAY", inv.Total)
	if err == nil {
		t.Errorf("Expected error when processing payment for an already PAID invoice, got nil")
	}
}

