package inventory

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"

	"github.com/google/uuid"
)

// Service handles inventory business logic
type Service struct {
	repo *Repository
}

// NewService creates a new inventory service
func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// GetAllByUserID retrieves all inventory items for a user
func (s *Service) GetAllByUserID(userID uuid.UUID) ([]*InventoryItem, error) {
	return s.repo.GetAllByUserID(userID)
}

// GetByID retrieves an inventory item by ID
func (s *Service) GetByID(id uuid.UUID) (*InventoryItem, error) {
	return s.repo.GetByID(id)
}

// Create creates a new inventory item
func (s *Service) Create(userID uuid.UUID, req *CreateInventoryItemRequest) (*InventoryItem, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(
			errors.ErrValidationFailed.Code,
			"Validation failed: "+validationErrors[0],
			nil,
		)
	}

	item := &InventoryItem{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       req.Name,
		Quantity:   req.Quantity,
		Unit:       req.Unit,
		ExpiryDate: req.ExpiryDate,
		Category:   req.Category,
		Location:   req.Location,
		FoodItemID: req.FoodItemID,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

// Update updates an inventory item
func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateInventoryItemRequest) (*InventoryItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if item.UserID != userID {
		return nil, errors.ErrForbidden
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
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.ExpiryDate != nil {
		item.ExpiryDate = req.ExpiryDate
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Location != "" {
		item.Location = req.Location
	}
	if req.FoodItemID != nil {
		item.FoodItemID = req.FoodItemID
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

// Delete deletes an inventory item
func (s *Service) Delete(id uuid.UUID, userID uuid.UUID) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Check ownership
	if item.UserID != userID {
		return errors.ErrForbidden
	}

	return s.repo.Delete(id)
}

// GetExpiring retrieves items expiring within specified days (default 7)
func (s *Service) GetExpiring(userID uuid.UUID, days int) ([]*InventoryItem, error) {
	if days <= 0 {
		days = 7
	}
	return s.repo.GetExpiring(userID, days)
}

// GetExpired retrieves expired items
func (s *Service) GetExpired(userID uuid.UUID) ([]*InventoryItem, error) {
	return s.repo.GetExpired(userID)
}
