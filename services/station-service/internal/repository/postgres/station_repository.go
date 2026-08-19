package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"station-service/internal/domain"
)

type stationRepository struct {
	db       *sql.DB
	stations map[string]*domain.Station
	slots    map[string]*domain.ChargerSlot
	tariffs  map[string]*domain.Tariff
	mu       sync.RWMutex
}

func NewStationRepository(db *sql.DB) domain.StationRepository {
	return &stationRepository{
		db:       db,
		stations: make(map[string]*domain.Station),
		slots:    make(map[string]*domain.ChargerSlot),
		tariffs:  make(map[string]*domain.Tariff),
	}
}

func (r *stationRepository) GetAll(ctx context.Context) ([]domain.Station, error) {
	if r.db != nil {
		query := `SELECT id, name, location, latitude, longitude, total_power_kw, status, created_at, updated_at FROM stations ORDER BY name ASC`
		rows, err := r.db.QueryContext(ctx, query)
		if err == nil {
			defer rows.Close()
			var list []domain.Station
			for rows.Next() {
				var s domain.Station
				if err := rows.Scan(&s.ID, &s.Name, &s.Location, &s.Latitude, &s.Longitude, &s.TotalPower, &s.Status, &s.CreatedAt, &s.UpdatedAt); err == nil {
					list = append(list, s)
				}
			}
			if len(list) > 0 {
				return list, nil
			}
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Station
	for _, s := range r.stations {
		list = append(list, *s)
	}
	return list, nil
}

func (r *stationRepository) GetByID(ctx context.Context, id string) (*domain.Station, error) {
	if r.db != nil {
		query := `SELECT id, name, location, latitude, longitude, total_power_kw, status, created_at, updated_at FROM stations WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)

		var s domain.Station
		if err := row.Scan(&s.ID, &s.Name, &s.Location, &s.Latitude, &s.Longitude, &s.TotalPower, &s.Status, &s.CreatedAt, &s.UpdatedAt); err == nil {
			slots, _ := r.GetSlotsByStationID(ctx, id)
			s.Slots = slots
			tariff, _ := r.GetTariffByStationID(ctx, id)
			s.ActiveTariff = tariff
			return &s, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.stations[id]; ok {
		var slots []domain.ChargerSlot
		for _, sl := range r.slots {
			if sl.StationID == id {
				slots = append(slots, *sl)
			}
		}
		s.Slots = slots
		for _, tr := range r.tariffs {
			if tr.StationID == id && tr.IsActive {
				s.ActiveTariff = tr
				break
			}
		}
		return s, nil
	}
	return nil, errors.New("station not found")
}

func (r *stationRepository) Create(ctx context.Context, station *domain.Station) error {
	station.CreatedAt = time.Now()
	station.UpdatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO stations (id, name, location, latitude, longitude, total_power_kw, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := r.db.ExecContext(ctx, query, station.ID, station.Name, station.Location, station.Latitude, station.Longitude, station.TotalPower, station.Status, station.CreatedAt, station.UpdatedAt)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.stations[station.ID] = station
	return nil
}

func (r *stationRepository) Update(ctx context.Context, station *domain.Station) error {
	station.UpdatedAt = time.Now()

	if r.db != nil {
		query := `UPDATE stations SET name = $1, location = $2, latitude = $3, longitude = $4, total_power_kw = $5, status = $6, updated_at = $7 WHERE id = $8`
		_, err := r.db.ExecContext(ctx, query, station.Name, station.Location, station.Latitude, station.Longitude, station.TotalPower, station.Status, station.UpdatedAt, station.ID)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.stations[station.ID]; ok {
		existing.Name = station.Name
		existing.Location = station.Location
		existing.Latitude = station.Latitude
		existing.Longitude = station.Longitude
		existing.TotalPower = station.TotalPower
		existing.Status = station.Status
		existing.UpdatedAt = station.UpdatedAt
		return nil
	}
	return errors.New("station not found")
}

func (r *stationRepository) Delete(ctx context.Context, id string) error {
	if r.db != nil {
		query := `DELETE FROM stations WHERE id = $1`
		_, err := r.db.ExecContext(ctx, query, id)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.stations[id]; ok {
		delete(r.stations, id)
		return nil
	}
	return errors.New("station not found")
}

func (r *stationRepository) GetSlotsByStationID(ctx context.Context, stationID string) ([]domain.ChargerSlot, error) {
	if r.db != nil {
		query := `SELECT id, station_id, slot_number, connector_type, max_power_kw, status, created_at, updated_at 
		          FROM charger_slots WHERE station_id = $1 ORDER BY slot_number ASC`
		rows, err := r.db.QueryContext(ctx, query, stationID)
		if err == nil {
			defer rows.Close()
			var slots []domain.ChargerSlot
			for rows.Next() {
				var cs domain.ChargerSlot
				if err := rows.Scan(&cs.ID, &cs.StationID, &cs.SlotNumber, &cs.ConnectorType, &cs.MaxPowerKw, &cs.Status, &cs.CreatedAt, &cs.UpdatedAt); err == nil {
					slots = append(slots, cs)
				}
			}
			return slots, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.ChargerSlot
	for _, sl := range r.slots {
		if sl.StationID == stationID {
			list = append(list, *sl)
		}
	}
	return list, nil
}

func (r *stationRepository) CreateSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	slot.CreatedAt = time.Now()
	slot.UpdatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO charger_slots (id, station_id, slot_number, connector_type, max_power_kw, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := r.db.ExecContext(ctx, query, slot.ID, slot.StationID, slot.SlotNumber, slot.ConnectorType, slot.MaxPowerKw, slot.Status, slot.CreatedAt, slot.UpdatedAt)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.slots[slot.ID] = slot
	return nil
}

func (r *stationRepository) UpdateSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	slot.UpdatedAt = time.Now()

	if r.db != nil {
		query := `UPDATE charger_slots SET connector_type = $1, max_power_kw = $2, status = $3, updated_at = $4 WHERE id = $5 AND station_id = $6`
		_, err := r.db.ExecContext(ctx, query, slot.ConnectorType, slot.MaxPowerKw, slot.Status, slot.UpdatedAt, slot.ID, slot.StationID)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.slots[slot.ID]; ok && existing.StationID == slot.StationID {
		existing.ConnectorType = slot.ConnectorType
		existing.MaxPowerKw = slot.MaxPowerKw
		existing.Status = slot.Status
		existing.UpdatedAt = slot.UpdatedAt
		return nil
	}
	return errors.New("slot not found")
}

func (r *stationRepository) GetTariffByStationID(ctx context.Context, stationID string) (*domain.Tariff, error) {
	if r.db != nil {
		query := `SELECT id, station_id, price_per_kwh, currency, valid_from, valid_to, is_active, created_at 
		          FROM tariffs WHERE station_id = $1 AND is_active = TRUE ORDER BY valid_from DESC LIMIT 1`
		row := r.db.QueryRowContext(ctx, query, stationID)
		var t domain.Tariff
		if err := row.Scan(&t.ID, &t.StationID, &t.PricePerKwh, &t.Currency, &t.ValidFrom, &t.ValidTo, &t.IsActive, &t.CreatedAt); err == nil {
			return &t, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, tr := range r.tariffs {
		if tr.StationID == stationID && tr.IsActive {
			return tr, nil
		}
	}
	return nil, errors.New("tariff not found")
}

func (r *stationRepository) CreateTariff(ctx context.Context, tariff *domain.Tariff) error {
	tariff.CreatedAt = time.Now()
	tariff.ValidFrom = time.Now()

	if r.db != nil {
		query := `INSERT INTO tariffs (id, station_id, price_per_kwh, currency, valid_from, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := r.db.ExecContext(ctx, query, tariff.ID, tariff.StationID, tariff.PricePerKwh, tariff.Currency, tariff.ValidFrom, tariff.IsActive, tariff.CreatedAt)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.tariffs[tariff.ID] = tariff
	return nil
}

func (r *stationRepository) GetAllTariffs(ctx context.Context) ([]domain.Tariff, error) {
	if r.db != nil {
		query := `SELECT id, station_id, price_per_kwh, currency, valid_from, valid_to, is_active, created_at FROM tariffs ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query)
		if err == nil {
			defer rows.Close()
			var list []domain.Tariff
			for rows.Next() {
				var t domain.Tariff
				if err := rows.Scan(&t.ID, &t.StationID, &t.PricePerKwh, &t.Currency, &t.ValidFrom, &t.ValidTo, &t.IsActive, &t.CreatedAt); err == nil {
					list = append(list, t)
				}
			}
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Tariff
	for _, tr := range r.tariffs {
		list = append(list, *tr)
	}
	return list, nil
}
