package postgres

import (
	"booking-service/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type bookingRepository struct {
	db       *sql.DB
	bookings map[string]*domain.Booking
	waitlist map[string]*domain.Waitlist
	mu       sync.RWMutex
}

func NewBookingRepository(db *sql.DB) domain.BookingRepository {
	repo := &bookingRepository{
		db:       db,
		bookings: make(map[string]*domain.Booking),
		waitlist: make(map[string]*domain.Waitlist),
	}
	repo.seedInitialData()
	return repo
}

func (r *bookingRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// 11 Seeded Bookings
	sampleBookings := []domain.Booking{
		{ID: "bkg-001", UserID: "usr-001", StationID: "stn-001", SlotID: "slot-001", StartTime: now.Add(-2 * time.Hour), EndTime: now.Add(-1 * time.Hour), Status: "COMPLETED", IdempotencyKey: "idem-001", CreatedAt: now.Add(-3 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)},
		{ID: "bkg-002", UserID: "usr-002", StationID: "stn-002", SlotID: "slot-003", StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour), Status: "IN_SESSION", IdempotencyKey: "idem-002", CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)},
		{ID: "bkg-003", UserID: "usr-003", StationID: "stn-003", SlotID: "slot-004", StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-003", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-004", UserID: "usr-004", StationID: "stn-004", SlotID: "slot-006", StartTime: now.Add(2 * time.Hour), EndTime: now.Add(3 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-004", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-005", UserID: "usr-005", StationID: "stn-005", SlotID: "slot-007", StartTime: now.Add(-4 * time.Hour), EndTime: now.Add(-3 * time.Hour), Status: "COMPLETED", IdempotencyKey: "idem-005", CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: now.Add(-3 * time.Hour)},
		{ID: "bkg-006", UserID: "usr-006", StationID: "stn-006", SlotID: "slot-008", StartTime: now.Add(3 * time.Hour), EndTime: now.Add(4 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-006", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-007", UserID: "usr-007", StationID: "stn-007", SlotID: "slot-009", StartTime: now.Add(-5 * time.Hour), EndTime: now.Add(-4 * time.Hour), Status: "EXPIRED_NO_SHOW", IdempotencyKey: "idem-007", CreatedAt: now.Add(-6 * time.Hour), UpdatedAt: now.Add(-4 * time.Hour)},
		{ID: "bkg-008", UserID: "usr-008", StationID: "stn-008", SlotID: "slot-010", StartTime: now.Add(4 * time.Hour), EndTime: now.Add(5 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-008", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-009", UserID: "usr-009", StationID: "stn-010", SlotID: "slot-011", StartTime: now.Add(-3 * time.Hour), EndTime: now.Add(-2 * time.Hour), Status: "CANCELLED", IdempotencyKey: "idem-009", CreatedAt: now.Add(-4 * time.Hour), UpdatedAt: now.Add(-3 * time.Hour)},
		{ID: "bkg-010", UserID: "usr-010", StationID: "stn-001", SlotID: "slot-002", StartTime: now.Add(5 * time.Hour), EndTime: now.Add(6 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-010", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-011", UserID: "usr-011", StationID: "stn-003", SlotID: "slot-005", StartTime: now.Add(6 * time.Hour), EndTime: now.Add(7 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-011", CreatedAt: now, UpdatedAt: now},
	}

	for i := range sampleBookings {
		b := sampleBookings[i]
		r.bookings[b.ID] = &b
	}

	// 11 Seeded Waitlists
	sampleWaitlists := []domain.Waitlist{
		{ID: "wt-001", StationID: "stn-001", UserID: "usr-002", RequestedStart: now.Add(1 * time.Hour), RequestedEnd: now.Add(2 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-002", StationID: "stn-001", UserID: "usr-003", RequestedStart: now.Add(1 * time.Hour), RequestedEnd: now.Add(2 * time.Hour), QueueNumber: 2, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-003", StationID: "stn-002", UserID: "usr-004", RequestedStart: now.Add(2 * time.Hour), RequestedEnd: now.Add(3 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-004", StationID: "stn-003", UserID: "usr-005", RequestedStart: now.Add(1 * time.Hour), RequestedEnd: now.Add(2 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-005", StationID: "stn-003", UserID: "usr-006", RequestedStart: now.Add(1 * time.Hour), RequestedEnd: now.Add(2 * time.Hour), QueueNumber: 2, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-006", StationID: "stn-004", UserID: "usr-007", RequestedStart: now.Add(3 * time.Hour), RequestedEnd: now.Add(4 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-007", StationID: "stn-005", UserID: "usr-008", RequestedStart: now.Add(2 * time.Hour), RequestedEnd: now.Add(3 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-008", StationID: "stn-006", UserID: "usr-009", RequestedStart: now.Add(4 * time.Hour), RequestedEnd: now.Add(5 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-009", StationID: "stn-007", UserID: "usr-010", RequestedStart: now.Add(1 * time.Hour), RequestedEnd: now.Add(2 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-010", StationID: "stn-008", UserID: "usr-011", RequestedStart: now.Add(5 * time.Hour), RequestedEnd: now.Add(6 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
		{ID: "wt-011", StationID: "stn-010", UserID: "usr-001", RequestedStart: now.Add(2 * time.Hour), RequestedEnd: now.Add(3 * time.Hour), QueueNumber: 1, Status: "WAITING", CreatedAt: now, UpdatedAt: now},
	}

	for i := range sampleWaitlists {
		w := sampleWaitlists[i]
		r.waitlist[w.ID] = &w
	}
}

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	if r.db != nil {
		tx, err := r.db.BeginTx(ctx, nil)
		if err == nil {
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
			tx.ExecContext(ctx, outboxQuery, uuid.New().String(), booking.ID, payloadBytes)
			if err := tx.Commit(); err == nil {
				return nil
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.bookings[booking.ID] = booking
	return nil
}

func (r *bookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	if r.db != nil {
		query := `SELECT id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key, created_at, updated_at
		          FROM bookings WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)

		var b domain.Booking
		var idemKey sql.NullString
		if err := row.Scan(&b.ID, &b.UserID, &b.StationID, &b.SlotID, &b.StartTime, &b.EndTime, &b.Status, &idemKey, &b.CreatedAt, &b.UpdatedAt); err == nil {
			b.IdempotencyKey = idemKey.String
			return &b, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if b, ok := r.bookings[id]; ok {
		return b, nil
	}
	return nil, errors.New("booking not found")
}

func (r *bookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]domain.Booking, error) {
	if r.db != nil {
		query := `SELECT id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key, created_at, updated_at
		          FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err == nil {
			defer rows.Close()
			var list []domain.Booking
			for rows.Next() {
				var b domain.Booking
				var idemKey sql.NullString
				if err := rows.Scan(&b.ID, &b.UserID, &b.StationID, &b.SlotID, &b.StartTime, &b.EndTime, &b.Status, &idemKey, &b.CreatedAt, &b.UpdatedAt); err == nil {
					b.IdempotencyKey = idemKey.String
					list = append(list, b)
				}
			}
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Booking
	for _, b := range r.bookings {
		if b.UserID == userID {
			list = append(list, *b)
		}
	}
	return list, nil
}

func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	if r.db != nil {
		query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
		res, err := r.db.ExecContext(ctx, query, status, id)
		if err == nil {
			if rows, _ := res.RowsAffected(); rows > 0 {
				return nil
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if b, ok := r.bookings[id]; ok {
		b.Status = status
		b.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("booking not found")
}

func (r *bookingRepository) CheckSlotAvailability(ctx context.Context, slotID string, start, end time.Time) (bool, error) {
	if r.db != nil {
		query := `SELECT COUNT(*) FROM bookings 
		          WHERE slot_id = $1 
		          AND status IN ('REQUESTED', 'CONFIRMED', 'IN_SESSION')
		          AND tstzrange(start_time, end_time, '[)') && tstzrange($2, $3, '[)')`
		var count int
		if err := r.db.QueryRowContext(ctx, query, slotID, start, end).Scan(&count); err == nil {
			return count == 0, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.bookings {
		if b.SlotID == slotID && (b.Status == "REQUESTED" || b.Status == "CONFIRMED" || b.Status == "IN_SESSION") {
			if start.Before(b.EndTime) && end.After(b.StartTime) {
				return false, nil
			}
		}
	}
	return true, nil
}

func (r *bookingRepository) AddToWaitlist(ctx context.Context, waitlist *domain.Waitlist) error {
	waitlist.CreatedAt = time.Now()
	waitlist.UpdatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO waitlists (id, station_id, user_id, requested_start, requested_end, queue_number, status)
		          VALUES ($1, $2, $3, $4, $5, 
		          (SELECT COALESCE(MAX(queue_number), 0) + 1 FROM waitlists WHERE station_id = $2 AND status = 'WAITING'), $6)`
		_, err := r.db.ExecContext(ctx, query, waitlist.ID, waitlist.StationID, waitlist.UserID, waitlist.RequestedStart, waitlist.RequestedEnd, waitlist.Status)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	queueNum := 1
	for _, w := range r.waitlist {
		if w.StationID == waitlist.StationID && w.Status == "WAITING" {
			queueNum++
		}
	}
	waitlist.QueueNumber = queueNum
	r.waitlist[waitlist.ID] = waitlist
	return nil
}

func (r *bookingRepository) GetWaitlistByStation(ctx context.Context, stationID string) ([]domain.Waitlist, error) {
	if r.db != nil {
		query := `SELECT id, station_id, user_id, requested_start, requested_end, queue_number, status, created_at, updated_at
		          FROM waitlists WHERE station_id = $1 AND status = 'WAITING' ORDER BY queue_number ASC`
		rows, err := r.db.QueryContext(ctx, query, stationID)
		if err == nil {
			defer rows.Close()
			var list []domain.Waitlist
			for rows.Next() {
				var w domain.Waitlist
				if err := rows.Scan(&w.ID, &w.StationID, &w.UserID, &w.RequestedStart, &w.RequestedEnd, &w.QueueNumber, &w.Status, &w.CreatedAt, &w.UpdatedAt); err == nil {
					list = append(list, w)
				}
			}
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Waitlist
	for _, w := range r.waitlist {
		if w.StationID == stationID && w.Status == "WAITING" {
			list = append(list, *w)
		}
	}
	return list, nil
}

func (r *bookingRepository) AutoReleaseNoShowBookings(ctx context.Context, graceMinutes int) (int, error) {
	cutoff := time.Now().Add(-time.Duration(graceMinutes) * time.Minute)

	if r.db != nil {
		query := `UPDATE bookings SET status = 'EXPIRED_NO_SHOW', updated_at = CURRENT_TIMESTAMP 
		          WHERE status = 'CONFIRMED' AND start_time <= $1`
		res, err := r.db.ExecContext(ctx, query, cutoff)
		if err == nil {
			rows, _ := res.RowsAffected()
			return int(rows), nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for _, b := range r.bookings {
		if b.Status == "CONFIRMED" && b.StartTime.Before(cutoff) {
			b.Status = "EXPIRED_NO_SHOW"
			b.UpdatedAt = time.Now()
			count++
		}
	}
	return count, nil
}

func (r *bookingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	if r.db != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		query := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status)
		          VALUES ($1, $2, $3, $4, $5, 'PENDING')`
		_, err = r.db.ExecContext(ctx, query, uuid.New().String(), aggregateType, aggregateID, eventType, payloadBytes)
		return err
	}
	return nil
}
