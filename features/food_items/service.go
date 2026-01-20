package food_items

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"

	"github.com/google/uuid"
)

// Service handles food items business logic
type Service struct {
	repo *Repository
}

// NewService creates a new food items service
func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// GetAll retrieves all food items
func (s *Service) GetAll() ([]*FoodItem, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a food item by ID
func (s *Service) GetByID(id uuid.UUID) (*FoodItem, error) {
	return s.repo.GetByID(id)
}

// Create creates a new food item
func (s *Service) Create(req *CreateFoodItemRequest) (*FoodItem, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(
			errors.ErrValidationFailed.Code,
			"Validation failed: "+validationErrors[0],
			nil,
		)
	}

	item := &FoodItem{
		ID:               uuid.New(),
		Name:             req.Name,
		Category:         req.Category,
		TypicalExpiryDays: req.TypicalExpiryDays,
		StorageTips:      req.StorageTips,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

// Update updates a food item
func (s *Service) Update(id uuid.UUID, req *UpdateFoodItemRequest) (*FoodItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(
			errors.ErrValidationFailed.Code,
			"Validation failed: "+validationErrors[0],
			nil,
		)
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.TypicalExpiryDays != nil {
		item.TypicalExpiryDays = *req.TypicalExpiryDays
	}
	if req.StorageTips != "" {
		item.StorageTips = req.StorageTips
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

// Delete deletes a food item
func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
