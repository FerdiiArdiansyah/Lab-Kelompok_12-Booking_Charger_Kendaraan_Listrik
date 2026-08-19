package usecase

import (
	"context"
	"errors"
	"session-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type sessionUsecase struct {
	repo domain.SessionRepository
}

func NewSessionUsecase(repo domain.SessionRepository) domain.SessionUsecase {
	return &sessionUsecase{repo: repo}
}

func (u *sessionUsecase) StartSession(ctx context.Context, bookingID, slotID, userID string) (*domain.ChargingSession, error) {
	if bookingID == "" || slotID == "" || userID == "" {
		return nil, errors.New("booking_id, slot_id, and user_id are required")
	}

	session := &domain.ChargingSession{
		ID:          "ses-" + uuid.New().String(),
		BookingID:   bookingID,
		SlotID:      slotID,
		UserID:      userID,
		StartedAt:   time.Now(),
		ConsumedKwh: 0.0,
		Status:      "IN_PROGRESS",
	}

	if err := u.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	_ = u.repo.SaveOutboxEvent(ctx, "ChargingSession", session.ID, "SessionStarted", map[string]interface{}{
		"session_id": session.ID,
		"booking_id": bookingID,
		"started_at": session.StartedAt,
	})

	return session, nil
}

func (u *sessionUsecase) GetSessionByID(ctx context.Context, id string) (*domain.ChargingSession, error) {
	if id == "" {
		return nil, errors.New("session ID is required")
	}
	return u.repo.GetSessionByID(ctx, id)
}

func (u *sessionUsecase) RecordMeter(ctx context.Context, sessionID string, currentKwh, powerKw, voltage, ampere float64) error {
	session, err := u.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.Status != "IN_PROGRESS" {
		return errors.New("cannot record meter readings on non-active session")
	}

	reading := &domain.MeterReading{
		SessionID:     sessionID,
		RecordedAt:    time.Now(),
		CurrentKwh:    currentKwh,
		PowerKw:       powerKw,
		Voltage:       voltage,
		CurrentAmpere: ampere,
	}

	return u.repo.AddMeterReading(ctx, reading)
}

func (u *sessionUsecase) FinishSession(ctx context.Context, sessionID string, finalKwh float64) (*domain.ChargingSession, error) {
	session, err := u.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.Status != "IN_PROGRESS" {
		return nil, errors.New("session is not currently IN_PROGRESS")
	}

	now := time.Now()
	if err := u.repo.FinishSession(ctx, sessionID, now, finalKwh); err != nil {
		return nil, err
	}

	session.Status = "COMPLETED"
	session.EndedAt = &now
	session.ConsumedKwh = finalKwh

	_ = u.repo.SaveOutboxEvent(ctx, "ChargingSession", sessionID, "SessionFinished", map[string]interface{}{
		"session_id":   sessionID,
		"booking_id":   session.BookingID,
		"user_id":      session.UserID,
		"consumed_kwh": finalKwh,
		"ended_at":     now,
	})

	return session, nil
}
