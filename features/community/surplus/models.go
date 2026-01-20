package surplus

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type SurplusPost struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name" db:"user_name"`
	AvatarURL    string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Title        string    `json:"title" db:"title"`
	Description  string   `json:"description" db:"description"`
	Category     string    `json:"category" db:"category"`
	Tags         []string  `json:"tags,omitempty" db:"tags"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	Unit         string    `json:"unit" db:"unit"`
	PickupWindow JSONB     `json:"pickup_window" db:"pickup_window"`
	PickupLocation string  `json:"pickup_location" db:"pickup_location"`
	DistanceKm   *float64  `json:"distance_km,omitempty" db:"distance_km"`
	Image        string    `json:"image,omitempty" db:"image"`
	Status       string    `json:"status" db:"status"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type SurplusRequest struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name" db:"user_name"`
	Message   string    `json:"message,omitempty" db:"message"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type SurplusComment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name" db:"user_name"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateSurplusPostRequest struct {
	Title         string                 `json:"title" validate:"required,min=1,max=255"`
	Description   string                 `json:"description" validate:"required,min=1"`
	Category      string                 `json:"category" validate:"required,min=1,max=100"`
	Tags          []string               `json:"tags,omitempty"`
	Quantity      float64                `json:"quantity" validate:"required,gt=0"`
	Unit          string                 `json:"unit" validate:"required,min=1,max=50"`
	PickupWindow  map[string]interface{} `json:"pickup_window" validate:"required"`
	PickupLocation string                `json:"pickup_location" validate:"required,min=1"`
	DistanceKm    *float64               `json:"distance_km,omitempty"`
	Image         string                 `json:"image,omitempty"`
	ExpiresAt     time.Time              `json:"expires_at" validate:"required"`
}

type UpdateSurplusPostRequest struct {
	Title         string                 `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description   string                 `json:"description,omitempty" validate:"omitempty,min=1"`
	Category      string                 `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	Tags          []string               `json:"tags,omitempty"`
	Quantity      *float64               `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit          string                 `json:"unit,omitempty" validate:"omitempty,min=1,max=50"`
	PickupWindow  map[string]interface{} `json:"pickup_window,omitempty"`
	PickupLocation string                `json:"pickup_location,omitempty" validate:"omitempty,min=1"`
	Status        string                 `json:"status,omitempty"`
	Image         string                 `json:"image,omitempty"`
}

type CreateSurplusRequestRequest struct {
	Message string `json:"message,omitempty"`
}

type UpdateSurplusRequestRequest struct {
	Status string `json:"status" validate:"required,oneof=pending approved declined"`
}

type CreateSurplusCommentRequest struct {
	Message string `json:"message" validate:"required,min=1"`
}
