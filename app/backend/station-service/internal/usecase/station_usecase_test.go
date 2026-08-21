package usecase

import (
	"context"
	"errors"
	"station-service/internal/domain"
	"testing"
)

type MockStationRepository struct {
	stations map[string]*domain.Station
	slots    map[string][]domain.ChargerSlot
	tariffs  map[string]*domain.Tariff
}

func NewMockStationRepository() *MockStationRepository {
	return &MockStationRepository{
		stations: make(map[string]*domain.Station),
		slots:    make(map[string][]domain.ChargerSlot),
		tariffs:  make(map[string]*domain.Tariff),
	}
}

func (m *MockStationRepository) GetAll(ctx context.Context) ([]domain.Station, error) {
	var list []domain.Station
	for _, s := range m.stations {
		list = append(list, *s)
	}
	return list, nil
}

func (m *MockStationRepository) GetByID(ctx context.Context, id string) (*domain.Station, error) {
	if s, ok := m.stations[id]; ok {
		return s, nil
	}
	return nil, errors.New("station not found")
}

func (m *MockStationRepository) Create(ctx context.Context, station *domain.Station) error {
	m.stations[station.ID] = station
	return nil
}

func (m *MockStationRepository) Update(ctx context.Context, station *domain.Station) error {
	if _, ok := m.stations[station.ID]; ok {
		m.stations[station.ID] = station
		return nil
	}
	return errors.New("station not found")
}

func (m *MockStationRepository) Delete(ctx context.Context, id string) error {
	delete(m.stations, id)
	return nil
}

func (m *MockStationRepository) GetSlotsByStationID(ctx context.Context, stationID string) ([]domain.ChargerSlot, error) {
	return m.slots[stationID], nil
}

func (m *MockStationRepository) CreateSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	m.slots[slot.StationID] = append(m.slots[slot.StationID], *slot)
	return nil
}

func (m *MockStationRepository) UpdateSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	return nil
}

func (m *MockStationRepository) GetTariffByStationID(ctx context.Context, stationID string) (*domain.Tariff, error) {
	if t, ok := m.tariffs[stationID]; ok {
		return t, nil
	}
	return nil, errors.New("tariff not found")
}

func (m *MockStationRepository) CreateTariff(ctx context.Context, tariff *domain.Tariff) error {
	m.tariffs[tariff.StationID] = tariff
	return nil
}

func (m *MockStationRepository) GetAllTariffs(ctx context.Context) ([]domain.Tariff, error) {
	var list []domain.Tariff
	for _, t := range m.tariffs {
		list = append(list, *t)
	}
	return list, nil
}

// Unit Tests for TICKET-07 (Station & Charger Slot Connector Compatibility)

func TestCreateStation_Success(t *testing.T) {
	repo := NewMockStationRepository()
	uc := NewStationUsecase(repo)
	ctx := context.Background()

	stn := &domain.Station{
		Name:       "SPKLU PLN Gambir",
		Location:   "Jakarta Pusat",
		TotalPower: 200.0,
	}

	err := uc.CreateStation(ctx, stn)
	if err != nil {
		t.Fatalf("Expected no error creating station, got: %v", err)
	}

	if stn.Status != "ACTIVE" {
		t.Errorf("Expected status ACTIVE, got %s", stn.Status)
	}
}

func TestAddSlot_ValidConnectorType_Success(t *testing.T) {
	repo := NewMockStationRepository()
	uc := NewStationUsecase(repo)
	ctx := context.Background()

	slot := &domain.ChargerSlot{
		StationID:     "stn-1",
		SlotNumber:    1,
		ConnectorType: "CCS2",
		MaxPowerKw:    150.0,
	}

	err := uc.AddSlot(ctx, slot)
	if err != nil {
		t.Fatalf("Expected no error adding slot with CCS2 connector, got: %v", err)
	}

	slots, _ := repo.GetSlotsByStationID(ctx, "stn-1")
	if len(slots) != 1 {
		t.Fatalf("Expected 1 slot added, got %d", len(slots))
	}
}

func TestAddSlot_InvalidConnectorType_ReturnsError(t *testing.T) {
	repo := NewMockStationRepository()
	uc := NewStationUsecase(repo)
	ctx := context.Background()

	slot := &domain.ChargerSlot{
		StationID:     "stn-1",
		SlotNumber:    2,
		ConnectorType: "INVALID_CONNECTOR",
		MaxPowerKw:    50.0,
	}

	err := uc.AddSlot(ctx, slot)
	if err == nil {
		t.Errorf("Expected error when adding slot with invalid connector type, got nil")
	}
}

