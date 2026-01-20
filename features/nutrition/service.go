package nutrition

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

func (s *Service) GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate time.Time) ([]*NutritionData, error) {
	return s.repo.GetByUserIDAndDateRange(userID, startDate, endDate)
}

func (s *Service) GetToday(userID uuid.UUID) (*NutritionData, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.repo.GetByUserIDAndDate(userID, today)
}

func (s *Service) GetByID(id uuid.UUID) (*NutritionData, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, req *CreateNutritionDataRequest) (*NutritionData, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	date := req.Date
	if date.IsZero() {
		date = time.Now().Truncate(24 * time.Hour)
	}
	d := &NutritionData{
		ID:             uuid.New(),
		UserID:         userID,
		Date:           date,
		Calories:       req.Calories,
		Protein:        req.Protein,
		Carbs:          req.Carbs,
		Fats:           req.Fats,
		Fiber:          req.Fiber,
		Sugar:          req.Sugar,
		Sodium:         req.Sodium,
		VitaminA:       req.VitaminA,
		VitaminB:       req.VitaminB,
		VitaminC:       req.VitaminC,
		VitaminD:       req.VitaminD,
		Iron:           req.Iron,
		Calcium:        req.Calcium,
		NutritionScore: req.NutritionScore,
	}
	if err := s.repo.Create(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateNutritionDataRequest) (*NutritionData, error) {
	d, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if d.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.Calories != nil {
		d.Calories = *req.Calories
	}
	if req.Protein != nil {
		d.Protein = *req.Protein
	}
	if req.Carbs != nil {
		d.Carbs = *req.Carbs
	}
	if req.Fats != nil {
		d.Fats = *req.Fats
	}
	if req.Fiber != nil {
		d.Fiber = *req.Fiber
	}
	if req.Sugar != nil {
		d.Sugar = *req.Sugar
	}
	if req.Sodium != nil {
		d.Sodium = *req.Sodium
	}
	if req.VitaminA != nil {
		d.VitaminA = *req.VitaminA
	}
	if req.VitaminB != nil {
		d.VitaminB = *req.VitaminB
	}
	if req.VitaminC != nil {
		d.VitaminC = *req.VitaminC
	}
	if req.VitaminD != nil {
		d.VitaminD = *req.VitaminD
	}
	if req.Iron != nil {
		d.Iron = *req.Iron
	}
	if req.Calcium != nil {
		d.Calcium = *req.Calcium
	}
	if req.NutritionScore != nil {
		d.NutritionScore = req.NutritionScore
	}
	if err := s.repo.Update(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetStats(userID uuid.UUID, days int) (*NutritionStats, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetStats(userID, days)
}
