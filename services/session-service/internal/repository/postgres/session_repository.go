package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"session-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.ChargingSession) error {
	query := `INSERT INTO charging_sessions (id, booking_id, slot_id, user_id, started_at, consumed_kwh, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.BookingID, session.SlotID, session.UserID,
		session.StartedAt, session.ConsumedKwh, session.Status)
	return err
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*domain.ChargingSession, error) {
	query := `SELECT id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status, created_at, updated_at
	          FROM charging_sessions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.ChargingSession
	var endedAt sql.NullTime
	if err := row.Scan(&s.ID, &s.BookingID, &s.SlotID, &s.UserID, &s.StartedAt, &endedAt, &s.ConsumedKwh, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	if endedAt.Valid {
		s.EndedAt = &endedAt.Time
	}

	readings, _ := r.getMeterReadings(ctx, id)
	s.Readings = readings

	return &s, nil
}

func (r *sessionRepository) GetSessionByBookingID(ctx context.Context, bookingID string) (*domain.ChargingSession, error) {
	query := `SELECT id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status, created_at, updated_at
	          FROM charging_sessions WHERE booking_id = $1`
	row := r.db.QueryRowContext(ctx, query, bookingID)

	var s domain.ChargingSession
	var endedAt sql.NullTime
	if err := row.Scan(&s.ID, &s.BookingID, &s.SlotID, &s.UserID, &s.StartedAt, &endedAt, &s.ConsumedKwh, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("session not found for booking")
		}
		return nil, err
	}
	if endedAt.Valid {
		s.EndedAt = &endedAt.Time
	}
	return &s, nil
}

func (r *sessionRepository) getMeterReadings(ctx context.Context, sessionID string) ([]domain.MeterReading, error) {
	query := `SELECT id, session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere
	          FROM meter_readings WHERE session_id = $1 ORDER BY recorded_at ASC`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.MeterReading
	for rows.Next() {
		var mr domain.MeterReading
		var v, a sql.NullFloat64
		if err := rows.Scan(&mr.ID, &mr.SessionID, &mr.RecordedAt, &mr.CurrentKwh, &mr.PowerKw, &v, &a); err != nil {
			return nil, err
		}
		mr.Voltage = v.Float64
		mr.CurrentAmpere = a.Float64
		list = append(list, mr)
	}
	return list, nil
}

func (r *sessionRepository) AddMeterReading(ctx context.Context, reading *domain.MeterReading) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertQuery := `INSERT INTO meter_readings (session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere)
	                VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, insertQuery, reading.SessionID, reading.RecordedAt, reading.CurrentKwh, reading.PowerKw, reading.Voltage, reading.CurrentAmpere)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE charging_sessions SET consumed_kwh = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, reading.CurrentKwh, reading.SessionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *sessionRepository) FinishSession(ctx context.Context, id string, endedAt time.Time, finalKwh float64) error {
	query := `UPDATE charging_sessions SET status = 'COMPLETED', ended_at = $1, consumed_kwh = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	res, err := r.db.ExecContext(ctx, query, endedAt, finalKwh, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("session not found")
	}
	return nil
}

func (r *sessionRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	query := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status)
	          VALUES ($1, $2, $3, $4, $5, 'PENDING')`
	_, err = r.db.ExecContext(ctx, query, uuid.New().String(), aggregateType, aggregateID, eventType, payloadBytes)
	return err
}
