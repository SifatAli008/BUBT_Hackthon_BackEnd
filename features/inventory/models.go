package inventory

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents an inventory item
type InventoryItem struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	Quantity    float64    `json:"quantity" db:"quantity"`
	Unit        string     `json:"unit,omitempty" db:"unit"`
	ExpiryDate  *time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	Category    string     `json:"category,omitempty" db:"category"`
	Location    string     `json:"location,omitempty" db:"location"`
	FoodItemID  *uuid.UUID `json:"food_item_id,omitempty" db:"food_item_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateInventoryItemRequest represents a request to create an inventory item
type CreateInventoryItemRequest struct {
	Name       string     `json:"name" validate:"required,min=1,max=255"`
	Quantity   float64    `json:"quantity" validate:"required,gt=0"`
	Unit       string     `json:"unit,omitempty" validate:"omitempty,max=50"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	Category   string     `json:"category,omitempty" validate:"omitempty,max=100"`
	Location   string     `json:"location,omitempty" validate:"omitempty,max=100"`
	FoodItemID *uuid.UUID `json:"food_item_id,omitempty"`
}

// UpdateInventoryItemRequest represents a request to update an inventory item
type UpdateInventoryItemRequest struct {
	Name       string     `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Quantity   *float64   `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit       string     `json:"unit,omitempty" validate:"omitempty,max=50"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	Category   string     `json:"category,omitempty" validate:"omitempty,max=100"`
	Location   string     `json:"location,omitempty" validate:"omitempty,max=100"`
	FoodItemID *uuid.UUID `json:"food_item_id,omitempty"`
}
