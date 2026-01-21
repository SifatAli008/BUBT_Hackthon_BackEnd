package shopping_list

import (
	"time"

	"github.com/google/uuid"
)

// ShoppingListItem represents a shopping list item
type ShoppingListItem struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	HouseholdID   *uuid.UUID `json:"household_id,omitempty" db:"household_id"`
	Name          string     `json:"name" db:"name"`
	Quantity      float64    `json:"quantity" db:"quantity"`
	Unit          string     `json:"unit,omitempty" db:"unit"`
	Category      string     `json:"category,omitempty" db:"category"`
	Priority      string     `json:"priority,omitempty" db:"priority"` // low|medium|high
	Purchased     bool       `json:"purchased" db:"purchased"`
	PurchasedAt   *time.Time `json:"purchased_at,omitempty" db:"purchased_at"`
	EstimatedPrice *float64   `json:"estimated_price,omitempty" db:"estimated_price"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateShoppingListItemRequest struct {
	Name           string   `json:"name" validate:"required,min=1,max=255"`
	Quantity       float64  `json:"quantity" validate:"required,gt=0"`
	Unit           string   `json:"unit,omitempty" validate:"omitempty,max=50"`
	Category       string   `json:"category,omitempty" validate:"omitempty,max=100"`
	Priority       string   `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	EstimatedPrice *float64 `json:"estimated_price,omitempty"`
}

type UpdateShoppingListItemRequest struct {
	Name           string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Quantity       *float64  `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit           string    `json:"unit,omitempty" validate:"omitempty,max=50"`
	Category       string    `json:"category,omitempty" validate:"omitempty,max=100"`
	Priority       string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	Purchased      *bool     `json:"purchased,omitempty"`
	EstimatedPrice *float64  `json:"estimated_price,omitempty"`
	PurchasedAt    *time.Time `json:"purchased_at,omitempty"`
}

