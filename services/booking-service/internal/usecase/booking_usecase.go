package usecase

import (
	"booking-service/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type bookingUsecase struct {
	repo domain.BookingRepository
}

func NewBookingUsecase(repo domain.BookingRepository) domain.BookingUsecase {
	return &bookingUsecase{repo: repo}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	if booking.UserID == "" || booking.StationID == "" || booking.SlotID == "" {
		return nil, errors.New("user_id, station_id, and slot_id are required")
	}
	if booking.EndTime.Before(booking.StartTime) || booking.EndTime.Equal(booking.StartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	if booking.ID == "" {
		booking.ID = "bkg-" + uuid.New().String()
	}
	booking.Status = "CONFIRMED"

	err := u.repo.CreateBooking(ctx, booking)
	if err != nil {
		if err.Error() == "SLOT_OVERLAP_CONFLICT" {
			// Masukkan ke waitlist jika slot penuh
			waitlist := &domain.Waitlist{
				ID:             "wt-" + uuid.New().String(),
				StationID:      booking.StationID,
				UserID:         booking.UserID,
				RequestedStart: booking.StartTime,
				RequestedEnd:   booking.EndTime,
				Status:         "WAITING",
			}
			_ = u.repo.AddToWaitlist(ctx, waitlist)

			booking.Status = "WAITLISTED"
			return booking, errors.New("slot is already booked for this timeframe; user placed on WAITLIST FIFO queue")
		}
		return nil, err
	}

	return booking, nil
}

func (u *bookingUsecase) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	if id == "" {
		return nil, errors.New("booking ID is required")
	}
	return u.repo.GetBookingByID(ctx, id)
}

func (u *bookingUsecase) CheckIn(ctx context.Context, bookingID string) error {
	booking, err := u.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking.Status != "CONFIRMED" {
		return errors.New("only CONFIRMED bookings can check in")
	}

	// Update status ke IN_SESSION
	if err := u.repo.UpdateBookingStatus(ctx, bookingID, "IN_SESSION"); err != nil {
		return err
	}

	_ = u.repo.SaveOutboxEvent(ctx, "Booking", bookingID, "BookingCheckedIn", map[string]string{
		"booking_id": bookingID,
		"user_id":    booking.UserID,
		"slot_id":    booking.SlotID,
	})
	return nil
}

func (u *bookingUsecase) CancelBooking(ctx context.Context, bookingID string) error {
	booking, err := u.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking.Status == "COMPLETED" || booking.Status == "CANCELLED" {
		return errors.New("cannot cancel a completed or already cancelled booking")
	}

	if err := u.repo.UpdateBookingStatus(ctx, bookingID, "CANCELLED"); err != nil {
		return err
	}

	_ = u.repo.SaveOutboxEvent(ctx, "Booking", bookingID, "BookingCancelled", map[string]string{
		"booking_id": bookingID,
		"slot_id":    booking.SlotID,
	})
	return nil
}

func (u *bookingUsecase) GetAvailability(ctx context.Context, stationID string, start, end time.Time) ([]domain.SlotAvailability, error) {
	// Simple availability checker for slots
	return []domain.SlotAvailability{}, nil
}
