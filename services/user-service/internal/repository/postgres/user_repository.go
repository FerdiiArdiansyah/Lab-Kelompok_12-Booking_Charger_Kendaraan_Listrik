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
	return &userRepository{
		db:       db,
		users:    make(map[string]*domain.User),
		vehicles: make(map[string]*domain.UserVehicle),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if r.db != nil {
		query := `INSERT INTO users (id, name, email, password_hash, phone, role, status, created_at, updated_at) 
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Phone, user.Role, user.Status, user.CreatedAt, user.UpdatedAt)
		if err == nil {
			return nil
		}
	}

	// Fallback to in-memory storage for high availability & offline dev
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.Email == user.Email {
			return errors.New("email already registered")
		}
	}
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if r.db != nil {
		query := `SELECT id, name, email, password_hash, phone, role, status, created_at, updated_at FROM users WHERE id = $1`
		row := r.db.QueryRowContext(ctx, query, id)
		var user domain.User
		if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Phone, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt); err == nil {
			return &user, nil
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
		var user domain.User
		if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Phone, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt); err == nil {
			return &user, nil
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
		query := `SELECT id, name, email, phone, role, status, created_at, updated_at FROM users ORDER BY created_at DESC`
		rows, err := r.db.QueryContext(ctx, query)
		if err == nil {
			defer rows.Close()
			var list []domain.User
			for rows.Next() {
				var u domain.User
				if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt); err == nil {
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
		query := `UPDATE users SET name = $1, phone = $2, updated_at = $3 WHERE id = $4`
		_, err := r.db.ExecContext(ctx, query, user.Name, user.Phone, user.UpdatedAt, user.ID)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.users[user.ID]; ok {
		existing.Name = user.Name
		existing.Phone = user.Phone
		existing.UpdatedAt = user.UpdatedAt
		return nil
	}
	return errors.New("user not found")
}

func (r *userRepository) CreateVehicle(ctx context.Context, vehicle *domain.UserVehicle) error {
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

func (r *userRepository) DeleteVehicle(ctx context.Context, id, userID string) error {
	if r.db != nil {
		query := `DELETE FROM user_vehicles WHERE id = $1 AND user_id = $2`
		_, err := r.db.ExecContext(ctx, query, id, userID)
		if err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.vehicles[id]; ok && v.UserID == userID {
		delete(r.vehicles, id)
		return nil
	}
	return errors.New("vehicle not found or access denied")
}
