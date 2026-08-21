package http

import (
	"billing-service/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BillingHandler struct {
	usecase domain.BillingUsecase
}

func NewBillingHandler(u domain.BillingUsecase) *BillingHandler {
	return &BillingHandler{usecase: u}
}

type GenerateInvoiceRequest struct {
	SessionID   string  `json:"session_id"`
	UserID      string  `json:"user_id"`
	TariffID    string  `json:"tariff_id"`
	ConsumedKwh float64 `json:"consumed_kwh"`
	PricePerKwh float64 `json:"price_per_kwh"`
	ServiceFee  float64 `json:"service_fee"`
}

type ProcessPaymentRequest struct {
	InvoiceID     string  `json:"invoice_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}

type ConfirmPaymentRequest struct {
	TransactionRef string `json:"transaction_ref"`
}

func (h *BillingHandler) GenerateInvoice(c echo.Context) error {
	var req GenerateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	invoice, err := h.usecase.GenerateInvoice(c.Request().Context(), req.SessionID, req.UserID, req.TariffID, req.ConsumedKwh, req.PricePerKwh, req.ServiceFee)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "invoice generated successfully",
		"data":    invoice,
	})
}

func (h *BillingHandler) GetInvoiceByID(c echo.Context) error {
	id := c.Param("id")
	invoice, err := h.usecase.GetInvoiceByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   invoice,
	})
}

func (h *BillingHandler) GetInvoiceBySessionID(c echo.Context) error {
	sessionID := c.Param("session_id")
	invoice, err := h.usecase.GetInvoiceBySessionID(c.Request().Context(), sessionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   invoice,
	})
}

func (h *BillingHandler) GetInvoicesByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	invoices, err := h.usecase.GetInvoicesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   invoices,
	})
}

func (h *BillingHandler) ProcessPayment(c echo.Context) error {
	var req ProcessPaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	payment, err := h.usecase.ProcessPayment(c.Request().Context(), req.InvoiceID, req.PaymentMethod, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "payment process initiated",
		"data":    payment,
	})
}

func (h *BillingHandler) GetPaymentByID(c echo.Context) error {
	id := c.Param("id")
	payment, err := h.usecase.GetPaymentByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   payment,
	})
}

func (h *BillingHandler) ConfirmPayment(c echo.Context) error {
	id := c.Param("id")
	var req ConfirmPaymentRequest
	_ = c.Bind(&req)

	if req.TransactionRef == "" {
		req.TransactionRef = "TX-REF-" + id
	}

	payment, err := h.usecase.ConfirmPayment(c.Request().Context(), id, req.TransactionRef)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "payment confirmed and invoice settled",
		"data":    payment,
	})
}
