package leftovers

import (
	"time"

	"github.com/google/uuid"
)

type LeftoverItem struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name" db:"user_name"`
	AvatarURL    string    `json:"avatar_url,omitempty" db:"avatar_url"`
	DishName     string    `json:"dish_name" db:"dish_name"`
	Description  string    `json:"description" db:"description"`
	Portions     int       `json:"portions" db:"portions"`
	DistanceKm   float64   `json:"distance_km" db:"distance_km"`
	DietaryTags  []string  `json:"dietary_tags,omitempty" db:"dietary_tags"`
	Allergens    []string  `json:"allergens,omitempty" db:"allergens"`
	PickupWindow string    `json:"pickup_window" db:"pickup_window"`
	Status       string    `json:"status" db:"status"`
	Image        string    `json:"image,omitempty" db:"image"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type LeftoverClaim struct {
	ID            uuid.UUID `json:"id" db:"id"`
	LeftoverItemID uuid.UUID `json:"leftover_item_id" db:"leftover_item_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	UserName      string    `json:"user_name" db:"user_name"`
	Message       string    `json:"message,omitempty" db:"message"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type CreateLeftoverItemRequest struct {
	DishName     string   `json:"dish_name" validate:"required,min=1,max=255"`
	Description  string   `json:"description" validate:"required,min=1"`
	Portions     int      `json:"portions" validate:"required,gt=0"`
	DistanceKm   float64  `json:"distance_km" validate:"required,gt=0"`
	DietaryTags  []string `json:"dietary_tags,omitempty"`
	Allergens    []string `json:"allergens,omitempty"`
	PickupWindow string   `json:"pickup_window" validate:"required,min=1"`
	Image        string   `json:"image,omitempty"`
}

type UpdateLeftoverItemRequest struct {
	DishName     string   `json:"dish_name,omitempty" validate:"omitempty,min=1,max=255"`
	Description  string   `json:"description,omitempty" validate:"omitempty,min=1"`
	Portions     *int     `json:"portions,omitempty" validate:"omitempty,gt=0"`
	DistanceKm   *float64 `json:"distance_km,omitempty" validate:"omitempty,gt=0"`
	DietaryTags  []string `json:"dietary_tags,omitempty"`
	Allergens    []string `json:"allergens,omitempty"`
	PickupWindow string   `json:"pickup_window,omitempty" validate:"omitempty,min=1"`
	Status       string   `json:"status,omitempty"`
	Image        string   `json:"image,omitempty"`
}

type CreateLeftoverClaimRequest struct {
	Message string `json:"message,omitempty"`
}
