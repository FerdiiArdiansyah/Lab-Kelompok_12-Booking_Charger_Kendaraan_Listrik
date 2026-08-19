package http

import (
	"net/http"
	"session-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	usecase domain.SessionUsecase
}

func NewSessionHandler(e *echo.Echo, u domain.SessionUsecase) {
	h := &SessionHandler{usecase: u}

	e.POST("/sessions/start", h.StartSession)
	e.GET("/sessions/:id", h.GetSessionByID)
	e.POST("/sessions/:id/meter", h.RecordMeter)
	e.POST("/sessions/:id/finish", h.FinishSession)
}

type StartSessionRequest struct {
	BookingID string `json:"booking_id"`
	SlotID    string `json:"slot_id"`
	UserID    string `json:"user_id"`
}

type MeterReadingRequest struct {
	CurrentKwh float64 `json:"current_kwh"`
	PowerKw    float64 `json:"power_kw"`
	Voltage    float64 `json:"voltage"`
	Ampere     float64 `json:"current_ampere"`
}

type FinishSessionRequest struct {
	FinalKwh float64 `json:"consumed_kwh"`
}

func (h *SessionHandler) StartSession(c echo.Context) error {
	var req StartSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	session, err := h.usecase.StartSession(c.Request().Context(), req.BookingID, req.SlotID, req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "charging session started",
		"data":    session,
	})
}

func (h *SessionHandler) GetSessionByID(c echo.Context) error {
	id := c.Param("id")
	session, err := h.usecase.GetSessionByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   session,
	})
}

func (h *SessionHandler) RecordMeter(c echo.Context) error {
	id := c.Param("id")
	var req MeterReadingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if err := h.usecase.RecordMeter(c.Request().Context(), id, req.CurrentKwh, req.PowerKw, req.Voltage, req.Ampere); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "meter reading recorded",
	})
}

func (h *SessionHandler) FinishSession(c echo.Context) error {
	id := c.Param("id")
	var req FinishSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	session, err := h.usecase.FinishSession(c.Request().Context(), id, req.FinalKwh)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "charging session finished",
		"data":    session,
	})
}
