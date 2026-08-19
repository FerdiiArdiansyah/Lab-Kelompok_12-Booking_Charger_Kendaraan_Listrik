package domain

import (
	"context"
	"time"
)

type Station struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Location    string        `json:"location"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	TotalPower  float64       `json:"total_power_kw"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Slots       []ChargerSlot `json:"slots,omitempty"`
	ActiveTariff *Tariff      `json:"active_tariff,omitempty"`
}

type ChargerSlot struct {
	ID            string    `json:"id"`
	StationID     string    `json:"station_id"`
	SlotNumber    int       `json:"slot_number"`
	ConnectorType string    `json:"connector_type"`
	MaxPowerKw    float64   `json:"max_power_kw"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Tariff struct {
	ID          string    `json:"id"`
	StationID   string    `json:"station_id"`
	PricePerKwh float64   `json:"price_per_kwh"`
	Currency    string    `json:"currency"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     *time.Time`json:"valid_to,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// Interfaces untuk Clean Architecture (Repository & Usecase)
type StationRepository interface {
	GetAll(ctx context.Context) ([]Station, error)
	GetByID(ctx context.Context, id string) (*Station, error)
	Create(ctx context.Context, station *Station) error
	GetSlotsByStationID(ctx context.Context, stationID string) ([]ChargerSlot, error)
	CreateSlot(ctx context.Context, slot *ChargerSlot) error
	GetTariffByStationID(ctx context.Context, stationID string) (*Tariff, error)
	CreateTariff(ctx context.Context, tariff *Tariff) error
}

type StationUsecase interface {
	GetAllStations(ctx context.Context) ([]Station, error)
	GetStationByID(ctx context.Context, id string) (*Station, error)
	CreateStation(ctx context.Context, station *Station) error
	GetSlots(ctx context.Context, stationID string) ([]ChargerSlot, error)
	AddSlot(ctx context.Context, slot *ChargerSlot) error
	GetTariff(ctx context.Context, stationID string) (*Tariff, error)
	AddTariff(ctx context.Context, tariff *Tariff) error
}
