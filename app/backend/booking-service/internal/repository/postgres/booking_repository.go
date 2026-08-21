package postgres

import (
	"booking-service/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"strings"
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

type bookingRepository struct {
	gormDB   *gorm.DB
	bookings map[string]*domain.Booking
	waitlist map[string]*domain.Waitlist
	mu       sync.RWMutex
}

func NewBookingRepository(gormDB *gorm.DB) domain.BookingRepository {
	repo := &bookingRepository{
		gormDB:   gormDB,
		bookings: make(map[string]*domain.Booking),
		waitlist: make(map[string]*domain.Waitlist),
	}

	if gormDB != nil {
		// AutoMigrate database tables directly from domain models
		_ = gormDB.AutoMigrate(&domain.Booking{}, &domain.Waitlist{}, &OutboxEventModel{})
	}

	repo.seedInitialData()
	return repo
}

func (r *bookingRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	sampleBookings := []domain.Booking{
		{ID: "bkg-001", UserID: "usr-001", StationID: "stn-001", SlotID: "slot-001", StartTime: now.Add(-2 * time.Hour), EndTime: now.Add(-1 * time.Hour), Status: "COMPLETED", IdempotencyKey: "idem-001", CreatedAt: now.Add(-3 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)},
		{ID: "bkg-002", UserID: "usr-002", StationID: "stn-002", SlotID: "slot-003", StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour), Status: "IN_SESSION", IdempotencyKey: "idem-002", CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)},
		{ID: "bkg-003", UserID: "usr-003", StationID: "stn-003", SlotID: "slot-004", StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-003", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-004", UserID: "usr-004", StationID: "stn-004", SlotID: "slot-006", StartTime: now.Add(2 * time.Hour), EndTime: now.Add(3 * time.Hour), Status: "CONFIRMED", IdempotencyKey: "idem-004", CreatedAt: now, UpdatedAt: now},
		{ID: "bkg-005", UserID: "usr-005", StationID: "stn-005", SlotID: "slot-007", StartTime: now.Add(-4 * time.Hour), EndTime: now.Add(-3 * time.Hour), Status: "COMPLETED", IdempotencyKey: "idem-005", CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: now.Add(-3 * time.Hour)},
	}

	for i := range sampleBookings {
		b := sampleBookings[i]
		r.bookings[b.ID] = &b
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&b, domain.Booking{ID: b.ID})
		}
	}
}

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(booking).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.bookings[booking.ID] = booking
	return nil
}

func (r *bookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	if r.gormDB != nil {
		var b domain.Booking
		if err := r.gormDB.WithContext(ctx).First(&b, "id = ?", id).Error; err == nil {
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
	if r.gormDB != nil {
		var list []domain.Booking
		if err := r.gormDB.WithContext(ctx).Where("user_id = ?", userID).Find(&list).Error; err == nil {
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

func (r *bookingRepository) GetAllBookings(ctx context.Context) ([]domain.Booking, error) {
	if r.gormDB != nil {
		var list []domain.Booking
		if err := r.gormDB.WithContext(ctx).Find(&list).Error; err == nil {
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Booking
	for _, b := range r.bookings {
		list = append(list, *b)
	}
	return list, nil
}

func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.Booking{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})
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
	if r.gormDB != nil {
		var count int64
		r.gormDB.WithContext(ctx).Model(&domain.Booking{}).
			Where("slot_id = ? AND status NOT IN ('CANCELLED', 'EXPIRED_NO_SHOW') AND start_time < ? AND end_time > ?", slotID, end, start).
			Count(&count)
		if count > 0 {
			return false, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.bookings {
		if b.SlotID == slotID && !strings.EqualFold(b.Status, "CANCELLED") && !strings.EqualFold(b.Status, "EXPIRED_NO_SHOW") {
			if start.Before(b.EndTime) && end.After(b.StartTime) {
				return false, nil
			}
		}
	}
	return true, nil
}

func (r *bookingRepository) AddToWaitlist(ctx context.Context, waitlist *domain.Waitlist) error {
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Create(waitlist)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.waitlist[waitlist.ID] = waitlist
	return nil
}

func (r *bookingRepository) GetWaitlistByStation(ctx context.Context, stationID string) ([]domain.Waitlist, error) {
	if r.gormDB != nil {
		var list []domain.Waitlist
		if err := r.gormDB.WithContext(ctx).Where("station_id = ? AND status = 'WAITING'", stationID).Order("queue_number ASC").Find(&list).Error; err == nil {
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

	if r.gormDB != nil {
		res := r.gormDB.WithContext(ctx).Model(&domain.Booking{}).
			Where("status = 'CONFIRMED' AND start_time <= ?", cutoff).
			Updates(map[string]interface{}{"status": "EXPIRED_NO_SHOW", "updated_at": time.Now()})
		if res.Error == nil {
			return int(res.RowsAffected), nil
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

func (r *bookingRepository) PromoteNextWaitlist(ctx context.Context, stationID string, slotID string) (*domain.Waitlist, error) {
	if r.gormDB != nil {
		var w domain.Waitlist
		err := r.gormDB.WithContext(ctx).Where("station_id = ? AND status = 'WAITING'", stationID).Order("queue_number ASC").First(&w).Error
		if err == nil {
			r.gormDB.WithContext(ctx).Model(&domain.Waitlist{}).Where("id = ?", w.ID).Updates(map[string]interface{}{"status": "PROMOTED", "updated_at": time.Now()})

			newBooking := &domain.Booking{
				ID:        "bkg-" + uuid.New().String(),
				UserID:    w.UserID,
				StationID: w.StationID,
				SlotID:    slotID,
				StartTime: w.RequestedStart,
				EndTime:   w.RequestedEnd,
				Status:    "CONFIRMED",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			r.gormDB.WithContext(ctx).Create(newBooking)

			w.Status = "PROMOTED"
			return &w, nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, w := range r.waitlist {
		if w.StationID == stationID && w.Status == "WAITING" {
			w.Status = "PROMOTED"
			w.UpdatedAt = time.Now()
			newBooking := &domain.Booking{
				ID:        "bkg-" + uuid.New().String(),
				UserID:    w.UserID,
				StationID: w.StationID,
				SlotID:    slotID,
				StartTime: w.RequestedStart,
				EndTime:   w.RequestedEnd,
				Status:    "CONFIRMED",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			r.bookings[newBooking.ID] = newBooking
			return w, nil
		}
	}
	return nil, nil
}

func (r *bookingRepository) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload interface{}) error {
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
