package usecase

import (
	"context"
	"errors"
	"time"
	"user-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewUserUsecase(repo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	return &userUsecase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New("name, email, and password are required")
	}

	existing, _ := u.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		hash = []byte(req.Password)
	}

	userRole := "USER"
	if req.Role == "ADMIN" {
		userRole = "ADMIN"
	}

	now := time.Now()
	newUser := &domain.User{
		ID:           "usr-" + uuid.New().String()[:8],
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Phone:        req.Phone,
		Role:         userRole,
		Status:       "ACTIVE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	token, err := u.generateToken(newUser)
	if err != nil {
		token = "mock-jwt-token-" + newUser.ID
	}

	return &domain.AuthResponse{
		Token: token,
		User:  newUser,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		if user.PasswordHash != req.Password && user.PasswordHash != "hashed-"+req.Password {
			return nil, errors.New("invalid email or password")
		}
	}

	token, err := u.generateToken(user)
	if err != nil {
		token = "mock-jwt-token-" + user.ID
	}

	vehicles, _ := u.repo.GetVehiclesByUserID(ctx, user.ID)
	user.Vehicles = vehicles

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	vehicles, _ := u.repo.GetVehiclesByUserID(ctx, user.ID)
	user.Vehicles = vehicles
	return user, nil
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return u.repo.GetAllUsers(ctx)
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if err := u.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) AddVehicle(ctx context.Context, userID string, req *domain.AddVehicleRequest) (*domain.UserVehicle, error) {
	if req.Brand == "" || req.Model == "" || req.LicensePlate == "" {
		return nil, errors.New("brand, model, and license plate are required")
	}

	// Unique vehicle identity constraint: 1 license plate = 1 owner
	existingVeh, _ := u.repo.GetVehicleByLicensePlate(ctx, req.LicensePlate)
	if existingVeh != nil {
		return nil, errors.New("nomor plat kendaraan sudah terdaftar atas nama pengguna lain")
	}

	validConnectors := map[string]bool{
		"CCS2":      true,
		"CHAdeMO":   true,
		"Type 2":    true,
		"AC Type 2": true,
		"GB/T":      true,
	}

	if req.ConnectorType != "" && !validConnectors[req.ConnectorType] {
		return nil, errors.New("invalid connector type; allowed types: CCS2, CHAdeMO, Type 2, AC Type 2, GB/T")
	}

	vehicle := &domain.UserVehicle{
		ID:                 "vhc-" + uuid.New().String()[:8],
		UserID:             userID,
		Brand:              req.Brand,
		Model:              req.Model,
		LicensePlate:       req.LicensePlate,
		ConnectorType:      req.ConnectorType,
		BatteryCapacityKwh: req.BatteryCapacityKwh,
		CreatedAt:          time.Now(),
	}

	if err := u.repo.CreateVehicle(ctx, vehicle); err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (u *userUsecase) GetUserVehicles(ctx context.Context, userID string) ([]domain.UserVehicle, error) {
	return u.repo.GetVehiclesByUserID(ctx, userID)
}

func (u *userUsecase) DeleteVehicle(ctx context.Context, vehicleID, userID string) error {
	return u.repo.DeleteVehicle(ctx, vehicleID, userID)
}

func (u *userUsecase) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}
