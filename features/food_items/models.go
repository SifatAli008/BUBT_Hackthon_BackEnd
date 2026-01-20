package food_items

import (
	"time"

	"github.com/google/uuid"
)

// FoodItem represents a food item reference
type FoodItem struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Category         string    `json:"category" db:"category"`
	TypicalExpiryDays int      `json:"typical_expiry_days" db:"typical_expiry_days"`
	StorageTips      string    `json:"storage_tips,omitempty" db:"storage_tips"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// CreateFoodItemRequest represents a request to create a food item
type CreateFoodItemRequest struct {
	Name             string `json:"name" validate:"required,min=1,max=255"`
	Category         string `json:"category" validate:"required,min=1,max=100"`
	TypicalExpiryDays int   `json:"typical_expiry_days" validate:"required,min=1"`
	StorageTips      string `json:"storage_tips,omitempty"`
}

// UpdateFoodItemRequest represents a request to update a food item
type UpdateFoodItemRequest struct {
	Name             string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Category         string `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	TypicalExpiryDays *int  `json:"typical_expiry_days,omitempty" validate:"omitempty,min=1"`
	StorageTips      string `json:"storage_tips,omitempty"`
}
