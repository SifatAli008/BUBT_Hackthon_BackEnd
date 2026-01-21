package meal_plans

import (
	"time"

	"github.com/google/uuid"
)

// MealPlan represents a meal plan entry
type MealPlan struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	HouseholdID *uuid.UUID `json:"household_id,omitempty" db:"household_id"`
	Date        time.Time  `json:"date" db:"date"` // date-only
	MealType    string     `json:"meal_type" db:"meal_type"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	Ingredients []string   `json:"ingredients,omitempty" db:"ingredients"`
	Servings    *int       `json:"servings,omitempty" db:"servings"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type UpsertMealPlanRequest struct {
	Date        string   `json:"date" validate:"required"` // YYYY-MM-DD
	MealType    string   `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description string   `json:"description,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	Servings    *int     `json:"servings,omitempty"`
}

// WeeklyMealPlanResponse groups meals by date and type.
type WeeklyMealPlanResponse map[string]map[string]*MealPlan

