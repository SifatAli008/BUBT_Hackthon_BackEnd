package preferences

import (
	"database/sql"
	"encoding/json"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetByHouseholdID(householdID uuid.UUID) (*FamilyPreferences, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	prefs := &FamilyPreferences{}
	var ageGroupsJSON, eatingScheduleJSON, budgetRangeJSON, macroGoalJSON []byte

	query := `
		SELECT id, household_id, household_size, age_groups, cooking_frequency, eating_schedule,
			dietary_type, dietary_restrictions, allergies, health_conditions, weekly_budget,
			budget_range, preferred_stores, price_sensitivity, preferred_cuisines, meal_prep_preference,
			waste_sensitivity_level, sustainability_preference, leftover_comfort_level, daily_calories,
			macro_goal, vitamins_focus, avoid_excess, created_at, updated_at
		FROM family_preferences
		WHERE household_id = $1
	`

	err := r.db.QueryRow(query, householdID).Scan(
		&prefs.ID, &prefs.HouseholdID, &prefs.HouseholdSize,
		&ageGroupsJSON, &prefs.CookingFrequency, &eatingScheduleJSON,
		&prefs.DietaryType, pq.Array(&prefs.DietaryRestrictions), pq.Array(&prefs.Allergies),
		pq.Array(&prefs.HealthConditions), &prefs.WeeklyBudget, &budgetRangeJSON,
		pq.Array(&prefs.PreferredStores), &prefs.PriceSensitivity, pq.Array(&prefs.PreferredCuisines),
		&prefs.MealPrepPreference, &prefs.WasteSensitivityLevel, &prefs.SustainabilityPreference,
		&prefs.LeftoverComfortLevel, &prefs.DailyCalories, &macroGoalJSON,
		pq.Array(&prefs.VitaminsFocus), pq.Array(&prefs.AvoidExcess),
		&prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}

	if len(ageGroupsJSON) > 0 {
		json.Unmarshal(ageGroupsJSON, &prefs.AgeGroups)
	}
	if len(eatingScheduleJSON) > 0 {
		json.Unmarshal(eatingScheduleJSON, &prefs.EatingSchedule)
	}
	if len(budgetRangeJSON) > 0 {
		json.Unmarshal(budgetRangeJSON, &prefs.BudgetRange)
	}
	if len(macroGoalJSON) > 0 {
		json.Unmarshal(macroGoalJSON, &prefs.MacroGoal)
	}

	return prefs, nil
}

func (r *Repository) Create(prefs *FamilyPreferences) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	ageGroupsJSON, _ := json.Marshal(prefs.AgeGroups)
	eatingScheduleJSON, _ := json.Marshal(prefs.EatingSchedule)
	budgetRangeJSON, _ := json.Marshal(prefs.BudgetRange)
	macroGoalJSON, _ := json.Marshal(prefs.MacroGoal)

	query := `
		INSERT INTO family_preferences (
			id, household_id, household_size, age_groups, cooking_frequency, eating_schedule,
			dietary_type, dietary_restrictions, allergies, health_conditions, weekly_budget,
			budget_range, preferred_stores, price_sensitivity, preferred_cuisines, meal_prep_preference,
			waste_sensitivity_level, sustainability_preference, leftover_comfort_level, daily_calories,
			macro_goal, vitamins_focus, avoid_excess, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
		)
		ON CONFLICT (household_id) DO UPDATE SET
			household_size = EXCLUDED.household_size,
			age_groups = EXCLUDED.age_groups,
			cooking_frequency = EXCLUDED.cooking_frequency,
			eating_schedule = EXCLUDED.eating_schedule,
			dietary_type = EXCLUDED.dietary_type,
			dietary_restrictions = EXCLUDED.dietary_restrictions,
			allergies = EXCLUDED.allergies,
			health_conditions = EXCLUDED.health_conditions,
			weekly_budget = EXCLUDED.weekly_budget,
			budget_range = EXCLUDED.budget_range,
			preferred_stores = EXCLUDED.preferred_stores,
			price_sensitivity = EXCLUDED.price_sensitivity,
			preferred_cuisines = EXCLUDED.preferred_cuisines,
			meal_prep_preference = EXCLUDED.meal_prep_preference,
			waste_sensitivity_level = EXCLUDED.waste_sensitivity_level,
			sustainability_preference = EXCLUDED.sustainability_preference,
			leftover_comfort_level = EXCLUDED.leftover_comfort_level,
			daily_calories = EXCLUDED.daily_calories,
			macro_goal = EXCLUDED.macro_goal,
			vitamins_focus = EXCLUDED.vitamins_focus,
			avoid_excess = EXCLUDED.avoid_excess,
			updated_at = EXCLUDED.updated_at
		RETURNING id, household_id, household_size, age_groups, cooking_frequency, eating_schedule,
			dietary_type, dietary_restrictions, allergies, health_conditions, weekly_budget,
			budget_range, preferred_stores, price_sensitivity, preferred_cuisines, meal_prep_preference,
			waste_sensitivity_level, sustainability_preference, leftover_comfort_level, daily_calories,
			macro_goal, vitamins_focus, avoid_excess, created_at, updated_at
	`

	now := time.Now()
	var ageGroupsJSONOut, eatingScheduleJSONOut, budgetRangeJSONOut, macroGoalJSONOut []byte

	err := r.db.QueryRow(query,
		prefs.ID, prefs.HouseholdID, prefs.HouseholdSize,
		ageGroupsJSON, prefs.CookingFrequency, eatingScheduleJSON,
		prefs.DietaryType, pq.Array(prefs.DietaryRestrictions), pq.Array(prefs.Allergies),
		pq.Array(prefs.HealthConditions), prefs.WeeklyBudget, budgetRangeJSON,
		pq.Array(prefs.PreferredStores), prefs.PriceSensitivity, pq.Array(prefs.PreferredCuisines),
		prefs.MealPrepPreference, prefs.WasteSensitivityLevel, prefs.SustainabilityPreference,
		prefs.LeftoverComfortLevel, prefs.DailyCalories, macroGoalJSON,
		pq.Array(prefs.VitaminsFocus), pq.Array(prefs.AvoidExcess),
		now, now,
	).Scan(
		&prefs.ID, &prefs.HouseholdID, &prefs.HouseholdSize,
		&ageGroupsJSONOut, &prefs.CookingFrequency, &eatingScheduleJSONOut,
		&prefs.DietaryType, pq.Array(&prefs.DietaryRestrictions), pq.Array(&prefs.Allergies),
		pq.Array(&prefs.HealthConditions), &prefs.WeeklyBudget, &budgetRangeJSONOut,
		pq.Array(&prefs.PreferredStores), &prefs.PriceSensitivity, pq.Array(&prefs.PreferredCuisines),
		&prefs.MealPrepPreference, &prefs.WasteSensitivityLevel, &prefs.SustainabilityPreference,
		&prefs.LeftoverComfortLevel, &prefs.DailyCalories, &macroGoalJSONOut,
		pq.Array(&prefs.VitaminsFocus), pq.Array(&prefs.AvoidExcess),
		&prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}

	if len(ageGroupsJSONOut) > 0 {
		json.Unmarshal(ageGroupsJSONOut, &prefs.AgeGroups)
	}
	if len(eatingScheduleJSONOut) > 0 {
		json.Unmarshal(eatingScheduleJSONOut, &prefs.EatingSchedule)
	}
	if len(budgetRangeJSONOut) > 0 {
		json.Unmarshal(budgetRangeJSONOut, &prefs.BudgetRange)
	}
	if len(macroGoalJSONOut) > 0 {
		json.Unmarshal(macroGoalJSONOut, &prefs.MacroGoal)
	}

	return nil
}
