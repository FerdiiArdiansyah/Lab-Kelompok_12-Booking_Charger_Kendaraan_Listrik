package usecase

import (
	"context"
	"errors"
	"testing"
	"user-service/internal/domain"
)

type MockUserRepository struct {
	users    map[string]*domain.User
	vehicles map[string][]domain.UserVehicle
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:    make(map[string]*domain.User),
		vehicles: make(map[string][]domain.UserVehicle),
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var list []domain.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; ok {
		m.users[user.ID] = user
		return nil
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) CreateVehicle(ctx context.Context, vehicle *domain.UserVehicle) error {
	m.vehicles[vehicle.UserID] = append(m.vehicles[vehicle.UserID], *vehicle)
	return nil
}

func (m *MockUserRepository) GetVehiclesByUserID(ctx context.Context, userID string) ([]domain.UserVehicle, error) {
	return m.vehicles[userID], nil
}

func (m *MockUserRepository) DeleteVehicle(ctx context.Context, id, userID string) error {
	var remaining []domain.UserVehicle
	for _, v := range m.vehicles[userID] {
		if v.ID != id {
			remaining = append(remaining, v)
		}
	}
	m.vehicles[userID] = remaining
	return nil
}

// Unit Tests for TICKET-07 (User & EV Connector Validation)

func TestRegisterAndLogin_Success(t *testing.T) {
	repo := NewMockUserRepository()
	uc := NewUserUsecase(repo, "secret-jwt-key")
	ctx := context.Background()

	regReq := &domain.RegisterRequest{
		Name:     "Driver Test",
		Email:    "driver@spklu.id",
		Password: "password123",
	}

	res, err := uc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("Expected registration success, got: %v", err)
	}

	if res.Token == "" {
		t.Errorf("Expected valid JWT token in auth response")
	}

	loginRes, err := uc.Login(ctx, &domain.LoginRequest{
		Email:    "driver@spklu.id",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Expected login success, got: %v", err)
	}

	if loginRes.User.Email != "driver@spklu.id" {
		t.Errorf("Expected logged in email driver@spklu.id, got %s", loginRes.User.Email)
	}
}

func TestAddVehicle_ValidConnectorType_Success(t *testing.T) {
	repo := NewMockUserRepository()
	uc := NewUserUsecase(repo, "secret-jwt-key")
	ctx := context.Background()

	req := &domain.AddVehicleRequest{
		Brand:              "Hyundai",
		Model:              "Ioniq 5",
		LicensePlate:       "B 1234 EV",
		ConnectorType:      "CCS2",
		BatteryCapacityKwh: 72.6,
	}

	vehicle, err := uc.AddVehicle(ctx, "usr-1", req)
	if err != nil {
		t.Fatalf("Expected no error adding valid vehicle, got: %v", err)
	}

	if vehicle.ConnectorType != "CCS2" {
		t.Errorf("Expected ConnectorType CCS2, got %s", vehicle.ConnectorType)
	}
}

func TestAddVehicle_InvalidConnectorType_ReturnsError(t *testing.T) {
	repo := NewMockUserRepository()
	uc := NewUserUsecase(repo, "secret-jwt-key")
	ctx := context.Background()

	req := &domain.AddVehicleRequest{
		Brand:              "Custom EV",
		Model:              "Prototype",
		LicensePlate:       "B 9999 EV",
		ConnectorType:      "UNKNOWN_CONNECTOR",
		BatteryCapacityKwh: 50.0,
	}

	_, err := uc.AddVehicle(ctx, "usr-1", req)
	if err == nil {
		t.Errorf("Expected error when adding vehicle with invalid connector type, got nil")
	}
}
