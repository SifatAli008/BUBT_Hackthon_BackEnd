package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	Name        string    `json:"name" db:"name"`
	PasswordHash string   `json:"-" db:"password_hash"`
	HouseholdID *uuid.UUID `json:"household_id,omitempty" db:"household_id"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=family restaurant shop ngo admin"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserResponse represents a user response (without sensitive data)
type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	HouseholdID *uuid.UUID `json:"household_id,omitempty"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		HouseholdID: u.HouseholdID,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
