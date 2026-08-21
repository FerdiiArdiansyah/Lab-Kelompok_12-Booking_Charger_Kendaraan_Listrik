package usecase

import (
	"booking-service/internal/domain"
	"context"
	"errors"
	"testing"
	"time"
)

// MockBookingRepository implements domain.BookingRepository for unit testing
type MockBookingRepository struct {
	existingBookings []domain.Booking
	waitlist         []domain.Waitlist
	outboxEvents     []map[string]interface{}
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, b *domain.Booking) error {
	avail, err := m.CheckSlotAvailability(ctx, b.SlotID, b.StartTime, b.EndTime)
	if err != nil {
		return err
	}
	if !avail {
		return errors.New("SLOT_OVERLAP_CONFLICT")
	}
	m.existingBookings = append(m.existingBookings, *b)
	return nil
}

func (m *MockBookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	for _, b := range m.existingBookings {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, errors.New("booking not found")
}

func (m *MockBookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]domain.Booking, error) {
	var result []domain.Booking
	for _, b := range m.existingBookings {
		if b.UserID == userID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockBookingRepository) GetAllBookings(ctx context.Context) ([]domain.Booking, error) {
	return m.existingBookings, nil
}

func (m *MockBookingRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	for i, b := range m.existingBookings {
		if b.ID == id {
			m.existingBookings[i].Status = status
			return nil
		}
	}
	return errors.New("booking not found")
}

func (m *MockBookingRepository) CheckSlotAvailability(ctx context.Context, slotID string, start, end time.Time) (bool, error) {
	for _, b := range m.existingBookings {
		if b.SlotID == slotID && (b.Status == "CONFIRMED" || b.Status == "IN_SESSION" || b.Status == "REQUESTED") {
			// Check overlap interval: [start, end) overlaps with [b.StartTime, b.EndTime)
			if start.Before(b.EndTime) && end.After(b.StartTime) {
				return false, nil
			}
		}
	}
	return true, nil
}

func (m *MockBookingRepository) AddToWaitlist(ctx context.Context, w *domain.Waitlist) error {
	w.QueueNumber = len(m.waitlist) + 1
	m.waitlist = append(m.waitlist, *w)
	return nil
}

func (m *MockBookingRepository) GetWaitlistByStation(ctx context.Context, stationID string) ([]domain.Waitlist, error) {
	var res []domain.Waitlist
	for _, w := range m.waitlist {
		if w.StationID == stationID {
			res = append(res, w)
		}
	}
	return res, nil
}

func (m *MockBookingRepository) AutoReleaseNoShowBookings(ctx context.Context, graceMinutes int) (int, error) {
	cutoff := time.Now().Add(-time.Duration(graceMinutes) * time.Minute)
	releasedCount := 0
	for i, b := range m.existingBookings {
		if b.Status == "CONFIRMED" && b.StartTime.Before(cutoff) {
			m.existingBookings[i].Status = "EXPIRED_NO_SHOW"
			releasedCount++
		}
	}
	return releasedCount, nil
}

func (m *MockBookingRepository) PromoteNextWaitlist(ctx context.Context, stationID string, slotID string) (*domain.Waitlist, error) {
	for i, w := range m.waitlist {
		if w.StationID == stationID && w.Status == "WAITING" {
			m.waitlist[i].Status = "PROMOTED"
			m.existingBookings = append(m.existingBookings, domain.Booking{
				ID:        "bkg-promoted-" + w.ID,
				UserID:    w.UserID,
				StationID: w.StationID,
				SlotID:    slotID,
				StartTime: w.RequestedStart,
				EndTime:   w.RequestedEnd,
				Status:    "CONFIRMED",
			})
			return &m.waitlist[i], nil
		}
	}
	return nil, nil
}

func (m *MockBookingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
	m.outboxEvents = append(m.outboxEvents, map[string]interface{}{
		"aggregate_type": aggregateType,
		"aggregate_id":   aggregateID,
		"event_type":     eventType,
		"payload":        payload,
	})
	return nil
}


// Tests for Anti-Overlap Engine

func TestCreateBooking_ValidTimeSlot_Success(t *testing.T) {
	repo := &MockBookingRepository{}
	uc := NewBookingUsecase(repo)
	ctx := context.Background()

	now := time.Now().Truncate(time.Hour)
	req := &domain.Booking{
		UserID:    "usr-001",
		StationID: "stn-001",
		SlotID:    "slt-001",
		StartTime: now.Add(1 * time.Hour),
		EndTime:   now.Add(2 * time.Hour),
	}

	result, err := uc.CreateBooking(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Status != "CONFIRMED" {
		t.Errorf("Expected status CONFIRMED, got %s", result.Status)
	}

	if len(repo.existingBookings) != 1 {
		t.Errorf("Expected 1 booking in repo, got %d", len(repo.existingBookings))
	}
}

func TestCreateBooking_OverlapWithExistingBooking_PutsOnWaitlist(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	existing := domain.Booking{
		ID:        "bkg-existing",
		UserID:    "usr-existing",
		StationID: "stn-001",
		SlotID:    "slt-001",
		StartTime: now.Add(10 * time.Hour),
		EndTime:   now.Add(12 * time.Hour),
		Status:    "CONFIRMED",
	}

	repo := &MockBookingRepository{
		existingBookings: []domain.Booking{existing},
	}
	uc := NewBookingUsecase(repo)
	ctx := context.Background()

	testCases := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
	}{
		{
			name:      "Exact Same Time Range",
			startTime: now.Add(10 * time.Hour),
			endTime:   now.Add(12 * time.Hour),
		},
		{
			name:      "Overlap Start Inside Existing",
			startTime: now.Add(11 * time.Hour),
			endTime:   now.Add(13 * time.Hour),
		},
		{
			name:      "Overlap End Inside Existing",
			startTime: now.Add(9 * time.Hour),
			endTime:   now.Add(11 * time.Hour),
		},
		{
			name:      "Enclosing Existing Booking",
			startTime: now.Add(9 * time.Hour),
			endTime:   now.Add(13 * time.Hour),
		},
		{
			name:      "Inner Subset Of Existing Booking",
			startTime: now.Add(10*time.Hour + 30*time.Minute),
			endTime:   now.Add(11 * time.Hour),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &domain.Booking{
				UserID:    "usr-new",
				StationID: "stn-001",
				SlotID:    "slt-001",
				StartTime: tc.startTime,
				EndTime:   tc.endTime,
			}

			res, err := uc.CreateBooking(ctx, req)
			if err == nil {
				t.Fatalf("Expected error due to slot overlap, got nil")
			}

			if res.Status != "WAITLISTED" {
				t.Errorf("Expected booking status WAITLISTED, got %s", res.Status)
			}

			waitlists, _ := repo.GetWaitlistByStation(ctx, "stn-001")
			if len(waitlists) == 0 {
				t.Errorf("Expected user to be added to waitlist queue")
			}
		})
	}
}

func TestCreateBooking_AdjacentTimeSlots_Success(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	existing := domain.Booking{
		ID:        "bkg-existing",
		UserID:    "usr-existing",
		StationID: "stn-001",
		SlotID:    "slt-001",
		StartTime: now.Add(10 * time.Hour),
		EndTime:   now.Add(12 * time.Hour),
		Status:    "CONFIRMED",
	}

	repo := &MockBookingRepository{
		existingBookings: []domain.Booking{existing},
	}
	uc := NewBookingUsecase(repo)
	ctx := context.Background()

	// Booking ending exactly when existing starts (8:00 - 10:00)
	t.Run("Adjacent Before", func(t *testing.T) {
		req := &domain.Booking{
			UserID:    "usr-adj-1",
			StationID: "stn-001",
			SlotID:    "slt-001",
			StartTime: now.Add(8 * time.Hour),
			EndTime:   now.Add(10 * time.Hour),
		}

		res, err := uc.CreateBooking(ctx, req)
		if err != nil {
			t.Fatalf("Adjacent before should not overlap, got error: %v", err)
		}
		if res.Status != "CONFIRMED" {
			t.Errorf("Expected CONFIRMED status, got %s", res.Status)
		}
	})

	// Booking starting exactly when existing ends (12:00 - 14:00)
	t.Run("Adjacent After", func(t *testing.T) {
		req := &domain.Booking{
			UserID:    "usr-adj-2",
			StationID: "stn-001",
			SlotID:    "slt-001",
			StartTime: now.Add(12 * time.Hour),
			EndTime:   now.Add(14 * time.Hour),
		}

		res, err := uc.CreateBooking(ctx, req)
		if err != nil {
			t.Fatalf("Adjacent after should not overlap, got error: %v", err)
		}
		if res.Status != "CONFIRMED" {
			t.Errorf("Expected CONFIRMED status, got %s", res.Status)
		}
	})
}

func TestCreateBooking_InvalidParameters_ReturnsError(t *testing.T) {
	repo := &MockBookingRepository{}
	uc := NewBookingUsecase(repo)
	ctx := context.Background()

	now := time.Now()

	t.Run("Missing Required Fields", func(t *testing.T) {
		req := &domain.Booking{
			UserID: "",
		}
		_, err := uc.CreateBooking(ctx, req)
		if err == nil {
			t.Errorf("Expected error for missing parameters")
		}
	})

	t.Run("EndTime Before StartTime", func(t *testing.T) {
		req := &domain.Booking{
			UserID:    "usr-1",
			StationID: "stn-1",
			SlotID:    "slt-1",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(1 * time.Hour),
		}
		_, err := uc.CreateBooking(ctx, req)
		if err == nil {
			t.Errorf("Expected error when EndTime is before StartTime")
		}
	})
}

// Tests for TICKET-02: Auto-Release No-Show Worker

func TestTriggerAutoRelease_ExpiredBookingsReleased(t *testing.T) {
	now := time.Now()
	repo := &MockBookingRepository{
		existingBookings: []domain.Booking{
			{
				ID:        "bkg-noshow-1",
				UserID:    "usr-1",
				StationID: "stn-1",
				SlotID:    "slt-1",
				StartTime: now.Add(-25 * time.Minute),
				EndTime:   now.Add(35 * time.Minute),
				Status:    "CONFIRMED",
			},
			{
				ID:        "bkg-valid-2",
				UserID:    "usr-2",
				StationID: "stn-1",
				SlotID:    "slt-2",
				StartTime: now.Add(-5 * time.Minute),
				EndTime:   now.Add(55 * time.Minute),
				Status:    "CONFIRMED",
			},
			{
				ID:        "bkg-insession-3",
				UserID:    "usr-3",
				StationID: "stn-1",
				SlotID:    "slt-3",
				StartTime: now.Add(-30 * time.Minute),
				EndTime:   now.Add(30 * time.Minute),
				Status:    "IN_SESSION",
			},
		},
	}

	uc := NewBookingUsecase(repo)
	ctx := context.Background()

	count, err := uc.TriggerAutoRelease(ctx, 15)
	if err != nil {
		t.Fatalf("Expected no error from TriggerAutoRelease, got: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 booking to be auto-released, got %d", count)
	}

	bkg1, _ := repo.GetBookingByID(ctx, "bkg-noshow-1")
	if bkg1.Status != "EXPIRED_NO_SHOW" {
		t.Errorf("Expected bkg-noshow-1 status EXPIRED_NO_SHOW, got %s", bkg1.Status)
	}

	bkg2, _ := repo.GetBookingByID(ctx, "bkg-valid-2")
	if bkg2.Status != "CONFIRMED" {
		t.Errorf("Expected bkg-valid-2 status CONFIRMED, got %s", bkg2.Status)
	}

	bkg3, _ := repo.GetBookingByID(ctx, "bkg-insession-3")
	if bkg3.Status != "IN_SESSION" {
		t.Errorf("Expected bkg-insession-3 status IN_SESSION, got %s", bkg3.Status)
	}
}

// Tests for TICKET-03: FIFO Waitlist Allocation Engine

func TestWaitlist_FIFO_Promotion_OnCancellation(t *testing.T) {
	now := time.Now()
	repo := &MockBookingRepository{
		existingBookings: []domain.Booking{
			{
				ID:        "bkg-to-cancel",
				UserID:    "usr-original",
				StationID: "stn-fifo-1",
				SlotID:    "slt-fifo-1",
				StartTime: now.Add(1 * time.Hour),
				EndTime:   now.Add(2 * time.Hour),
				Status:    "CONFIRMED",
			},
		},
		waitlist: []domain.Waitlist{
			{
				ID:             "wt-user-first",
				StationID:      "stn-fifo-1",
				UserID:         "usr-waitlist-1",
				RequestedStart: now.Add(1 * time.Hour),
				RequestedEnd:   now.Add(2 * time.Hour),
				QueueNumber:    1,
				Status:         "WAITING",
			},
			{
				ID:             "wt-user-second",
				StationID:      "stn-fifo-1",
				UserID:         "usr-waitlist-2",
				RequestedStart: now.Add(1 * time.Hour),
				RequestedEnd:   now.Add(2 * time.Hour),
				QueueNumber:    2,
				Status:         "WAITING",
			},
		},
	}

	uc := NewBookingUsecase(repo)
	ctx := context.Background()

	// Cancel booking
	err := uc.CancelBooking(ctx, "bkg-to-cancel")
	if err != nil {
		t.Fatalf("Expected no error when cancelling booking, got: %v", err)
	}

	// Verify original booking is CANCELLED
	cancelledBkg, _ := repo.GetBookingByID(ctx, "bkg-to-cancel")
	if cancelledBkg.Status != "CANCELLED" {
		t.Errorf("Expected original booking status CANCELLED, got %s", cancelledBkg.Status)
	}

	// Verify User #1 in Waitlist is PROMOTED
	waitlist, _ := repo.GetWaitlistByStation(ctx, "stn-fifo-1")
	foundPromoted := false
	for _, w := range waitlist {
		if w.ID == "wt-user-first" && w.Status == "PROMOTED" {
			foundPromoted = true
		}
	}
	if !foundPromoted {
		t.Errorf("Expected Queue #1 user (wt-user-first) to be PROMOTED")
	}

	// Verify Outbox Event for Promotion exists
	foundOutbox := false
	for _, evt := range repo.outboxEvents {
		if evt["event_type"] == "WaitlistPromoted" {
			foundOutbox = true
		}
	}
	if !foundOutbox {
		t.Errorf("Expected WaitlistPromoted event to be saved in Outbox")
	}
}


