package usecase

import (
	"billing-service/internal/domain"
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type billingUsecase struct {
	repo domain.BillingRepository
}

func NewBillingUsecase(repo domain.BillingRepository) domain.BillingUsecase {
	return &billingUsecase{repo: repo}
}

func (u *billingUsecase) GenerateInvoice(ctx context.Context, sessionID, userID, tariffID string, consumedKwh, pricePerKwh float64) (*domain.Invoice, error) {
	if sessionID == "" || userID == "" {
		return nil, errors.New("session_id and user_id are required")
	}

	subtotal := math.Round(consumedKwh*pricePerKwh*100) / 100
	tax := math.Round(subtotal*0.11*100) / 100 // PPN 11%
	total := subtotal + tax

	invoice := &domain.Invoice{
		ID:          "inv-" + uuid.New().String(),
		SessionID:   sessionID,
		UserID:      userID,
		TariffID:    tariffID,
		ConsumedKwh: consumedKwh,
		PricePerKwh: pricePerKwh,
		Subtotal:    subtotal,
		Tax:         tax,
		Total:       total,
		Status:      "UNPAID",
	}

	if err := u.repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, err
	}

	_ = u.repo.SaveOutboxEvent(ctx, "Invoice", invoice.ID, "InvoiceCreated", map[string]interface{}{
		"invoice_id": invoice.ID,
		"session_id": sessionID,
		"user_id":    userID,
		"total":      total,
	})

	return invoice, nil
}

func (u *billingUsecase) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	if id == "" {
		return nil, errors.New("invoice ID is required")
	}
	return u.repo.GetInvoiceByID(ctx, id)
}

func (u *billingUsecase) ProcessPayment(ctx context.Context, invoiceID, paymentMethod string, amount float64) (*domain.Payment, error) {
	invoice, err := u.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if invoice.Status == "PAID" {
		return nil, errors.New("invoice is already paid")
	}

	payment := &domain.Payment{
		ID:            "pay-" + uuid.New().String(),
		InvoiceID:     invoiceID,
		PaymentMethod: paymentMethod,
		Amount:        amount,
		Status:        "PENDING",
	}

	if err := u.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (u *billingUsecase) ConfirmPayment(ctx context.Context, paymentID, transactionRef string) (*domain.Payment, error) {
	payment, err := u.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	payment.Status = "SUCCESS"
	payment.TransactionRef = transactionRef
	payment.PaidAt = &now

	if err := u.repo.UpdatePaymentStatus(ctx, paymentID, "SUCCESS", now); err != nil {
		return nil, err
	}

	// Update Invoice to PAID
	_ = u.repo.UpdateInvoiceStatus(ctx, payment.InvoiceID, "PAID")

	_ = u.repo.SaveAuditLog(ctx, "Invoice", payment.InvoiceID, "UPDATE_STATUS", map[string]string{"status": "UNPAID"}, map[string]string{"status": "PAID"})
	_ = u.repo.SaveOutboxEvent(ctx, "Payment", paymentID, "PaymentCompleted", map[string]interface{}{
		"payment_id": paymentID,
		"invoice_id": payment.InvoiceID,
		"amount":     payment.Amount,
		"paid_at":    now,
	})

	return payment, nil
}
