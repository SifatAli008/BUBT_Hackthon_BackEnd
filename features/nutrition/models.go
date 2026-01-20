package nutrition

import (
	"time"

	"github.com/google/uuid"
)

// NutritionData represents daily nutrition data
type NutritionData struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Date           time.Time `json:"date" db:"date"`
	Calories       float64   `json:"calories" db:"calories"`
	Protein        float64   `json:"protein" db:"protein"`
	Carbs          float64   `json:"carbs" db:"carbs"`
	Fats           float64   `json:"fats" db:"fats"`
	Fiber          float64   `json:"fiber" db:"fiber"`
	Sugar          float64   `json:"sugar" db:"sugar"`
	Sodium         float64   `json:"sodium" db:"sodium"`
	VitaminA       float64   `json:"vitamin_a" db:"vitamin_a"`
	VitaminB       float64   `json:"vitamin_b" db:"vitamin_b"`
	VitaminC       float64   `json:"vitamin_c" db:"vitamin_c"`
	VitaminD       float64   `json:"vitamin_d" db:"vitamin_d"`
	Iron           float64   `json:"iron" db:"iron"`
	Calcium        float64   `json:"calcium" db:"calcium"`
	NutritionScore *int      `json:"nutrition_score,omitempty" db:"nutrition_score"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreateNutritionDataRequest represents a request to create nutrition data
type CreateNutritionDataRequest struct {
	Date           time.Time `json:"date"`
	Calories       float64   `json:"calories"`
	Protein        float64   `json:"protein"`
	Carbs          float64   `json:"carbs"`
	Fats           float64   `json:"fats"`
	Fiber          float64   `json:"fiber"`
	Sugar          float64   `json:"sugar"`
	Sodium         float64   `json:"sodium"`
	VitaminA       float64   `json:"vitamin_a"`
	VitaminB       float64   `json:"vitamin_b"`
	VitaminC       float64   `json:"vitamin_c"`
	VitaminD       float64   `json:"vitamin_d"`
	Iron           float64   `json:"iron"`
	Calcium        float64   `json:"calcium"`
	NutritionScore *int      `json:"nutrition_score,omitempty"`
}

// UpdateNutritionDataRequest represents a request to update nutrition data
type UpdateNutritionDataRequest struct {
	Calories       *float64 `json:"calories,omitempty"`
	Protein        *float64 `json:"protein,omitempty"`
	Carbs          *float64 `json:"carbs,omitempty"`
	Fats           *float64 `json:"fats,omitempty"`
	Fiber          *float64 `json:"fiber,omitempty"`
	Sugar          *float64 `json:"sugar,omitempty"`
	Sodium         *float64 `json:"sodium,omitempty"`
	VitaminA       *float64 `json:"vitamin_a,omitempty"`
	VitaminB       *float64 `json:"vitamin_b,omitempty"`
	VitaminC       *float64 `json:"vitamin_c,omitempty"`
	VitaminD       *float64 `json:"vitamin_d,omitempty"`
	Iron           *float64 `json:"iron,omitempty"`
	Calcium        *float64 `json:"calcium,omitempty"`
	NutritionScore *int     `json:"nutrition_score,omitempty"`
}

// NutritionStats represents nutrition statistics
type NutritionStats struct {
	AvgCalories       float64 `json:"avg_calories"`
	AvgProtein        float64 `json:"avg_protein"`
	AvgCarbs          float64 `json:"avg_carbs"`
	AvgFats           float64 `json:"avg_fats"`
	TotalDays         int     `json:"total_days"`
	AvgNutritionScore *int    `json:"avg_nutrition_score,omitempty"`
}
