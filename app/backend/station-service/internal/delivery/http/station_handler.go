package http

import (
	"net/http"
	"station-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type StationHandler struct {
	usecase domain.StationUsecase
}

func NewStationHandler(u domain.StationUsecase) *StationHandler {
	return &StationHandler{usecase: u}
}

func (h *StationHandler) GetAll(c echo.Context) error {
	stations, err := h.usecase.GetAllStations(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   stations,
	})
}

func (h *StationHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	station, err := h.usecase.GetStationByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   station,
	})
}

func (h *StationHandler) Create(c echo.Context) error {
	var station domain.Station
	if err := c.Bind(&station); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if err := h.usecase.CreateStation(c.Request().Context(), &station); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "station created successfully",
		"data":    station,
	})
}

func (h *StationHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var station domain.Station
	if err := c.Bind(&station); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	updated, err := h.usecase.UpdateStation(c.Request().Context(), id, &station)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "station updated successfully",
		"data":    updated,
	})
}

func (h *StationHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.usecase.DeleteStation(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "station deleted successfully",
	})
}

func (h *StationHandler) GetSlots(c echo.Context) error {
	stationID := c.Param("id")
	slots, err := h.usecase.GetSlots(c.Request().Context(), stationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   slots,
	})
}

func (h *StationHandler) AddSlot(c echo.Context) error {
	stationID := c.Param("id")
	var slot domain.ChargerSlot
	if err := c.Bind(&slot); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	slot.StationID = stationID

	if err := h.usecase.AddSlot(c.Request().Context(), &slot); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "slot added successfully",
		"data":    slot,
	})
}

func (h *StationHandler) UpdateSlot(c echo.Context) error {
	stationID := c.Param("id")
	slotID := c.Param("slot_id")

	var slot domain.ChargerSlot
	if err := c.Bind(&slot); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	updated, err := h.usecase.UpdateSlot(c.Request().Context(), stationID, slotID, &slot)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "slot updated successfully",
		"data":    updated,
	})
}

func (h *StationHandler) GetTariff(c echo.Context) error {
	stationID := c.Param("id")
	tariff, err := h.usecase.GetTariff(c.Request().Context(), stationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   tariff,
	})
}

func (h *StationHandler) AddTariff(c echo.Context) error {
	stationID := c.Param("id")
	var tariff domain.Tariff
	if err := c.Bind(&tariff); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	tariff.StationID = stationID

	if err := h.usecase.AddTariff(c.Request().Context(), &tariff); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "tariff added successfully",
		"data":    tariff,
	})
}

func (h *StationHandler) GetAllTariffs(c echo.Context) error {
	tariffs, err := h.usecase.GetAllTariffs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   tariffs,
	})
}
