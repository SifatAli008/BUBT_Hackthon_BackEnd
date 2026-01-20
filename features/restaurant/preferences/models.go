package preferences

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantPreferences struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	CuisineType        string    `json:"cuisine_type,omitempty" db:"cuisine_type"`
	OperatingHours     string    `json:"operating_hours,omitempty" db:"operating_hours"`
	DonationPreferences []string `json:"donation_preferences,omitempty" db:"donation_preferences"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRestaurantPreferencesRequest struct {
	CuisineType        string   `json:"cuisine_type,omitempty" validate:"omitempty,max=100"`
	OperatingHours     string   `json:"operating_hours,omitempty"`
	DonationPreferences []string `json:"donation_preferences,omitempty"`
}
