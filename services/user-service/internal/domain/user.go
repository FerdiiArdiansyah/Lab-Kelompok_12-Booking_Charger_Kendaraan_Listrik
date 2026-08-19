package domain

import (
	"context"
	"time"
)

type User struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	PasswordHash string        `json:"-"`
	Phone        string        `json:"phone,omitempty"`
	Role         string        `json:"role"` // USER, ADMIN, OPERATOR
	Status       string        `json:"status"` // ACTIVE, SUSPENDED
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Vehicles     []UserVehicle `json:"vehicles,omitempty"`
}

type UserVehicle struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Brand              string    `json:"brand"`
	Model              string    `json:"model"`
	LicensePlate       string    `json:"license_plate"`
	ConnectorType      string    `json:"connector_type"` // Type 2, CCS2, CHAdeMO, GB/T
	BatteryCapacityKwh float64   `json:"battery_capacity_kwh"`
	CreatedAt          time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type AddVehicleRequest struct {
	Brand              string  `json:"brand"`
	Model              string  `json:"model"`
	LicensePlate       string  `json:"license_plate"`
	ConnectorType      string  `json:"connector_type"`
	BatteryCapacityKwh float64 `json:"battery_capacity_kwh"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, user *User) error
	CreateVehicle(ctx context.Context, vehicle *UserVehicle) error
	GetVehiclesByUserID(ctx context.Context, userID string) ([]UserVehicle, error)
	DeleteVehicle(ctx context.Context, id, userID string) error
}

type UserUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateUserRequest) (*User, error)
	AddVehicle(ctx context.Context, userID string, req *AddVehicleRequest) (*UserVehicle, error)
	GetUserVehicles(ctx context.Context, userID string) ([]UserVehicle, error)
	DeleteVehicle(ctx context.Context, vehicleID, userID string) error
}
