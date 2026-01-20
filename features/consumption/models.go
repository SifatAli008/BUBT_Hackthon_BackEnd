package consumption

import (
	"time"

	"github.com/google/uuid"
)

// ConsumptionLog represents a consumption log entry
type ConsumptionLog struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	InventoryItemID *uuid.UUID `json:"inventory_item_id,omitempty" db:"inventory_item_id"`
	FoodName        string     `json:"food_name" db:"food_name"`
	Quantity        float64    `json:"quantity" db:"quantity"`
	Unit            string     `json:"unit,omitempty" db:"unit"`
	Category        string     `json:"category,omitempty" db:"category"`
	ConsumedAt      time.Time  `json:"consumed_at" db:"consumed_at"`
	WasWasted       bool       `json:"was_wasted" db:"was_wasted"`
	Notes           string     `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateConsumptionLogRequest represents a request to create a consumption log
type CreateConsumptionLogRequest struct {
	InventoryItemID *uuid.UUID `json:"inventory_item_id,omitempty"`
	FoodName        string     `json:"food_name" validate:"required,min=1,max=255"`
	Quantity        float64    `json:"quantity" validate:"required,gt=0"`
	Unit            string     `json:"unit,omitempty" validate:"omitempty,max=50"`
	Category        string     `json:"category,omitempty" validate:"omitempty,max=100"`
	ConsumedAt      time.Time  `json:"consumed_at"`
	WasWasted       bool       `json:"was_wasted"`
	Notes           string     `json:"notes,omitempty"`
}

// UpdateConsumptionLogRequest represents a request to update a consumption log
type UpdateConsumptionLogRequest struct {
	FoodName   string     `json:"food_name,omitempty" validate:"omitempty,min=1,max=255"`
	Quantity   *float64   `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit       string     `json:"unit,omitempty" validate:"omitempty,max=50"`
	Category   string     `json:"category,omitempty" validate:"omitempty,max=100"`
	ConsumedAt *time.Time `json:"consumed_at,omitempty"`
	WasWasted  *bool      `json:"was_wasted,omitempty"`
	Notes      string     `json:"notes,omitempty"`
}

// ConsumptionStats represents consumption statistics
type ConsumptionStats struct {
	TotalConsumed   float64 `json:"total_consumed"`
	TotalWasted     float64 `json:"total_wasted"`
	WastePercentage float64 `json:"waste_percentage"`
	TotalLogs       int     `json:"total_logs"`
}
