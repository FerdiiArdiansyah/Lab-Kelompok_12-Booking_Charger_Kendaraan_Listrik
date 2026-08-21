package postgres

import (
	"context"
	"errors"
	"sync"
	"time"
	"user-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	gormDB   *gorm.DB
	users    map[string]*domain.User
	vehicles map[string]*domain.UserVehicle
	mu       sync.RWMutex
}

func NewUserRepository(gormDB *gorm.DB) domain.UserRepository {
	repo := &userRepository{
		gormDB:   gormDB,
		users:    make(map[string]*domain.User),
		vehicles: make(map[string]*domain.UserVehicle),
	}

	if gormDB != nil {
		// AutoMigrate tables directly from domain models
		_ = gormDB.AutoMigrate(&domain.User{}, &domain.UserVehicle{})
	}

	repo.seedInitialData()
	return repo
}

func (r *userRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	hashPassword := func(pwd string) string {
		h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return pwd
		}
		return string(h)
	}

	sampleUsers := []domain.User{
		{ID: "usr-admin", Name: "Administrator VoltHub", Email: "admin@spklu.co.id", PasswordHash: hashPassword("admin123"), Phone: "081100009999", Role: "ADMIN", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-operator", Name: "Operator SPKLU PLN", Email: "operator@spklu.co.id", PasswordHash: hashPassword("operator123"), Phone: "081122223333", Role: "OPERATOR", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-driver", Name: "Ferdi Ardiansyah", Email: "ferdi@gmail.com", PasswordHash: hashPassword("user123"), Phone: "081234567890", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-001", Name: "Budi Santoso", Email: "budi.santoso@gmail.com", PasswordHash: hashPassword("password123"), Phone: "081234567890", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-002", Name: "Siti Aminah", Email: "siti.aminah@yahoo.com", PasswordHash: hashPassword("password123"), Phone: "081298765432", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-003", Name: "Farhan Ridwan", Email: "farhan.ridwan@outlook.com", PasswordHash: hashPassword("password123"), Phone: "081311223344", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
	}

	for i := range sampleUsers {
		u := sampleUsers[i]
		r.users[u.ID] = &u
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&u, domain.User{ID: u.ID})
		}
	}

	sampleVehicles := []domain.UserVehicle{
		{ID: "vhc-001", UserID: "usr-driver", Brand: "BMW", Model: "iX xDrive50", LicensePlate: "B 8888 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 105.2, CreatedAt: now},
		{ID: "vhc-002", UserID: "usr-driver", Brand: "Hyundai", Model: "Ioniq 5 Long Range", LicensePlate: "B 5678 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 72.6, CreatedAt: now},
	}

	for i := range sampleVehicles {
		v := sampleVehicles[i]
		r.vehicles[v.ID] = &v
		if r.gormDB != nil {
			r.gormDB.FirstOrCreate(&v, domain.UserVehicle{ID: v.ID})
		}
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if r.gormDB != nil {
		var u domain.User
		if err := r.gormDB.WithContext(ctx).Preload("Vehicles").First(&u, "id = ?", id).Error; err == nil {
			return &u, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if r.gormDB != nil {
		var u domain.User
		if err := r.gormDB.WithContext(ctx).Preload("Vehicles").First(&u, "email = ?", email).Error; err == nil {
			return &u, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	if r.gormDB != nil {
		var list []domain.User
		if err := r.gormDB.WithContext(ctx).Preload("Vehicles").Find(&list).Error; err == nil {
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.User
	for _, u := range r.users {
		list = append(list, *u)
	}
	return list, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(user)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.users[user.ID]; ok {
		existing.Name = user.Name
		existing.Phone = user.Phone
		existing.Role = user.Role
		existing.Status = user.Status
		existing.UpdatedAt = user.UpdatedAt
		return nil
	}
	return errors.New("user not found")
}

func (r *userRepository) CreateVehicle(ctx context.Context, vehicle *domain.UserVehicle) error {
	vehicle.CreatedAt = time.Now()

	if r.gormDB != nil {
		if err := r.gormDB.WithContext(ctx).Create(vehicle).Error; err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.vehicles[vehicle.ID] = vehicle
	return nil
}

func (r *userRepository) GetVehiclesByUserID(ctx context.Context, userID string) ([]domain.UserVehicle, error) {
	if r.gormDB != nil {
		var list []domain.UserVehicle
		if err := r.gormDB.WithContext(ctx).Where("user_id = ?", userID).Find(&list).Error; err == nil {
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.UserVehicle
	for _, v := range r.vehicles {
		if v.UserID == userID {
			list = append(list, *v)
		}
	}
	return list, nil
}

func (r *userRepository) DeleteVehicle(ctx context.Context, vehicleID, userID string) error {
	if r.gormDB != nil {
		r.gormDB.WithContext(ctx).Where("id = ? AND user_id = ?", vehicleID, userID).Delete(&domain.UserVehicle{})
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.vehicles[vehicleID]; ok && v.UserID == userID {
		delete(r.vehicles, vehicleID)
		return nil
	}
	return errors.New("vehicle not found")
}
