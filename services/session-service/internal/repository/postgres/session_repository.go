package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"session-service/internal/domain"
	"sync"
	"time"

	"github.com/google/uuid"
)

type sessionRepository struct {
	db       *sql.DB
	sessions map[string]*domain.ChargingSession
	readings map[string][]domain.MeterReading
	mu       sync.RWMutex
}

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	repo := &sessionRepository{
		db:       db,
		sessions: make(map[string]*domain.ChargingSession),
		readings: make(map[string][]domain.MeterReading),
	}
	repo.seedInitialData()
	return repo
}

func (r *sessionRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// 11 Seeded Charging Sessions
	ended1 := now.Add(-1 * time.Hour)
	ended5 := now.Add(-3 * time.Hour)

	sampleSessions := []domain.ChargingSession{
		{ID: "ses-001", BookingID: "bkg-001", SlotID: "slot-001", UserID: "usr-001", StartedAt: now.Add(-2 * time.Hour), EndedAt: &ended1, ConsumedKwh: 25.5, Status: "COMPLETED", CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: ended1},
		{ID: "ses-002", BookingID: "bkg-002", SlotID: "slot-003", UserID: "usr-002", StartedAt: now.Add(-1 * time.Hour), EndedAt: nil, ConsumedKwh: 18.2, Status: "IN_PROGRESS", CreatedAt: now.Add(-1 * time.Hour), UpdatedAt: now},
		{ID: "ses-003", BookingID: "bkg-005", SlotID: "slot-007", UserID: "usr-005", StartedAt: now.Add(-4 * time.Hour), EndedAt: &ended5, ConsumedKwh: 42.0, Status: "COMPLETED", CreatedAt: now.Add(-4 * time.Hour), UpdatedAt: ended5},
		{ID: "ses-004", BookingID: "bkg-003", SlotID: "slot-004", UserID: "usr-003", StartedAt: now.Add(-5 * time.Hour), EndedAt: &ended1, ConsumedKwh: 55.0, Status: "COMPLETED", CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: ended1},
		{ID: "ses-005", BookingID: "bkg-004", SlotID: "slot-006", UserID: "usr-004", StartedAt: now.Add(-6 * time.Hour), EndedAt: &ended5, ConsumedKwh: 30.1, Status: "COMPLETED", CreatedAt: now.Add(-6 * time.Hour), UpdatedAt: ended5},
		{ID: "ses-006", BookingID: "bkg-006", SlotID: "slot-008", UserID: "usr-006", StartedAt: now.Add(-7 * time.Hour), EndedAt: &ended1, ConsumedKwh: 22.8, Status: "COMPLETED", CreatedAt: now.Add(-7 * time.Hour), UpdatedAt: ended1},
		{ID: "ses-007", BookingID: "bkg-007", SlotID: "slot-009", UserID: "usr-007", StartedAt: now.Add(-8 * time.Hour), EndedAt: &ended5, ConsumedKwh: 15.4, Status: "COMPLETED", CreatedAt: now.Add(-8 * time.Hour), UpdatedAt: ended5},
		{ID: "ses-008", BookingID: "bkg-008", SlotID: "slot-010", UserID: "usr-008", StartedAt: now.Add(-9 * time.Hour), EndedAt: &ended1, ConsumedKwh: 60.0, Status: "COMPLETED", CreatedAt: now.Add(-9 * time.Hour), UpdatedAt: ended1},
		{ID: "ses-009", BookingID: "bkg-009", SlotID: "slot-011", UserID: "usr-009", StartedAt: now.Add(-10 * time.Hour), EndedAt: &ended5, ConsumedKwh: 38.9, Status: "COMPLETED", CreatedAt: now.Add(-10 * time.Hour), UpdatedAt: ended5},
		{ID: "ses-010", BookingID: "bkg-010", SlotID: "slot-002", UserID: "usr-010", StartedAt: now.Add(-11 * time.Hour), EndedAt: &ended1, ConsumedKwh: 19.5, Status: "COMPLETED", CreatedAt: now.Add(-11 * time.Hour), UpdatedAt: ended1},
		{ID: "ses-011", BookingID: "bkg-011", SlotID: "slot-005", UserID: "usr-011", StartedAt: now.Add(-12 * time.Hour), EndedAt: &ended5, ConsumedKwh: 48.0, Status: "COMPLETED", CreatedAt: now.Add(-12 * time.Hour), UpdatedAt: ended5},
	}

	for i := range sampleSessions {
		s := sampleSessions[i]
		r.sessions[s.ID] = &s

		// Seed initial meter readings for each session
		r.readings[s.ID] = []domain.MeterReading{
			{ID: int64(i*10 + 1), SessionID: s.ID, RecordedAt: s.StartedAt, CurrentKwh: s.ConsumedKwh / 2, PowerKw: 50.0, Voltage: 400.0, CurrentAmpere: 125.0},
			{ID: int64(i*10 + 2), SessionID: s.ID, RecordedAt: s.StartedAt.Add(30 * time.Minute), CurrentKwh: s.ConsumedKwh, PowerKw: 50.0, Voltage: 400.0, CurrentAmpere: 125.0},
		}
	}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.ChargingSession) error {
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO charging_sessions (id, booking_id, slot_id, user_id, started_at, consumed_kwh, status)
		          VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := r.db.ExecContext(ctx, query,
			session.ID, session.BookingID, session.SlotID, session.UserID,
			session.StartedAt, session.ConsumedKwh, session.Status)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*domain.ChargingSession, error) {
	if r.db != nil {
		query := `SELECT id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status, created_at, updated_at
		          FROM charging_sessions WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)

		var s domain.ChargingSession
		var endedAt sql.NullTime
		if err := row.Scan(&s.ID, &s.BookingID, &s.SlotID, &s.UserID, &s.StartedAt, &endedAt, &s.ConsumedKwh, &s.Status, &s.CreatedAt, &s.UpdatedAt); err == nil {
			if endedAt.Valid {
				s.EndedAt = &endedAt.Time
			}
			readings, _ := r.getMeterReadings(ctx, id)
			s.Readings = readings
			return &s, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.sessions[id]; ok {
		s.Readings = r.readings[id]
		return s, nil
	}
	return nil, errors.New("session not found")
}

func (r *sessionRepository) GetSessionByBookingID(ctx context.Context, bookingID string) (*domain.ChargingSession, error) {
	if r.db != nil {
		query := `SELECT id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status, created_at, updated_at
		          FROM charging_sessions WHERE booking_id = $1`
		row := r.db.QueryRowContext(ctx, query, bookingID)

		var s domain.ChargingSession
		var endedAt sql.NullTime
		if err := row.Scan(&s.ID, &s.BookingID, &s.SlotID, &s.UserID, &s.StartedAt, &endedAt, &s.ConsumedKwh, &s.Status, &s.CreatedAt, &s.UpdatedAt); err == nil {
			if endedAt.Valid {
				s.EndedAt = &endedAt.Time
			}
			return &s, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.sessions {
		if s.BookingID == bookingID {
			s.Readings = r.readings[s.ID]
			return s, nil
		}
	}
	return nil, errors.New("session not found for booking")
}

func (r *sessionRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]domain.ChargingSession, error) {
	if r.db != nil {
		query := `SELECT id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status, created_at, updated_at
		          FROM charging_sessions WHERE user_id = $1 ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err == nil {
			defer rows.Close()
			var list []domain.ChargingSession
			for rows.Next() {
				var s domain.ChargingSession
				var endedAt sql.NullTime
				if err := rows.Scan(&s.ID, &s.BookingID, &s.SlotID, &s.UserID, &s.StartedAt, &endedAt, &s.ConsumedKwh, &s.Status, &s.CreatedAt, &s.UpdatedAt); err == nil {
					if endedAt.Valid {
						s.EndedAt = &endedAt.Time
					}
					list = append(list, s)
				}
			}
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.ChargingSession
	for _, s := range r.sessions {
		if s.UserID == userID {
			list = append(list, *s)
		}
	}
	return list, nil
}

func (r *sessionRepository) getMeterReadings(ctx context.Context, sessionID string) ([]domain.MeterReading, error) {
	if r.db != nil {
		query := `SELECT id, session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere
		          FROM meter_readings WHERE session_id = $1 ORDER BY recorded_at ASC`
		rows, err := r.db.QueryContext(ctx, query, sessionID)
		if err == nil {
			defer rows.Close()
			var list []domain.MeterReading
			for rows.Next() {
				var mr domain.MeterReading
				var v, a sql.NullFloat64
				if err := rows.Scan(&mr.ID, &mr.SessionID, &mr.RecordedAt, &mr.CurrentKwh, &mr.PowerKw, &v, &a); err == nil {
					mr.Voltage = v.Float64
					mr.CurrentAmpere = a.Float64
					list = append(list, mr)
				}
			}
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readings[sessionID], nil
}

func (r *sessionRepository) AddMeterReading(ctx context.Context, reading *domain.MeterReading) error {
	if r.db != nil {
		tx, err := r.db.BeginTx(ctx, nil)
		if err == nil {
			defer tx.Rollback()
			insertQuery := `INSERT INTO meter_readings (session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere)
			                VALUES ($1, $2, $3, $4, $5, $6)`
			tx.ExecContext(ctx, insertQuery, reading.SessionID, reading.RecordedAt, reading.CurrentKwh, reading.PowerKw, reading.Voltage, reading.CurrentAmpere)
			updateQuery := `UPDATE charging_sessions SET consumed_kwh = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
			tx.ExecContext(ctx, updateQuery, reading.CurrentKwh, reading.SessionID)
			if err := tx.Commit(); err == nil {
				return nil
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[reading.SessionID]; ok {
		s.ConsumedKwh = reading.CurrentKwh
		s.UpdatedAt = time.Now()
		r.readings[reading.SessionID] = append(r.readings[reading.SessionID], *reading)
		return nil
	}
	return errors.New("session not found")
}

func (r *sessionRepository) FinishSession(ctx context.Context, id string, endedAt time.Time, finalKwh float64) error {
	if r.db != nil {
		query := `UPDATE charging_sessions SET status = 'COMPLETED', ended_at = $1, consumed_kwh = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
		res, err := r.db.ExecContext(ctx, query, endedAt, finalKwh, id)
		if err == nil {
			if rows, _ := res.RowsAffected(); rows > 0 {
				return nil
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[id]; ok {
		s.Status = "COMPLETED"
		s.EndedAt = &endedAt
		s.ConsumedKwh = finalKwh
		s.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("session not found")
}

func (r *sessionRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
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
