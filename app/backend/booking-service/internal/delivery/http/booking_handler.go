package http

import (
	"booking-service/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	usecase domain.BookingUsecase
}

func NewBookingHandler(u domain.BookingUsecase) *BookingHandler {
	return &BookingHandler{usecase: u}
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

func (h *BookingHandler) GetBookingsByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	list, err := h.usecase.GetBookingsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   list,
	})
}

func (h *BookingHandler) GetAllBookings(c echo.Context) error {
	list, err := h.usecase.GetAllBookings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   list,
	})
}

func (h *BookingHandler) GetAvailability(c echo.Context) error {
	stationID := c.Param("id")
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	var start, end time.Time
	if startStr != "" {
		start, _ = time.Parse(time.RFC3339, startStr)
	}
	if endStr != "" {
		end, _ = time.Parse(time.RFC3339, endStr)
	}

	list, err := h.usecase.GetAvailability(c.Request().Context(), stationID, start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   list,
	})
}

func (h *BookingHandler) GetWaitlist(c echo.Context) error {
	stationID := c.Param("id")
	list, err := h.usecase.GetWaitlist(c.Request().Context(), stationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   list,
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

func (h *BookingHandler) CompleteBooking(c echo.Context) error {
	id := c.Param("id")
	if err := h.usecase.CompleteBooking(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "booking marked as COMPLETED",
	})
}

func (h *BookingHandler) AutoRelease(c echo.Context) error {
	graceStr := c.QueryParam("grace_period_minutes")
	grace, _ := strconv.Atoi(graceStr)

	count, err := h.usecase.TriggerAutoRelease(c.Request().Context(), grace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":         "success",
		"message":        "auto-release processed",
		"released_count": count,
	})
}
