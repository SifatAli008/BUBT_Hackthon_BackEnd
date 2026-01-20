package consumption

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

func (s *Service) GetAllByUserID(userID uuid.UUID) ([]*ConsumptionLog, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *Service) GetByID(id uuid.UUID) (*ConsumptionLog, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, req *CreateConsumptionLogRequest) (*ConsumptionLog, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	log := &ConsumptionLog{
		ID:              uuid.New(),
		UserID:          userID,
		InventoryItemID: req.InventoryItemID,
		FoodName:        req.FoodName,
		Quantity:        req.Quantity,
		Unit:            req.Unit,
		Category:        req.Category,
		ConsumedAt:      req.ConsumedAt,
		WasWasted:       req.WasWasted,
		Notes:           req.Notes,
	}
	if log.ConsumedAt.IsZero() {
		log.ConsumedAt = time.Now()
	}
	if err := s.repo.Create(log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateConsumptionLogRequest) (*ConsumptionLog, error) {
	log, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if log.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.FoodName != "" {
		log.FoodName = req.FoodName
	}
	if req.Quantity != nil {
		log.Quantity = *req.Quantity
	}
	if req.Unit != "" {
		log.Unit = req.Unit
	}
	if req.Category != "" {
		log.Category = req.Category
	}
	if req.ConsumedAt != nil {
		log.ConsumedAt = *req.ConsumedAt
	}
	if req.WasWasted != nil {
		log.WasWasted = *req.WasWasted
	}
	if req.Notes != "" {
		log.Notes = req.Notes
	}
	if err := s.repo.Update(log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *Service) Delete(id uuid.UUID, userID uuid.UUID) error {
	log, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if log.UserID != userID {
		return errors.ErrForbidden
	}
	return s.repo.Delete(id)
}

func (s *Service) GetStats(userID uuid.UUID) (*ConsumptionStats, error) {
	return s.repo.GetStats(userID)
}
