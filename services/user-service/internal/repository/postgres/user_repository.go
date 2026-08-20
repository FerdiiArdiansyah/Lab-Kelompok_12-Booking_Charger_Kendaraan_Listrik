package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
	"user-service/internal/domain"
)

type userRepository struct {
	db       *sql.DB
	users    map[string]*domain.User
	vehicles map[string]*domain.UserVehicle
	mu       sync.RWMutex
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	repo := &userRepository{
		db:       db,
		users:    make(map[string]*domain.User),
		vehicles: make(map[string]*domain.UserVehicle),
	}
	repo.seedInitialData()
	return repo
}

func (r *userRepository) seedInitialData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// 11 Seeded Users
	sampleUsers := []domain.User{
		{ID: "usr-001", Name: "Budi Santoso", Email: "budi.santoso@gmail.com", PasswordHash: "hashed-password123", Phone: "081234567890", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-002", Name: "Siti Aminah", Email: "siti.aminah@yahoo.com", PasswordHash: "hashed-password123", Phone: "081298765432", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-003", Name: "Farhan Ridwan", Email: "farhan.ridwan@outlook.com", PasswordHash: "hashed-password123", Phone: "081311223344", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-004", Name: "Dewi Lestari", Email: "dewi.lestari@gmail.com", PasswordHash: "hashed-password123", Phone: "081355667788", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-005", Name: "Agus Setiawan", Email: "agus.setiawan@gmail.com", PasswordHash: "hashed-password123", Phone: "081499887766", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-006", Name: "Rina Wijaya", Email: "rina.wijaya@gmail.com", PasswordHash: "hashed-password123", Phone: "081512341234", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-007", Name: "Hendra Pratama", Email: "hendra.pratama@gmail.com", PasswordHash: "hashed-password123", Phone: "081623452345", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-008", Name: "Maya Indah", Email: "maya.indah@gmail.com", PasswordHash: "hashed-password123", Phone: "081734563456", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-009", Name: "Eko Prasetyo", Email: "eko.prasetyo@gmail.com", PasswordHash: "hashed-password123", Phone: "081845674567", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-010", Name: "Ani Suryani", Email: "ani.suryani@gmail.com", PasswordHash: "hashed-password123", Phone: "081956785678", Role: "USER", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		{ID: "usr-011", Name: "Denny Kurniawan", Email: "denny.kurniawan@gmail.com", PasswordHash: "hashed-password123", Phone: "082167896789", Role: "ADMIN", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
	}

	for i := range sampleUsers {
		u := sampleUsers[i]
		r.users[u.ID] = &u
	}

	// 11 Seeded User Vehicles
	sampleVehicles := []domain.UserVehicle{
		{ID: "vhc-001", UserID: "usr-001", Brand: "Wuling", Model: "Air EV Long Range", LicensePlate: "B 1234 EV", ConnectorType: "Type 2", BatteryCapacityKwh: 26.7, CreatedAt: now},
		{ID: "vhc-002", UserID: "usr-002", Brand: "Hyundai", Model: "Ioniq 5 Signature Long Range", LicensePlate: "B 5678 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 72.6, CreatedAt: now},
		{ID: "vhc-003", UserID: "usr-003", Brand: "BYD", Model: "Atto 3 Extended", LicensePlate: "B 9101 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 60.48, CreatedAt: now},
		{ID: "vhc-004", UserID: "usr-004", Brand: "Tesla", Model: "Model 3 Performance", LicensePlate: "B 1122 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 82.0, CreatedAt: now},
		{ID: "vhc-005", UserID: "usr-005", Brand: "Chery", Model: "Omoda E5", LicensePlate: "B 3344 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 61.0, CreatedAt: now},
		{ID: "vhc-006", UserID: "usr-006", Brand: "MG", Model: "MG4 EV Magnify", LicensePlate: "B 5566 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 51.0, CreatedAt: now},
		{ID: "vhc-007", UserID: "usr-007", Brand: "Wuling", Model: "Binguo EV Premium", LicensePlate: "B 7788 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 37.9, CreatedAt: now},
		{ID: "vhc-008", UserID: "usr-008", Brand: "Neta", Model: "Neta V", LicensePlate: "B 9900 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 40.7, CreatedAt: now},
		{ID: "vhc-009", UserID: "usr-009", Brand: "Toyota", Model: "bZ4X AWD", LicensePlate: "B 1213 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 71.4, CreatedAt: now},
		{ID: "vhc-010", UserID: "usr-010", Brand: "Nissan", Model: "Leaf", LicensePlate: "B 1415 EV", ConnectorType: "CHAdeMO", BatteryCapacityKwh: 40.0, CreatedAt: now},
		{ID: "vhc-011", UserID: "usr-011", Brand: "BMW", Model: "i4 eDrive40", LicensePlate: "B 1617 EV", ConnectorType: "CCS2", BatteryCapacityKwh: 83.9, CreatedAt: now},
	}

	for i := range sampleVehicles {
		v := sampleVehicles[i]
		r.vehicles[v.ID] = &v
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if r.db != nil {
		query := `INSERT INTO users (id, name, email, password_hash, phone, role, status, created_at, updated_at) 
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Phone, user.Role, user.Status, user.CreatedAt, user.UpdatedAt)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if r.db != nil {
		query := `SELECT id, name, email, password_hash, phone, role, status, created_at, updated_at FROM users WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)
		var u domain.User
		if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt); err == nil {
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
	if r.db != nil {
		query := `SELECT id, name, email, password_hash, phone, role, status, created_at, updated_at FROM users WHERE email = $1`
		row := r.db.QueryRowContext(ctx, query, email)
		var u domain.User
		if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt); err == nil {
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
	if r.db != nil {
		query := `SELECT id, name, email, password_hash, phone, role, status, created_at, updated_at FROM users ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query)
		if err == nil {
			defer rows.Close()
			var list []domain.User
			for rows.Next() {
				var u domain.User
				if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt); err == nil {
					list = append(list, u)
				}
			}
			if len(list) > 0 {
				return list, nil
			}
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

	if r.db != nil {
		query := `UPDATE users SET name = $1, phone = $2, role = $3, status = $4, updated_at = $5 WHERE id = $6`
		res, err := r.db.ExecContext(ctx, query, user.Name, user.Phone, user.Role, user.Status, user.UpdatedAt, user.ID)
		if err == nil {
			if rows, _ := res.RowsAffected(); rows > 0 {
				return nil
			}
		}
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

	if r.db != nil {
		query := `INSERT INTO user_vehicles (id, user_id, brand, model, license_plate, connector_type, battery_capacity_kwh, created_at)
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := r.db.ExecContext(ctx, query, vehicle.ID, vehicle.UserID, vehicle.Brand, vehicle.Model, vehicle.LicensePlate, vehicle.ConnectorType, vehicle.BatteryCapacityKwh, vehicle.CreatedAt)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.vehicles[vehicle.ID] = vehicle
	return nil
}

func (r *userRepository) GetVehiclesByUserID(ctx context.Context, userID string) ([]domain.UserVehicle, error) {
	if r.db != nil {
		query := `SELECT id, user_id, brand, model, license_plate, connector_type, battery_capacity_kwh, created_at
		          FROM user_vehicles WHERE user_id = $1 ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err == nil {
			defer rows.Close()
			var list []domain.UserVehicle
			for rows.Next() {
				var v domain.UserVehicle
				if err := rows.Scan(&v.ID, &v.UserID, &v.Brand, &v.Model, &v.LicensePlate, &v.ConnectorType, &v.BatteryCapacityKwh, &v.CreatedAt); err == nil {
					list = append(list, v)
				}
			}
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
	if r.db != nil {
		query := `DELETE FROM user_vehicles WHERE id = $1 AND user_id = $2`
		res, err := r.db.ExecContext(ctx, query, vehicleID, userID)
		if err == nil {
			if rows, _ := res.RowsAffected(); rows > 0 {
				return nil
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.vehicles[vehicleID]; ok && v.UserID == userID {
		delete(r.vehicles, vehicleID)
		return nil
	}
	return errors.New("vehicle not found")
}
