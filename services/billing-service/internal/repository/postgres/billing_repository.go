package postgres

import (
	"billing-service/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type billingRepository struct {
	db       *sql.DB
	invoices map[string]*domain.Invoice
	payments map[string]*domain.Payment
	mu       sync.RWMutex
}

func NewBillingRepository(db *sql.DB) domain.BillingRepository {
	return &billingRepository{
		db:       db,
		invoices: make(map[string]*domain.Invoice),
		payments: make(map[string]*domain.Payment),
	}
}

func (r *billingRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO invoices (id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status)
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err := r.db.ExecContext(ctx, query,
			invoice.ID, invoice.SessionID, invoice.UserID, invoice.TariffID,
			invoice.ConsumedKwh, invoice.PricePerKwh, invoice.Subtotal, invoice.Tax, invoice.Total, invoice.Status)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoices[invoice.ID] = invoice
	return nil
}

func (r *billingRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	if r.db != nil {
		query := `SELECT id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status, created_at, updated_at
		          FROM invoices WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)

		var inv domain.Invoice
		if err := row.Scan(&inv.ID, &inv.SessionID, &inv.UserID, &inv.TariffID, &inv.ConsumedKwh, &inv.PricePerKwh, &inv.Subtotal, &inv.Tax, &inv.Total, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt); err == nil {
			payments, _ := r.getPaymentsByInvoiceID(ctx, id)
			inv.Payments = payments
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
	if r.db != nil {
		query := `SELECT id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status, created_at, updated_at
		          FROM invoices WHERE session_id = $1`
		row := r.db.QueryRowContext(ctx, query, sessionID)

		var inv domain.Invoice
		if err := row.Scan(&inv.ID, &inv.SessionID, &inv.UserID, &inv.TariffID, &inv.ConsumedKwh, &inv.PricePerKwh, &inv.Subtotal, &inv.Tax, &inv.Total, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt); err == nil {
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
	if r.db != nil {
		query := `SELECT id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status, created_at, updated_at
		          FROM invoices WHERE user_id = $1 ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err == nil {
			defer rows.Close()
			var list []domain.Invoice
			for rows.Next() {
				var inv domain.Invoice
				if err := rows.Scan(&inv.ID, &inv.SessionID, &inv.UserID, &inv.TariffID, &inv.ConsumedKwh, &inv.PricePerKwh, &inv.Subtotal, &inv.Tax, &inv.Total, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt); err == nil {
					list = append(list, inv)
				}
			}
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
	if r.db != nil {
		query := `UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
		res, err := r.db.ExecContext(ctx, query, status, id)
		if err == nil {
			if rows, _ := res.RowsAffected(); rows > 0 {
				return nil
			}
		}
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

func (r *billingRepository) getPaymentsByInvoiceID(ctx context.Context, invoiceID string) ([]domain.Payment, error) {
	if r.db != nil {
		query := `SELECT id, invoice_id, payment_method, amount, status, transaction_ref, paid_at, created_at
		          FROM payments WHERE invoice_id = $1 ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, invoiceID)
		if err == nil {
			defer rows.Close()
			var list []domain.Payment
			for rows.Next() {
				var p domain.Payment
				var txRef sql.NullString
				var paidAt sql.NullTime
				if err := rows.Scan(&p.ID, &p.InvoiceID, &p.PaymentMethod, &p.Amount, &p.Status, &txRef, &paidAt, &p.CreatedAt); err == nil {
					p.TransactionRef = txRef.String
					if paidAt.Valid {
						p.PaidAt = &paidAt.Time
					}
					list = append(list, p)
				}
			}
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Payment
	for _, p := range r.payments {
		if p.InvoiceID == invoiceID {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (r *billingRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	payment.CreatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO payments (id, invoice_id, payment_method, amount, status)
		          VALUES ($1, $2, $3, $4, $5)`
		_, err := r.db.ExecContext(ctx, query, payment.ID, payment.InvoiceID, payment.PaymentMethod, payment.Amount, payment.Status)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments[payment.ID] = payment
	return nil
}

func (r *billingRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	if r.db != nil {
		query := `SELECT id, invoice_id, payment_method, amount, status, transaction_ref, paid_at, created_at
		          FROM payments WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)

		var p domain.Payment
		var txRef sql.NullString
		var paidAt sql.NullTime
		if err := row.Scan(&p.ID, &p.InvoiceID, &p.PaymentMethod, &p.Amount, &p.Status, &txRef, &paidAt, &p.CreatedAt); err == nil {
			p.TransactionRef = txRef.String
			if paidAt.Valid {
				p.PaidAt = &paidAt.Time
			}
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
	if r.db != nil {
		query := `UPDATE payments SET status = $1, paid_at = $2 WHERE id = $3`
		_, err := r.db.ExecContext(ctx, query, status, paidAt, id)
		if err == nil {
			return nil
		}
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
	if r.db != nil {
		oldBytes, _ := json.Marshal(oldVal)
		newBytes, _ := json.Marshal(newVal)
		query := `INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value, performed_by)
		          VALUES ($1, $2, $3, $4, $5, 'billing-service')`
		_, err := r.db.ExecContext(ctx, query, entityName, entityID, action, oldBytes, newBytes)
		return err
	}
	return nil
}

func (r *billingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	if r.db != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		query := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status)
		          VALUES ($1, $2, $3, $4, $5, 'PENDING')`
		_, err = r.db.ExecContext(ctx, query, uuid.New().String(), aggregateType, aggregateID, eventType, payloadBytes)
		return err
	}
	return nil
}
