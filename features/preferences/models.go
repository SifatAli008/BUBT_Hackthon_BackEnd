package preferences

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

// FamilyPreferences represents family preferences
type FamilyPreferences struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	HouseholdID             uuid.UUID `json:"household_id" db:"household_id"`
	HouseholdSize           int       `json:"household_size" db:"household_size"`
	AgeGroups               JSONB      `json:"age_groups,omitempty" db:"age_groups"`
	CookingFrequency        string    `json:"cooking_frequency,omitempty" db:"cooking_frequency"`
	EatingSchedule          JSONB      `json:"eating_schedule,omitempty" db:"eating_schedule"`
	DietaryType             string    `json:"dietary_type,omitempty" db:"dietary_type"`
	DietaryRestrictions     []string  `json:"dietary_restrictions,omitempty" db:"dietary_restrictions"`
	Allergies               []string  `json:"allergies,omitempty" db:"allergies"`
	HealthConditions        []string  `json:"health_conditions,omitempty" db:"health_conditions"`
	WeeklyBudget            *float64  `json:"weekly_budget,omitempty" db:"weekly_budget"`
	BudgetRange             JSONB      `json:"budget_range,omitempty" db:"budget_range"`
	PreferredStores         []string  `json:"preferred_stores,omitempty" db:"preferred_stores"`
	PriceSensitivity        string    `json:"price_sensitivity,omitempty" db:"price_sensitivity"`
	PreferredCuisines       []string  `json:"preferred_cuisines,omitempty" db:"preferred_cuisines"`
	MealPrepPreference      string    `json:"meal_prep_preference,omitempty" db:"meal_prep_preference"`
	WasteSensitivityLevel   string    `json:"waste_sensitivity_level,omitempty" db:"waste_sensitivity_level"`
	SustainabilityPreference string    `json:"sustainability_preference,omitempty" db:"sustainability_preference"`
	LeftoverComfortLevel    string    `json:"leftover_comfort_level,omitempty" db:"leftover_comfort_level"`
	DailyCalories           *int      `json:"daily_calories,omitempty" db:"daily_calories"`
	MacroGoal               JSONB      `json:"macro_goal,omitempty" db:"macro_goal"`
	VitaminsFocus           []string  `json:"vitamins_focus,omitempty" db:"vitamins_focus"`
	AvoidExcess             []string  `json:"avoid_excess,omitempty" db:"avoid_excess"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePreferencesRequest represents a request to create/update preferences
type CreatePreferencesRequest struct {
	HouseholdID             uuid.UUID              `json:"household_id" validate:"required"`
	HouseholdSize           int                    `json:"household_size" validate:"required,min=1"`
	AgeGroups               map[string]interface{} `json:"age_groups,omitempty"`
	CookingFrequency        string                 `json:"cooking_frequency,omitempty"`
	EatingSchedule          map[string]interface{} `json:"eating_schedule,omitempty"`
	DietaryType             string                 `json:"dietary_type,omitempty"`
	DietaryRestrictions     []string               `json:"dietary_restrictions,omitempty"`
	Allergies               []string               `json:"allergies,omitempty"`
	HealthConditions        []string               `json:"health_conditions,omitempty"`
	WeeklyBudget            *float64               `json:"weekly_budget,omitempty"`
	BudgetRange             map[string]interface{} `json:"budget_range,omitempty"`
	PreferredStores         []string               `json:"preferred_stores,omitempty"`
	PriceSensitivity        string                 `json:"price_sensitivity,omitempty"`
	PreferredCuisines       []string               `json:"preferred_cuisines,omitempty"`
	MealPrepPreference      string                 `json:"meal_prep_preference,omitempty"`
	WasteSensitivityLevel   string                 `json:"waste_sensitivity_level,omitempty"`
	SustainabilityPreference string                 `json:"sustainability_preference,omitempty"`
	LeftoverComfortLevel    string                 `json:"leftover_comfort_level,omitempty"`
	DailyCalories           *int                   `json:"daily_calories,omitempty"`
	MacroGoal               map[string]interface{} `json:"macro_goal,omitempty"`
	VitaminsFocus           []string               `json:"vitamins_focus,omitempty"`
	AvoidExcess             []string               `json:"avoid_excess,omitempty"`
}
