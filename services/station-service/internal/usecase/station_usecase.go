package usecase

import (
	"context"
	"errors"
	"station-service/internal/domain"

	"github.com/google/uuid"
)

type stationUsecase struct {
	repo domain.StationRepository
}

func NewStationUsecase(repo domain.StationRepository) domain.StationUsecase {
	return &stationUsecase{repo: repo}
}

func (u *stationUsecase) GetAllStations(ctx context.Context) ([]domain.Station, error) {
	return u.repo.GetAll(ctx)
}

func (u *stationUsecase) GetStationByID(ctx context.Context, id string) (*domain.Station, error) {
	if id == "" {
		return nil, errors.New("station ID is required")
	}
	return u.repo.GetByID(ctx, id)
}

func (u *stationUsecase) CreateStation(ctx context.Context, station *domain.Station) error {
	if station.Name == "" || station.Location == "" {
		return errors.New("name and location are required")
	}
	if station.ID == "" {
		station.ID = "stn-" + uuid.New().String()
	}
	if station.Status == "" {
		station.Status = "ACTIVE"
	}
	return u.repo.Create(ctx, station)
}

func (u *stationUsecase) GetSlots(ctx context.Context, stationID string) ([]domain.ChargerSlot, error) {
	if stationID == "" {
		return nil, errors.New("station ID is required")
	}
	return u.repo.GetSlotsByStationID(ctx, stationID)
}

func (u *stationUsecase) AddSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	if slot.StationID == "" || slot.ConnectorType == "" {
		return errors.New("station ID and connector type are required")
	}
	if slot.ID == "" {
		slot.ID = "slot-" + uuid.New().String()
	}
	if slot.Status == "" {
		slot.Status = "AVAILABLE"
	}
	return u.repo.CreateSlot(ctx, slot)
}

func (u *stationUsecase) GetTariff(ctx context.Context, stationID string) (*domain.Tariff, error) {
	if stationID == "" {
		return nil, errors.New("station ID is required")
	}
	return u.repo.GetTariffByStationID(ctx, stationID)
}

func (u *stationUsecase) AddTariff(ctx context.Context, tariff *domain.Tariff) error {
	if tariff.StationID == "" || tariff.PricePerKwh <= 0 {
		return errors.New("station ID and valid price per kwh are required")
	}
	if tariff.ID == "" {
		tariff.ID = "trf-" + uuid.New().String()
	}
	if tariff.Currency == "" {
		tariff.Currency = "IDR"
	}
	tariff.IsActive = true
	return u.repo.CreateTariff(ctx, tariff)
}
