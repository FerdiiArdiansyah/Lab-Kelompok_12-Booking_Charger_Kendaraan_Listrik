package usecase

import (
	"context"
	"errors"
	"session-service/internal/domain"
	"testing"
	"time"
)

type MockSessionRepository struct {
	sessions     map[string]*domain.ChargingSession
	readings     map[string][]domain.MeterReading
	outboxEvents []map[string]interface{}
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions:     make(map[string]*domain.ChargingSession),
		readings:     make(map[string][]domain.MeterReading),
		outboxEvents: make([]map[string]interface{}, 0),
	}
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, session *domain.ChargingSession) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *MockSessionRepository) GetSessionByID(ctx context.Context, id string) (*domain.ChargingSession, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, errors.New("session not found")
}

func (m *MockSessionRepository) GetSessionByBookingID(ctx context.Context, bookingID string) (*domain.ChargingSession, error) {
	for _, s := range m.sessions {
		if s.BookingID == bookingID {
			return s, nil
		}
	}
	return nil, errors.New("session not found for booking")
}

func (m *MockSessionRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]domain.ChargingSession, error) {
	var list []domain.ChargingSession
	for _, s := range m.sessions {
		if s.UserID == userID {
			list = append(list, *s)
		}
	}
	return list, nil
}

func (m *MockSessionRepository) AddMeterReading(ctx context.Context, reading *domain.MeterReading) error {
	m.readings[reading.SessionID] = append(m.readings[reading.SessionID], *reading)
	return nil
}

func (m *MockSessionRepository) FinishSession(ctx context.Context, id string, endedAt time.Time, finalKwh float64) error {
	if s, ok := m.sessions[id]; ok {
		s.Status = "COMPLETED"
		s.EndedAt = &endedAt
		s.ConsumedKwh = finalKwh
		return nil
	}
	return errors.New("session not found")
}

func (m *MockSessionRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	m.outboxEvents = append(m.outboxEvents, map[string]interface{}{
		"aggregate_type": aggregateType,
		"aggregate_id":   aggregateID,
		"event_type":     eventType,
		"payload":        payload,
	})
	return nil
}

// Unit Tests for TICKET-04

func TestStartSession_ValidParameters_Success(t *testing.T) {
	repo := NewMockSessionRepository()
	uc := NewSessionUsecase(repo)
	ctx := context.Background()

	session, err := uc.StartSession(ctx, "bkg-100", "slt-100", "usr-100")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if session.Status != "IN_PROGRESS" {
		t.Errorf("Expected status IN_PROGRESS, got %s", session.Status)
	}

	if len(repo.outboxEvents) != 1 {
		t.Errorf("Expected 1 outbox event SessionStarted, got %d", len(repo.outboxEvents))
	}
}

func TestRecordMeter_ActiveSession_Success(t *testing.T) {
	repo := NewMockSessionRepository()
	uc := NewSessionUsecase(repo)
	ctx := context.Background()

	session, _ := uc.StartSession(ctx, "bkg-101", "slt-101", "usr-101")

	err := uc.RecordMeter(ctx, session.ID, 12.5, 22.0, 380.0, 32.0)
	if err != nil {
		t.Fatalf("Expected no error recording meter, got: %v", err)
	}

	readings := repo.readings[session.ID]
	if len(readings) != 1 {
		t.Fatalf("Expected 1 meter reading recorded, got %d", len(readings))
	}

	if readings[0].CurrentKwh != 12.5 {
		t.Errorf("Expected CurrentKwh 12.5, got %f", readings[0].CurrentKwh)
	}
}

func TestFinishSession_CalculatesKwhAndPublishesEvent_Success(t *testing.T) {
	repo := NewMockSessionRepository()
	uc := NewSessionUsecase(repo)
	ctx := context.Background()

	session, _ := uc.StartSession(ctx, "bkg-102", "slt-102", "usr-102")

	// Record several telemetry meters
	_ = uc.RecordMeter(ctx, session.ID, 5.0, 22.0, 380.0, 32.0)
	_ = uc.RecordMeter(ctx, session.ID, 15.4, 22.0, 380.0, 32.0)

	// Finish session with 15.4 kWh final consumption
	finishedSession, err := uc.FinishSession(ctx, session.ID, 15.4)
	if err != nil {
		t.Fatalf("Expected no error finishing session, got: %v", err)
	}

	if finishedSession.Status != "COMPLETED" {
		t.Errorf("Expected status COMPLETED, got %s", finishedSession.Status)
	}

	if finishedSession.ConsumedKwh != 15.4 {
		t.Errorf("Expected ConsumedKwh 15.4, got %f", finishedSession.ConsumedKwh)
	}

	// Verify Outbox Event for SessionFinished
	foundFinishedEvent := false
	for _, evt := range repo.outboxEvents {
		if evt["event_type"] == "SessionFinished" {
			foundFinishedEvent = true
			payload := evt["payload"].(map[string]interface{})
			if payload["consumed_kwh"] != 15.4 {
				t.Errorf("Expected event payload consumed_kwh 15.4, got %v", payload["consumed_kwh"])
			}
		}
	}
	if !foundFinishedEvent {
		t.Errorf("Expected SessionFinished event to be saved in outbox")
	}
}

func TestFinishSession_AlreadyFinished_ReturnsError(t *testing.T) {
	repo := NewMockSessionRepository()
	uc := NewSessionUsecase(repo)
	ctx := context.Background()

	session, _ := uc.StartSession(ctx, "bkg-103", "slt-103", "usr-103")

	// Finish session first time
	_, _ = uc.FinishSession(ctx, session.ID, 20.0)

	// Try to finish session second time
	_, err := uc.FinishSession(ctx, session.ID, 25.0)
	if err == nil {
		t.Errorf("Expected error when trying to finish an already finished session")
	}
}
