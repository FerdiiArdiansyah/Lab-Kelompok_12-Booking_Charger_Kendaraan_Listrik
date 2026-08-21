package domain

import (
	"context"
	"time"
)

type Station struct {
	ID           string        `gorm:"primaryKey;size:64" json:"id"`
	Name         string        `gorm:"size:128" json:"name"`
	Location     string        `gorm:"size:256" json:"location"`
	Latitude     float64       `json:"latitude"`
	Longitude    float64       `json:"longitude"`
	TotalPower   float64       `json:"total_power_kw"`
	ImageURL     string        `gorm:"size:512" json:"image_url"`
	MapURL       string        `gorm:"size:512" json:"map_url"`
	Status       string        `gorm:"size:32;default:'ACTIVE'" json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Slots        []ChargerSlot `gorm:"foreignKey:StationID" json:"slots,omitempty"`
	ActiveTariff *Tariff       `gorm:"foreignKey:StationID" json:"active_tariff,omitempty"`
}

type ChargerSlot struct {
	ID            string    `gorm:"primaryKey;size:64" json:"id"`
	StationID     string    `gorm:"index;size:64" json:"station_id"`
	SlotNumber    int       `json:"slot_number"`
	ConnectorType string    `gorm:"size:32" json:"connector_type"`
	MaxPowerKw    float64   `json:"max_power_kw"`
	Status        string    `gorm:"size:32;default:'AVAILABLE'" json:"status"` // AVAILABLE, IN_USE, OUT_OF_SERVICE
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Tariff struct {
	ID          string     `gorm:"primaryKey;size:64" json:"id"`
	StationID   string     `gorm:"index;size:64" json:"station_id"`
	PricePerKwh float64    `json:"price_per_kwh"`
	Currency    string     `gorm:"size:8;default:'IDR'" json:"currency"`
	ValidFrom   time.Time  `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to,omitempty"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

type StationRepository interface {
	GetAll(ctx context.Context) ([]Station, error)
	GetByID(ctx context.Context, id string) (*Station, error)
	Create(ctx context.Context, station *Station) error
	Update(ctx context.Context, station *Station) error
	Delete(ctx context.Context, id string) error
	GetSlotsByStationID(ctx context.Context, stationID string) ([]ChargerSlot, error)
	CreateSlot(ctx context.Context, slot *ChargerSlot) error
	UpdateSlot(ctx context.Context, slot *ChargerSlot) error
	GetTariffByStationID(ctx context.Context, stationID string) (*Tariff, error)
	CreateTariff(ctx context.Context, tariff *Tariff) error
	GetAllTariffs(ctx context.Context) ([]Tariff, error)
}

type StationUsecase interface {
	GetAllStations(ctx context.Context) ([]Station, error)
	GetStationByID(ctx context.Context, id string) (*Station, error)
	CreateStation(ctx context.Context, station *Station) error
	UpdateStation(ctx context.Context, id string, station *Station) (*Station, error)
	DeleteStation(ctx context.Context, id string) error
	GetSlots(ctx context.Context, stationID string) ([]ChargerSlot, error)
	AddSlot(ctx context.Context, slot *ChargerSlot) error
	UpdateSlot(ctx context.Context, stationID, slotID string, slot *ChargerSlot) (*ChargerSlot, error)
	GetTariff(ctx context.Context, stationID string) (*Tariff, error)
	AddTariff(ctx context.Context, tariff *Tariff) error
	GetAllTariffs(ctx context.Context) ([]Tariff, error)
}
