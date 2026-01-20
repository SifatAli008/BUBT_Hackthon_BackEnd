package auth

import (
	"foodlink_backend/config"
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication business logic
type Service struct {
	repo     *Repository
	cfg      *config.Config
	jwtExpiry time.Duration
}

// NewService creates a new auth service
func NewService(cfg *config.Config) *Service {
	expiry, _ := utils.ParseExpiry(cfg.JWTExpiry)
	if expiry == 0 {
		expiry = 24 * time.Hour // Default 24 hours
	}

	return &Service{
		repo:      NewRepository(),
		cfg:       cfg,
		jwtExpiry: expiry,
	}
}

// Register registers a new user
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(
			errors.ErrValidationFailed.Code,
			"Validation failed: "+validationErrors[0],
			nil,
		)
	}

	// Check if email already exists
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrInternalServer)
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "family"
	}

	// Create user
	user := &User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtExpiry)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrInternalServer)
	}

	return &AuthResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.jwtExpiry.Seconds()),
	}, nil
}

// Login authenticates a user and returns a token
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	// Validate input
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(
			errors.ErrValidationFailed.Code,
			"Validation failed: "+validationErrors[0],
			nil,
		)
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtExpiry)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrInternalServer)
	}

	return &AuthResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.jwtExpiry.Seconds()),
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *Service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetUserByEmail(email)
}

// ValidateToken validates a JWT token and returns user info
func (s *Service) ValidateToken(tokenString string) (*User, error) {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	return user, nil
}
