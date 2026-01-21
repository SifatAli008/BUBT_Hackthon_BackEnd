package meal_plans

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func parseDateOnly(s string) (time.Time, error) {
	// Accept YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (s *Service) GetWeeklyByUserID(userID uuid.UUID, startDate string) (WeeklyMealPlanResponse, error) {
	var start time.Time
	var err error
	if startDate != "" {
		start, err = parseDateOnly(startDate)
		if err != nil {
			return nil, errors.NewAppErrorWithErr(errors.ErrBadRequest.Code, "Invalid startDate (expected YYYY-MM-DD)", err)
		}
	} else {
		// Default: start of week (Sunday)
		now := time.Now()
		start = now.AddDate(0, 0, -int(now.Weekday()))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	}

	end := start.AddDate(0, 0, 7)
	meals, err := s.repo.GetRangeByUserID(userID, start, end)
	if err != nil {
		return nil, err
	}

	out := WeeklyMealPlanResponse{}
	for _, mp := range meals {
		dateKey := mp.Date.Format("2006-01-02")
		if out[dateKey] == nil {
			out[dateKey] = map[string]*MealPlan{}
		}
		out[dateKey][mp.MealType] = mp
	}
	return out, nil
}

// Upsert creates or updates a meal plan slot for (user_id, date, meal_type).
func (s *Service) Upsert(userID uuid.UUID, householdID *uuid.UUID, req *UpsertMealPlanRequest) (*MealPlan, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}

	date, err := parseDateOnly(req.Date)
	if err != nil {
		return nil, errors.NewAppErrorWithErr(errors.ErrBadRequest.Code, "Invalid date (expected YYYY-MM-DD)", err)
	}

	existingID, err := s.repo.FindByUserDateType(userID, date, req.MealType)
	if err != nil {
		return nil, err
	}

	if existingID != nil {
		mp, err := s.repo.GetByID(*existingID)
		if err != nil {
			return nil, err
		}
		// ownership is enforced by query, but double-check
		if mp.UserID != userID {
			return nil, errors.ErrForbidden
		}
		mp.Name = req.Name
		mp.Description = req.Description
		mp.Ingredients = req.Ingredients
		mp.Servings = req.Servings
		if err := s.repo.Update(mp); err != nil {
			return nil, err
		}
		return mp, nil
	}

	mp := &MealPlan{
		ID:          uuid.New(),
		UserID:      userID,
		HouseholdID: householdID,
		Date:        date,
		MealType:    req.MealType,
		Name:        req.Name,
		Description: req.Description,
		Ingredients: req.Ingredients,
		Servings:    req.Servings,
	}
	if err := s.repo.Create(mp); err != nil {
		return nil, err
	}
	return mp, nil
}

func (s *Service) Delete(id uuid.UUID, userID uuid.UUID) error {
	mp, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if mp.UserID != userID {
		return errors.ErrForbidden
	}
	return s.repo.Delete(id)
}

