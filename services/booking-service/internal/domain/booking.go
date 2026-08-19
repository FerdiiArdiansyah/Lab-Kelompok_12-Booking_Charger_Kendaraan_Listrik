package domain

import (
	"context"
	"time"
)

type Booking struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	StationID      string    `json:"station_id"`
	SlotID         string    `json:"slot_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Status         string    `json:"status"` // REQUESTED, CONFIRMED, WAITLISTED, IN_SESSION, COMPLETED, EXPIRED_NO_SHOW, CANCELLED
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Waitlist struct {
	ID             string    `json:"id"`
	StationID      string    `json:"station_id"`
	UserID         string    `json:"user_id"`
	RequestedStart time.Time `json:"requested_start"`
	RequestedEnd   time.Time `json:"requested_end"`
	QueueNumber    int       `json:"queue_number"`
	Status         string    `json:"status"` // WAITING, PROMOTED, EXPIRED, CANCELLED
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SlotAvailability struct {
	SlotID    string `json:"slot_id"`
	Available bool   `json:"available"`
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBookingByID(ctx context.Context, id string) (*Booking, error)
	UpdateBookingStatus(ctx context.Context, id string, status string) error
	CheckSlotAvailability(ctx context.Context, slotID string, start, end time.Time) (bool, error)
	AddToWaitlist(ctx context.Context, waitlist *Waitlist) error
	GetWaitlistByStation(ctx context.Context, stationID string) ([]Waitlist, error)
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, booking *Booking) (*Booking, error)
	GetBookingByID(ctx context.Context, id string) (*Booking, error)
	CheckIn(ctx context.Context, bookingID string) error
	CancelBooking(ctx context.Context, bookingID string) error
	GetAvailability(ctx context.Context, stationID string, start, end time.Time) ([]SlotAvailability, error)
}
