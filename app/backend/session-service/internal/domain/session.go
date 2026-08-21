package domain

import (
	"context"
	"time"
)

type ChargingSession struct {
	ID          string         `gorm:"primaryKey;size:64" json:"id"`
	BookingID   string         `gorm:"index;size:64" json:"booking_id"`
	SlotID      string         `gorm:"index;size:64" json:"slot_id"`
	UserID      string         `gorm:"index;size:64" json:"user_id"`
	StartedAt   time.Time      `json:"started_at"`
	EndedAt     *time.Time     `json:"ended_at,omitempty"`
	ConsumedKwh float64        `json:"consumed_kwh"`
	Status      string         `gorm:"size:32;default:'IN_PROGRESS'" json:"status"` // IN_PROGRESS, COMPLETED, INTERRUPTED, CANCELLED
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Readings    []MeterReading `gorm:"foreignKey:SessionID" json:"readings,omitempty"`
}

type MeterReading struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID     string    `gorm:"index;size:64" json:"session_id"`
	RecordedAt    time.Time `json:"recorded_at"`
	CurrentKwh    float64   `json:"current_kwh"`
	PowerKw       float64   `json:"power_kw"`
	Voltage       float64   `json:"voltage,omitempty"`
	CurrentAmpere float64   `json:"current_ampere,omitempty"`
}


type SessionRepository interface {
	CreateSession(ctx context.Context, session *ChargingSession) error
	GetSessionByID(ctx context.Context, id string) (*ChargingSession, error)
	GetSessionByBookingID(ctx context.Context, bookingID string) (*ChargingSession, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]ChargingSession, error)
	AddMeterReading(ctx context.Context, reading *MeterReading) error
	FinishSession(ctx context.Context, id string, endedAt time.Time, finalKwh float64) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error
}

type SessionUsecase interface {
	StartSession(ctx context.Context, bookingID, slotID, userID string) (*ChargingSession, error)
	GetSessionByID(ctx context.Context, id string) (*ChargingSession, error)
	GetSessionByBookingID(ctx context.Context, bookingID string) (*ChargingSession, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]ChargingSession, error)
	RecordMeter(ctx context.Context, sessionID string, currentKwh, powerKw, voltage, ampere float64) error
	FinishSession(ctx context.Context, sessionID string, finalKwh float64) (*ChargingSession, error)
}
