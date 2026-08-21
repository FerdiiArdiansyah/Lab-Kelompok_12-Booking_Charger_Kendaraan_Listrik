package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"session-service/internal/domain"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutboxEventModel struct {
	ID            string    `gorm:"primaryKey;size:64"`
	AggregateType string    `gorm:"index;size:64"`
	AggregateID   string    `gorm:"index;size:64"`
	EventType     string    `gorm:"index;size:64"`
	Payload       string    `gorm:"type:text"`
	Status        string    `gorm:"index;size:32;default:'PENDING'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (OutboxEventModel) TableName() string {
	return "outbox_events"
}

type sessionRepository struct {
	gormDB   *gorm.DB
	sessions map[string]*domain.ChargingSession
	readings map[string][]domain.MeterReading
	mu       sync.RWMutex
}

func NewSessionRepository(gormDB *gorm.DB) domain.SessionRepository {
	repo := &sessionRepository{
		gormDB:   gormDB,
		sessions: make(map[string]*domain.ChargingSession),
		readings: make(map[string][]domain.MeterReading),
	}

	if gormDB != nil {
		// AutoMigrate database tables directly from domain models
		_ = gormDB.AutoMigrate(&domain.ChargingSession{}, &domain.MeterReading{}, &OutboxEventModel{})
	}

	repo.seedInitialData()
	return repo
}

func (r *sessionRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

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
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&s, domain.ChargingSession{ID: s.ID})
		}

		readings := []domain.MeterReading{
			{ID: int64(i*10 + 1), SessionID: s.ID, RecordedAt: s.StartedAt, CurrentKwh: s.ConsumedKwh / 2, PowerKw: 50.0, Voltage: 400.0, CurrentAmpere: 125.0},
			{ID: int64(i*10 + 2), SessionID: s.ID, RecordedAt: s.StartedAt.Add(30 * time.Minute), CurrentKwh: s.ConsumedKwh, PowerKw: 50.0, Voltage: 400.0, CurrentAmpere: 125.0},
		}
		r.readings[s.ID] = readings

		if r.gormDB != nil {
			for _, rd := range readings {
				r.gormDB.FirstOrCreate(&rd, domain.MeterReading{ID: rd.ID})
			}
		}
	}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.ChargingSession) error {
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(session).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*domain.ChargingSession, error) {
	if r.gormDB != nil {
		var s domain.ChargingSession
		if err := r.gormDB.WithContext(ctx).Preload("Readings").First(&s, "id = ?", id).Error; err == nil {
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
	if r.gormDB != nil {
		var s domain.ChargingSession
		if err := r.gormDB.WithContext(ctx).Preload("Readings").First(&s, "booking_id = ?", bookingID).Error; err == nil {
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
	if r.gormDB != nil {
		var list []domain.ChargingSession
		if err := r.gormDB.WithContext(ctx).Preload("Readings").Where("user_id = ?", userID).Find(&list).Error; err == nil {
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

func (r *sessionRepository) AddMeterReading(ctx context.Context, reading *domain.MeterReading) error {
	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(reading).Error; err == nil {
			r.gormDB.WithContext(ctx).Model(&domain.ChargingSession{}).Where("id = ?", reading.SessionID).Updates(map[string]interface{}{
				"consumed_kwh": reading.CurrentKwh,
				"updated_at":   time.Now(),
			})
			return nil
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
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.ChargingSession{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       "COMPLETED",
			"ended_at":     endedAt,
			"consumed_kwh": finalKwh,
			"updated_at":   time.Now(),
		})
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
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if r.gormDB != nil {
		model := &OutboxEventModel{
			ID:            uuid.New().String(),
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			EventType:     eventType,
			Payload:       string(payloadBytes),
			Status:        "PENDING",
			CreatedAt:     time.Now(),
		}
		return r.gormDB.WithContext(ctx).Create(model).Error
	}

	return nil
}
