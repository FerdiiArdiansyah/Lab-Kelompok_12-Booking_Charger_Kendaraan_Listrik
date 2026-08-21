package postgres

import (
	"context"
	"errors"
	"sync"
	"time"

	"station-service/internal/domain"
	"gorm.io/gorm"
)

type stationRepository struct {
	gormDB   *gorm.DB
	stations map[string]*domain.Station
	slots    map[string]*domain.ChargerSlot
	tariffs  map[string]*domain.Tariff
	mu       sync.RWMutex
}

func NewStationRepository(gormDB *gorm.DB) domain.StationRepository {
	repo := &stationRepository{
		gormDB:   gormDB,
		stations: make(map[string]*domain.Station),
		slots:    make(map[string]*domain.ChargerSlot),
		tariffs:  make(map[string]*domain.Tariff),
	}

	if gormDB != nil {
		// AutoMigrate database tables directly from domain models
		_ = gormDB.AutoMigrate(&domain.Station{}, &domain.ChargerSlot{}, &domain.Tariff{})
	}

	repo.seedInitialData()
	return repo
}

func (r *stationRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	sampleStations := []domain.Station{
		{
			ID:         "stn-001",
			Name:       "SPKLU PLN UID Jakarta Gambir",
			Location:   "Jl. M.I. Ridwan Rais No.1, Gambir, Jakarta Pusat",
			Latitude:   -6.1822,
			Longitude:  106.8344,
			TotalPower: 200.0,
			ImageURL:   "https://images.unsplash.com/photo-1558441719-443b38631157?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+PLN+UID+Jakarta+Gambir",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-002",
			Name:       "SPKLU PLN Senayan City Mall",
			Location:   "Jl. Asia Afrika No.19, Gelora, Tanah Abang, Jakarta Selatan",
			Latitude:   -6.2272,
			Longitude:  106.7974,
			TotalPower: 150.0,
			ImageURL:   "https://images.unsplash.com/photo-1593941707882-a5bba14938c7?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+Senayan+City+Mall",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-003",
			Name:       "SPKLU Rest Area KM 57A Tol Jakarta-Cikampek",
			Location:   "Tol Jakarta-Cikampek KM 57A, Karawang",
			Latitude:   -6.3683,
			Longitude:  107.3512,
			TotalPower: 200.0,
			ImageURL:   "https://images.unsplash.com/photo-1617788138017-80ad40651399?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+Rest+Area+KM+57A",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-004",
			Name:       "SPKLU PLN Rest Area KM 207A Tol Palimanan",
			Location:   "Tol Palikanci KM 207A, Cirebon",
			Latitude:   -6.7725,
			Longitude:  108.5361,
			TotalPower: 100.0,
			ImageURL:   "https://images.unsplash.com/photo-1542282088-72c9c27ed0cd?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+Rest+Area+KM+207A",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-005",
			Name:       "SPKLU PLN UID Jabar - Gedung Sate",
			Location:   "Jl. Asia Afrika No.63, Bandung",
			Latitude:   -6.9025,
			Longitude:  107.6186,
			TotalPower: 120.0,
			ImageURL:   "https://images.unsplash.com/photo-1571127236794-81c0bbfe1ce3?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+Gedung+Sate+Bandung",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-006",
			Name:       "SPKLU PLN UID Jawa Tengah Semarang",
			Location:   "Jl. Pemuda No.93, Semarang",
			Latitude:   -6.9822,
			Longitude:  110.4203,
			TotalPower: 150.0,
			ImageURL:   "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+PLN+UID+Semarang",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-007",
			Name:       "SPKLU PLN UID Jawa Timur Surabaya",
			Location:   "Jl. Yos Sudarso No.11, Surabaya",
			Latitude:   -7.2656,
			Longitude:  112.7483,
			TotalPower: 180.0,
			ImageURL:   "https://images.unsplash.com/photo-1558441719-443b38631157?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+PLN+UID+Surabaya",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-008",
			Name:       "SPKLU PLN UID Sulselrabar Hertasning Makassar",
			Location:   "Jl. Letjen Hertasning No.99, Kassi-Kassi, Kec. Rappocini, Kota Makassar, Sulawesi Selatan",
			Latitude:   -5.1678,
			Longitude:  119.4485,
			TotalPower: 200.0,
			ImageURL:   "https://images.unsplash.com/photo-1593941707882-a5bba14938c7?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+PLN+Hertasning+Makassar",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-009",
			Name:       "SPKLU PLN Mattoanging Makassar",
			Location:   "Jl. Andi Mappanyukki No.14, Kunjung Mae, Kec. Mariso, Kota Makassar, Sulawesi Selatan (PT PLN Persero Rayon Mattoanging)",
			Latitude:   -5.1553,
			Longitude:  119.4142,
			TotalPower: 50.0,
			ImageURL:   "/images/spklu-mattoanging.png",
			MapURL:     "https://maps.google.com/?q=PT+PLN+Persero+Rayon+Mattoanging",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-010",
			Name:       "SPKLU PLN UP3 Parepare",
			Location:   "Jl. Ahmad Yani No.51, Ujung Baru, Kec. Soreang, Kota Parepare, Sulawesi Selatan",
			Latitude:   -4.0152,
			Longitude:  119.6289,
			TotalPower: 50.0,
			ImageURL:   "https://images.unsplash.com/photo-1617788138017-80ad40651399?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+PLN+UP3+Parepare",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "stn-011",
			Name:       "SPKLU PLN UP3 Palopo",
			Location:   "Jl. Kelapa No.1, Dangerakko, Kec. Wara, Kota Palopo, Sulawesi Selatan",
			Latitude:   -2.9928,
			Longitude:  120.1983,
			TotalPower: 50.0,
			ImageURL:   "https://images.unsplash.com/photo-1558441719-443b38631157?q=80&w=800&auto=format&fit=crop",
			MapURL:     "https://maps.google.com/?q=SPKLU+PLN+UP3+Palopo",
			Status:     "ACTIVE",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	for i := range sampleStations {
		st := sampleStations[i]
		r.stations[st.ID] = &st
		if r.gormDB != nil {
			var existing domain.Station
			if err := r.gormDB.First(&existing, "id = ?", st.ID).Error; err != nil {
				r.gormDB.Create(&st)
			} else {
				r.gormDB.Model(&existing).Updates(domain.Station{
					Name:       st.Name,
					Location:   st.Location,
					TotalPower: st.TotalPower,
					ImageURL:   st.ImageURL,
					MapURL:     st.MapURL,
					Status:     st.Status,
				})
			}
		}
	}

	sampleSlots := []domain.ChargerSlot{
		{ID: "slot-001", StationID: "stn-001", SlotNumber: 1, ConnectorType: "CCS2", MaxPowerKw: 100.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-002", StationID: "stn-001", SlotNumber: 2, ConnectorType: "CHAdeMO", MaxPowerKw: 100.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-003", StationID: "stn-002", SlotNumber: 1, ConnectorType: "CCS2", MaxPowerKw: 150.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-004", StationID: "stn-003", SlotNumber: 1, ConnectorType: "CCS2", MaxPowerKw: 100.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-005", StationID: "stn-003", SlotNumber: 2, ConnectorType: "CHAdeMO", MaxPowerKw: 100.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-006", StationID: "stn-004", SlotNumber: 1, ConnectorType: "Type 2", MaxPowerKw: 22.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-007", StationID: "stn-005", SlotNumber: 1, ConnectorType: "CHAdeMO", MaxPowerKw: 50.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-008", StationID: "stn-006", SlotNumber: 1, ConnectorType: "Type 2", MaxPowerKw: 22.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-009", StationID: "stn-007", SlotNumber: 1, ConnectorType: "CHAdeMO", MaxPowerKw: 60.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-010", StationID: "stn-008", SlotNumber: 1, ConnectorType: "CCS2", MaxPowerKw: 100.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-011", StationID: "stn-008", SlotNumber: 2, ConnectorType: "CHAdeMO", MaxPowerKw: 100.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-012", StationID: "stn-009", SlotNumber: 1, ConnectorType: "Type 2", MaxPowerKw: 22.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-013", StationID: "stn-010", SlotNumber: 1, ConnectorType: "CHAdeMO", MaxPowerKw: 50.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
		{ID: "slot-014", StationID: "stn-011", SlotNumber: 1, ConnectorType: "CCS2", MaxPowerKw: 50.0, Status: "AVAILABLE", CreatedAt: now, UpdatedAt: now},
	}

	for i := range sampleSlots {
		sl := sampleSlots[i]
		r.slots[sl.ID] = &sl
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&sl, domain.ChargerSlot{ID: sl.ID})
		}
	}

	sampleTariffs := []domain.Tariff{
		{ID: "trf-001", StationID: "stn-001", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-002", StationID: "stn-002", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-003", StationID: "stn-003", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-004", StationID: "stn-004", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-005", StationID: "stn-005", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-006", StationID: "stn-006", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-007", StationID: "stn-007", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-008", StationID: "stn-008", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-009", StationID: "stn-009", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-010", StationID: "stn-010", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
		{ID: "trf-011", StationID: "stn-011", PricePerKwh: 2467.0, Currency: "IDR", ValidFrom: now, IsActive: true, CreatedAt: now},
	}

	for i := range sampleTariffs {
		tf := sampleTariffs[i]
		r.tariffs[tf.ID] = &tf
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&tf, domain.Tariff{ID: tf.ID})
		}
	}
}

func (r *stationRepository) GetAll(ctx context.Context) ([]domain.Station, error) {
	if r.gormDB != nil {
		var list []domain.Station
		if err := r.gormDB.WithContext(ctx).Preload("Slots").Preload("ActiveTariff").Find(&list).Error; err == nil {
			return list, nil
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
	if r.gormDB != nil {
		var s domain.Station
		if err := r.gormDB.WithContext(ctx).Preload("Slots").Preload("ActiveTariff").First(&s, "id = ?", id).Error; err == nil {
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

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(station).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.stations[station.ID] = station
	return nil
}

func (r *stationRepository) Update(ctx context.Context, station *domain.Station) error {
	station.UpdatedAt = time.Now()

	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.Station{}).Where("id = ?", station.ID).Updates(station)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.stations[station.ID]; ok {
		existing.Name = station.Name
		existing.Location = station.Location
		existing.Latitude = station.Latitude
		existing.Longitude = station.Longitude
		existing.TotalPower = station.TotalPower
		existing.ImageURL = station.ImageURL
		existing.MapURL = station.MapURL
		existing.Status = station.Status
		existing.UpdatedAt = station.UpdatedAt
		return nil
	}
	return errors.New("station not found")
}

func (r *stationRepository) Delete(ctx context.Context, id string) error {
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Where("id = ?", id).Delete(&domain.Station{})
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
	if r.gormDB != nil {
		var slots []domain.ChargerSlot
		if err := r.gormDB.WithContext(ctx).Where("station_id = ?", stationID).Order("slot_number ASC").Find(&slots).Error; err == nil {
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

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(slot).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.slots[slot.ID] = slot
	return nil
}

func (r *stationRepository) UpdateSlot(ctx context.Context, slot *domain.ChargerSlot) error {
	slot.UpdatedAt = time.Now()

	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.ChargerSlot{}).Where("id = ? AND station_id = ?", slot.ID, slot.StationID).Updates(slot)
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
	if r.gormDB != nil {
		var t domain.Tariff
		if err := r.gormDB.WithContext(ctx).Where("station_id = ? AND is_active = ?", stationID, true).Order("valid_from DESC").First(&t).Error; err == nil {
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

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(tariff).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.tariffs[tariff.ID] = tariff
	return nil
}

func (r *stationRepository) GetAllTariffs(ctx context.Context) ([]domain.Tariff, error) {
	if r.gormDB != nil {
		var list []domain.Tariff
		if err := r.gormDB.WithContext(ctx).Find(&list).Error; err == nil {
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
