package leaderboard

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

type Leaderboard struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	Entries   JSONB     `json:"entries" db:"entries"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Impact struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	TotalSurplusKg     float64   `json:"total_surplus_kg" db:"total_surplus_kg"`
	Donations           int       `json:"donations" db:"donations"`
	CO2PreventedKg      float64   `json:"co2_prevented_kg" db:"co2_prevented_kg"`
	WaterSavedLiters   float64   `json:"water_saved_liters" db:"water_saved_liters"`
	MealsProvided       int       `json:"meals_provided" db:"meals_provided"`
	WeeklyTrend         JSONB     `json:"weekly_trend,omitempty" db:"weekly_trend"`
	PersonalContribution JSONB    `json:"personal_contribution,omitempty" db:"personal_contribution"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
