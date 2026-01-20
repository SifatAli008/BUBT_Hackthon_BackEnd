package donations

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

type DonationLog struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Date          time.Time `json:"date" db:"date"`
	RecipientType string    `json:"recipient_type" db:"recipient_type"`
	RecipientName string    `json:"recipient_name" db:"recipient_name"`
	Items         string    `json:"items" db:"items"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	Unit          string    `json:"unit" db:"unit"`
	MealsProvided int       `json:"meals_provided" db:"meals_provided"`
	CO2SavedKg    float64   `json:"co2_saved_kg" db:"co2_saved_kg"`
	Notes         string    `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type ImpactMetrics struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	WastePreventedKg   float64   `json:"waste_prevented_kg" db:"waste_prevented_kg"`
	SurplusDonationRate float64  `json:"surplus_donation_rate" db:"surplus_donation_rate"`
	WaterSavedLiters   float64   `json:"water_saved_liters" db:"water_saved_liters"`
	CO2PreventedKg     float64   `json:"co2_prevented_kg" db:"co2_prevented_kg"`
	SustainabilityScore int      `json:"sustainability_score" db:"sustainability_score"`
	WeeklyTrend        JSONB     `json:"weekly_trend,omitempty" db:"weekly_trend"`
	MonthlyTrend       JSONB     `json:"monthly_trend,omitempty" db:"monthly_trend"`
	CategoryBreakdown  JSONB     `json:"category_breakdown,omitempty" db:"category_breakdown"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDonationLogRequest struct {
	Date          time.Time `json:"date" validate:"required"`
	RecipientType string    `json:"recipient_type" validate:"required,oneof=ngo community-kitchen"`
	RecipientName string    `json:"recipient_name" validate:"required,min=1"`
	Items         string    `json:"items" validate:"required,min=1"`
	Quantity      float64   `json:"quantity" validate:"required,gt=0"`
	Unit          string    `json:"unit" validate:"required,min=1"`
	MealsProvided int       `json:"meals_provided,omitempty"`
	CO2SavedKg    float64   `json:"co2_saved_kg,omitempty"`
	Notes         string    `json:"notes,omitempty"`
}
