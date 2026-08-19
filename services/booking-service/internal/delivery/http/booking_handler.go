package http

import (
	"booking-service/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	usecase domain.BookingUsecase
}

func NewBookingHandler(e *echo.Echo, u domain.BookingUsecase) {
	h := &BookingHandler{usecase: u}

	e.POST("/bookings", h.CreateBooking)
	e.GET("/bookings/:id", h.GetBookingByID)
	e.POST("/bookings/:id/check-in", h.CheckIn)
	e.POST("/bookings/:id/cancel", h.CancelBooking)
}

func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var booking domain.Booking
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	result, err := h.usecase.CreateBooking(c.Request().Context(), &booking)
	if err != nil {
		if result != nil && result.Status == "WAITLISTED" {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"status":  "WAITLISTED",
				"message": err.Error(),
				"data":    result,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "booking confirmed successfully",
		"data":    result,
	})
}

func (h *BookingHandler) GetBookingByID(c echo.Context) error {
	id := c.Param("id")
	booking, err := h.usecase.GetBookingByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   booking,
	})
}

func (h *BookingHandler) CheckIn(c echo.Context) error {
	id := c.Param("id")
	if err := h.usecase.CheckIn(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "check-in successful, session initiated",
	})
}

func (h *BookingHandler) CancelBooking(c echo.Context) error {
	id := c.Param("id")
	if err := h.usecase.CancelBooking(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "booking cancelled successfully",
	})
}
