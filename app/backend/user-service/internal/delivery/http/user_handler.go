package http

import (
	"net/http"
	"user-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	res, err := h.usecase.Register(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "user registered successfully",
		"data":    res,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	res, err := h.usecase.Login(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "login successful",
		"data":    res,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}
	if userID == "" {
		userID = "usr-demo-001"
	}

	user, err := h.usecase.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}
	if userID == "" {
		userID = "usr-demo-001"
	}

	var req domain.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	updated, err := h.usecase.UpdateProfile(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "profile updated successfully",
		"data":    updated,
	})
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.usecase.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   users,
	})
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.usecase.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) GetVehicles(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}
	if userID == "" {
		userID = "usr-demo-001"
	}

	vehicles, err := h.usecase.GetUserVehicles(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   vehicles,
	})
}

func (h *UserHandler) AddVehicle(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}
	if userID == "" {
		userID = "usr-demo-001"
	}

	var req domain.AddVehicleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	vehicle, err := h.usecase.AddVehicle(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "vehicle added successfully",
		"data":    vehicle,
	})
}

func (h *UserHandler) DeleteVehicle(c echo.Context) error {
	vehicleID := c.Param("vehicle_id")
	userID := c.QueryParam("user_id")
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}
	if userID == "" {
		userID = "usr-demo-001"
	}

	if err := h.usecase.DeleteVehicle(c.Request().Context(), vehicleID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "vehicle deleted successfully",
	})
}
