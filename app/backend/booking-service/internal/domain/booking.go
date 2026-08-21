package domain

import (
	"context"
	"time"
)

type Booking struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	UserID         string    `gorm:"index;size:64" json:"user_id"`
	StationID      string    `gorm:"index;size:64" json:"station_id"`
	SlotID         string    `gorm:"index;size:64" json:"slot_id"`
	StartTime      time.Time `gorm:"index" json:"start_time"`
	EndTime        time.Time `gorm:"index" json:"end_time"`
	Status         string    `gorm:"index;size:32" json:"status"` // REQUESTED, CONFIRMED, WAITLISTED, IN_SESSION, COMPLETED, EXPIRED_NO_SHOW, CANCELLED
	IdempotencyKey string    `gorm:"index;size:128" json:"idempotency_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Waitlist struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	StationID      string    `gorm:"index;size:64" json:"station_id"`
	UserID         string    `gorm:"index;size:64" json:"user_id"`
	RequestedStart time.Time `json:"requested_start"`
	RequestedEnd   time.Time `json:"requested_end"`
	QueueNumber    int       `gorm:"index" json:"queue_number"`
	Status         string    `gorm:"index;size:32" json:"status"` // WAITING, PROMOTED, EXPIRED, CANCELLED
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
	GetBookingsByUserID(ctx context.Context, userID string) ([]Booking, error)
	GetAllBookings(ctx context.Context) ([]Booking, error)
	UpdateBookingStatus(ctx context.Context, id string, status string) error
	CheckSlotAvailability(ctx context.Context, slotID string, start, end time.Time) (bool, error)
	AddToWaitlist(ctx context.Context, waitlist *Waitlist) error
	GetWaitlistByStation(ctx context.Context, stationID string) ([]Waitlist, error)
	AutoReleaseNoShowBookings(ctx context.Context, graceMinutes int) (int, error)
	PromoteNextWaitlist(ctx context.Context, stationID string, slotID string) (*Waitlist, error)
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, booking *Booking) (*Booking, error)
	GetBookingByID(ctx context.Context, id string) (*Booking, error)
	GetBookingsByUserID(ctx context.Context, userID string) ([]Booking, error)
	GetAllBookings(ctx context.Context) ([]Booking, error)
	CheckIn(ctx context.Context, bookingID string) error
	CancelBooking(ctx context.Context, bookingID string) error
	GetAvailability(ctx context.Context, stationID string, start, end time.Time) ([]SlotAvailability, error)
	GetWaitlist(ctx context.Context, stationID string) ([]Waitlist, error)
	TriggerAutoRelease(ctx context.Context, graceMinutes int) (int, error)
	PromoteNextWaitlist(ctx context.Context, stationID string, slotID string) (*Waitlist, error)
}
