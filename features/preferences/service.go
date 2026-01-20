package preferences

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) GetByHouseholdID(householdID uuid.UUID) (*FamilyPreferences, error) {
	return s.repo.GetByHouseholdID(householdID)
}

func (s *Service) CreateOrUpdate(req *CreatePreferencesRequest) (*FamilyPreferences, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}

	prefs := &FamilyPreferences{
		ID:                      uuid.New(),
		HouseholdID:             req.HouseholdID,
		HouseholdSize:           req.HouseholdSize,
		AgeGroups:               JSONB(req.AgeGroups),
		CookingFrequency:        req.CookingFrequency,
		EatingSchedule:          JSONB(req.EatingSchedule),
		DietaryType:             req.DietaryType,
		DietaryRestrictions:     req.DietaryRestrictions,
		Allergies:               req.Allergies,
		HealthConditions:        req.HealthConditions,
		WeeklyBudget:            req.WeeklyBudget,
		BudgetRange:             JSONB(req.BudgetRange),
		PreferredStores:         req.PreferredStores,
		PriceSensitivity:        req.PriceSensitivity,
		PreferredCuisines:       req.PreferredCuisines,
		MealPrepPreference:      req.MealPrepPreference,
		WasteSensitivityLevel:   req.WasteSensitivityLevel,
		SustainabilityPreference: req.SustainabilityPreference,
		LeftoverComfortLevel:    req.LeftoverComfortLevel,
		DailyCalories:           req.DailyCalories,
		MacroGoal:               JSONB(req.MacroGoal),
		VitaminsFocus:           req.VitaminsFocus,
		AvoidExcess:             req.AvoidExcess,
	}

	if err := s.repo.Create(prefs); err != nil {
		return nil, err
	}

	return prefs, nil
}
