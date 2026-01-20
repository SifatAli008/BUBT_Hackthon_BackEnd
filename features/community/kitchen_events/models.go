package kitchen_events

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

type KitchenEvent struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	Date          time.Time `json:"date" db:"date"`
	Time          string    `json:"time" db:"time"`
	Location      string    `json:"location" db:"location"`
	Tags          []string  `json:"tags,omitempty" db:"tags"`
	VolunteersNeeded int    `json:"volunteers_needed" db:"volunteers_needed"`
	Volunteers    JSONB     `json:"volunteers" db:"volunteers"`
	FoodSavedKg   float64   `json:"food_saved_kg" db:"food_saved_kg"`
	Status        string    `json:"status" db:"status"`
	Image         string    `json:"image,omitempty" db:"image"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateKitchenEventRequest struct {
	Title          string    `json:"title" validate:"required,min=1,max=255"`
	Description   string    `json:"description" validate:"required,min=1"`
	Date          time.Time `json:"date" validate:"required"`
	Time          string    `json:"time" validate:"required"`
	Location      string    `json:"location" validate:"required,min=1"`
	Tags          []string  `json:"tags,omitempty"`
	VolunteersNeeded int    `json:"volunteers_needed,omitempty"`
	Image         string    `json:"image,omitempty"`
}

type UpdateKitchenEventRequest struct {
	Title          string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description   string    `json:"description,omitempty" validate:"omitempty,min=1"`
	Date          *time.Time `json:"date,omitempty"`
	Time          string    `json:"time,omitempty"`
	Location      string    `json:"location,omitempty" validate:"omitempty,min=1"`
	Tags          []string  `json:"tags,omitempty"`
	VolunteersNeeded *int   `json:"volunteers_needed,omitempty"`
	Status        string    `json:"status,omitempty"`
	FoodSavedKg   *float64  `json:"food_saved_kg,omitempty"`
	Image         string    `json:"image,omitempty"`
}

type VolunteerRequest struct {
	Role string `json:"role,omitempty"`
}
