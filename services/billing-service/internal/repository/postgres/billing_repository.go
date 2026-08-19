package postgres

import (
	"billing-service/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type billingRepository struct {
	db *sql.DB
}

func NewBillingRepository(db *sql.DB) domain.BillingRepository {
	return &billingRepository{db: db}
}

func (r *billingRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	query := `INSERT INTO invoices (id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.SessionID, invoice.UserID, invoice.TariffID,
		invoice.ConsumedKwh, invoice.PricePerKwh, invoice.Subtotal, invoice.Tax, invoice.Total, invoice.Status)
	return err
}

func (r *billingRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	query := `SELECT id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status, created_at, updated_at
	          FROM invoices WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var inv domain.Invoice
	if err := row.Scan(&inv.ID, &inv.SessionID, &inv.UserID, &inv.TariffID, &inv.ConsumedKwh, &inv.PricePerKwh, &inv.Subtotal, &inv.Tax, &inv.Total, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	payments, _ := r.getPaymentsByInvoiceID(ctx, id)
	inv.Payments = payments

	return &inv, nil
}

func (r *billingRepository) GetInvoiceBySessionID(ctx context.Context, sessionID string) (*domain.Invoice, error) {
	query := `SELECT id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status, created_at, updated_at
	          FROM invoices WHERE session_id = $1`
	row := r.db.QueryRowContext(ctx, query, sessionID)

	var inv domain.Invoice
	if err := row.Scan(&inv.ID, &inv.SessionID, &inv.UserID, &inv.TariffID, &inv.ConsumedKwh, &inv.PricePerKwh, &inv.Subtotal, &inv.Tax, &inv.Total, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invoice not found for session")
		}
		return nil, err
	}
	return &inv, nil
}

func (r *billingRepository) UpdateInvoiceStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("invoice not found")
	}
	return nil
}

func (r *billingRepository) getPaymentsByInvoiceID(ctx context.Context, invoiceID string) ([]domain.Payment, error) {
	query := `SELECT id, invoice_id, payment_method, amount, status, transaction_ref, paid_at, created_at
	          FROM payments WHERE invoice_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Payment
	for rows.Next() {
		var p domain.Payment
		var txRef sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.InvoiceID, &p.PaymentMethod, &p.Amount, &p.Status, &txRef, &paidAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.TransactionRef = txRef.String
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *billingRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	query := `INSERT INTO payments (id, invoice_id, payment_method, amount, status)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, payment.ID, payment.InvoiceID, payment.PaymentMethod, payment.Amount, payment.Status)
	return err
}

func (r *billingRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `SELECT id, invoice_id, payment_method, amount, status, transaction_ref, paid_at, created_at
	          FROM payments WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var p domain.Payment
	var txRef sql.NullString
	var paidAt sql.NullTime
	if err := row.Scan(&p.ID, &p.InvoiceID, &p.PaymentMethod, &p.Amount, &p.Status, &txRef, &paidAt, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	p.TransactionRef = txRef.String
	if paidAt.Valid {
		p.PaidAt = &paidAt.Time
	}
	return &p, nil
}

func (r *billingRepository) UpdatePaymentStatus(ctx context.Context, id string, status string, paidAt time.Time) error {
	query := `UPDATE payments SET status = $1, paid_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, paidAt, id)
	return err
}

func (r *billingRepository) SaveAuditLog(ctx context.Context, entityName, entityID, action string, oldVal, newVal interface{}) error {
	oldBytes, _ := json.Marshal(oldVal)
	newBytes, _ := json.Marshal(newVal)
	query := `INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value, performed_by)
	          VALUES ($1, $2, $3, $4, $5, 'billing-service')`
	_, err := r.db.ExecContext(ctx, query, entityName, entityID, action, oldBytes, newBytes)
	return err
}

func (r *billingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	query := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status)
	          VALUES ($1, $2, $3, $4, $5, 'PENDING')`
	_, err = r.db.ExecContext(ctx, query, uuid.New().String(), aggregateType, aggregateID, eventType, payloadBytes)
	return err
}
