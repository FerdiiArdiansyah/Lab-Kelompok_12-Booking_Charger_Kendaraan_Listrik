package postgres

import (
	"booking-service/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) domain.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, query,
		booking.ID, booking.UserID, booking.StationID, booking.SlotID,
		booking.StartTime, booking.EndTime, booking.Status, booking.IdempotencyKey)
	if err != nil {
		if strings.Contains(err.Error(), "no_overlapping_bookings") {
			return errors.New("SLOT_OVERLAP_CONFLICT")
		}
		return err
	}

	payload := map[string]interface{}{
		"booking_id": booking.ID,
		"user_id":    booking.UserID,
		"slot_id":    booking.SlotID,
		"station_id": booking.StationID,
		"start_time": booking.StartTime,
		"end_time":   booking.EndTime,
		"status":     booking.Status,
	}
	payloadBytes, _ := json.Marshal(payload)

	outboxQuery := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status)
	                VALUES ($1, 'Booking', $2, 'BookingCreated', $3, 'PENDING')`
	_, err = tx.ExecContext(ctx, outboxQuery, uuid.New().String(), booking.ID, payloadBytes)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *bookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	query := `SELECT id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key, created_at, updated_at
	          FROM bookings WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var b domain.Booking
	var idemKey sql.NullString
	if err := row.Scan(&b.ID, &b.UserID, &b.StationID, &b.SlotID, &b.StartTime, &b.EndTime, &b.Status, &idemKey, &b.CreatedAt, &b.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	b.IdempotencyKey = idemKey.String
	return &b, nil
}

func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("booking not found")
	}
	return nil
}

func (r *bookingRepository) CheckSlotAvailability(ctx context.Context, slotID string, start, end time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM bookings 
	          WHERE slot_id = $1 
	          AND status IN ('REQUESTED', 'CONFIRMED', 'IN_SESSION')
	          AND tstzrange(start_time, end_time, '[)') && tstzrange($2, $3, '[)')`
	var count int
	err := r.db.QueryRowContext(ctx, query, slotID, start, end).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *bookingRepository) AddToWaitlist(ctx context.Context, waitlist *domain.Waitlist) error {
	query := `INSERT INTO waitlists (id, station_id, user_id, requested_start, requested_end, queue_number, status)
	          VALUES ($1, $2, $3, $4, $5, 
	          (SELECT COALESCE(MAX(queue_number), 0) + 1 FROM waitlists WHERE station_id = $2 AND status = 'WAITING'), $6)`
	_, err := r.db.ExecContext(ctx, query, waitlist.ID, waitlist.StationID, waitlist.UserID, waitlist.RequestedStart, waitlist.RequestedEnd, waitlist.Status)
	return err
}

func (r *bookingRepository) GetWaitlistByStation(ctx context.Context, stationID string) ([]domain.Waitlist, error) {
	query := `SELECT id, station_id, user_id, requested_start, requested_end, queue_number, status, created_at, updated_at
	          FROM waitlists WHERE station_id = $1 AND status = 'WAITING' ORDER BY queue_number ASC`
	rows, err := r.db.QueryContext(ctx, query, stationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Waitlist
	for rows.Next() {
		var w domain.Waitlist
		if err := rows.Scan(&w.ID, &w.StationID, &w.UserID, &w.RequestedStart, &w.RequestedEnd, &w.QueueNumber, &w.Status, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, nil
}

func (r *bookingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	query := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status)
	          VALUES ($1, $2, $3, $4, $5, 'PENDING')`
	_, err = r.db.ExecContext(ctx, query, uuid.New().String(), aggregateType, aggregateID, eventType, payloadBytes)
	return err
}
