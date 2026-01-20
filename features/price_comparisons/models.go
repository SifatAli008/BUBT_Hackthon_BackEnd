package price_comparisons

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONB is a helper type for PostgreSQL JSONB fields
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

// StorePrice represents a store price entry
type StorePrice struct {
	StoreName string  `json:"store_name"`
	Price     float64 `json:"price"`
	Unit      string  `json:"unit,omitempty"`
	Available bool    `json:"available"`
}

// PriceComparison represents a price comparison
type PriceComparison struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ItemName   string    `json:"item_name" db:"item_name"`
	Category   string    `json:"category,omitempty" db:"category"`
	Stores     JSONB     `json:"stores" db:"stores"`
	BestPrice  JSONB     `json:"best_price" db:"best_price"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePriceComparisonRequest represents a request to create a price comparison
type CreatePriceComparisonRequest struct {
	ItemName  string                   `json:"item_name" validate:"required,min=1,max=255"`
	Category  string                   `json:"category,omitempty" validate:"omitempty,max=100"`
	Stores    []map[string]interface{} `json:"stores" validate:"required"`
	BestPrice map[string]interface{}   `json:"best_price" validate:"required"`
}

// UpdatePriceComparisonRequest represents a request to update a price comparison
type UpdatePriceComparisonRequest struct {
	ItemName  string                   `json:"item_name,omitempty" validate:"omitempty,min=1,max=255"`
	Category  string                   `json:"category,omitempty" validate:"omitempty,max=100"`
	Stores    []map[string]interface{} `json:"stores,omitempty"`
	BestPrice map[string]interface{}   `json:"best_price,omitempty"`
}
