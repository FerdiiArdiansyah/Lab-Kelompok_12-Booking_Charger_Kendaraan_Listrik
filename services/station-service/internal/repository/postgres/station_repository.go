package postgres

import (
	"context"
	"database/sql"
	"errors"
	"station-service/internal/domain"
)

type stationRepository struct {
	db *sql.DB
}

func NewStationRepository(db *sql.DB) domain.StationRepository {
	return &stationRepository{db: db}
}

func (r *stationRepository) GetAll(ctx context.Context) ([]domain.Station, error) {
	query := `SELECT id, name, location, latitude, longitude, total_power_kw, status, created_at, updated_at FROM stations ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stations []domain.Station
	for rows.Next() {
		var s domain.Station
		if err := rows.Scan(&s.ID, &s.Name, &s.Location, &s.Latitude, &s.Longitude, &s.TotalPower, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		stations = append(stations, s)
	}
	return stations, nil
}

func (r *stationRepository) GetByID(ctx context.Context, id string) (*domain.Station, error) {
	query := `SELECT id, name, location, latitude, longitude, total_power_kw, status, created_at, updated_at FROM stations WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.Station
	if err := row.Scan(&s.ID, &s.Name, &s.Location, &s.Latitude, &s.Longitude, &s.TotalPower, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("station not found")
		}
		return nil, err
	}

	slots, _ := r.GetSlotsByStationID(ctx, id)
	s.Slots = slots

	tariff, _ := r.GetTariffByStationID(ctx, id)
	s.ActiveTariff = tariff

	return &s, nil
}

func (r *stationRepository) Create(ctx context.Context, station *domain.Station) error {
	query := `INSERT INTO stations (id, name, location, latitude, longitude, total_power_kw, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, station.ID, station.Name, station.Location, station.Latitude, station.Longitude, station.TotalPower, station.Status)
	return err
}

func (r *stationRepository) GetSlotsByStationID(ctx context.Context, stationID string) ([]domain.ChargerSlot, error) {
	query := `SELECT id, station_id, slot_number, connector_type, max_power_kw, status, created_at, updated_at 
	          FROM charger_slots WHERE station_id = $1 ORDER BY slot_number ASC`
	rows, err := r.db.QueryContext(ctx, query, stationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []domain.ChargerSlot
	for rows.Next() {
		var cs domain.ChargerSlot
		if err := rows.Scan(&cs.ID, &cs.StationID, &cs.SlotNumber, &cs.ConnectorType, &cs.MaxPowerKw, &cs.Status, &cs.CreatedAt, &cs.UpdatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, cs)
	}
	return slots, nil
}

func (r *stationRepository) CreateSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	query := `INSERT INTO charger_slots (id, station_id, slot_number, connector_type, max_power_kw, status) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, slot.ID, slot.StationID, slot.SlotNumber, slot.ConnectorType, slot.MaxPowerKw, slot.Status)
	return err
}

func (r *stationRepository) GetTariffByStationID(ctx context.Context, stationID string) (*domain.Tariff, error) {
	query := `SELECT id, station_id, price_per_kwh, currency, valid_from, valid_to, is_active, created_at 
	          FROM tariffs WHERE station_id = $1 AND is_active = TRUE ORDER BY valid_from DESC LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, stationID)

	var t domain.Tariff
	if err := row.Scan(&t.ID, &t.StationID, &t.PricePerKwh, &t.Currency, &t.ValidFrom, &t.ValidTo, &t.IsActive, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tariff not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *stationRepository) CreateTariff(ctx context.Context, tariff *domain.Tariff) error {
	query := `INSERT INTO tariffs (id, station_id, price_per_kwh, currency, valid_from, is_active) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, tariff.ID, tariff.StationID, tariff.PricePerKwh, tariff.Currency, tariff.ValidFrom, tariff.IsActive)
	return err
}
