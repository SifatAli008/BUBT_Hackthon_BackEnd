package shopping_list

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) GetAllByUserID(userID uuid.UUID, includePurchased bool) ([]*ShoppingListItem, error) {
	return s.repo.GetAllByUserID(userID, includePurchased)
}

func (s *Service) GetByID(id uuid.UUID) (*ShoppingListItem, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, householdID *uuid.UUID, req *CreateShoppingListItemRequest) (*ShoppingListItem, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}
	item := &ShoppingListItem{
		ID:            uuid.New(),
		UserID:        userID,
		HouseholdID:   householdID,
		Name:          req.Name,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		Category:      req.Category,
		Priority:      priority,
		Purchased:     false,
		PurchasedAt:   nil,
		EstimatedPrice: req.EstimatedPrice,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateShoppingListItemRequest) (*ShoppingListItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
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
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Priority != "" {
		item.Priority = req.Priority
	}
	if req.EstimatedPrice != nil {
		item.EstimatedPrice = req.EstimatedPrice
	}
	if req.Purchased != nil {
		item.Purchased = *req.Purchased
		if item.Purchased {
			now := time.Now()
			item.PurchasedAt = &now
		} else {
			item.PurchasedAt = nil
		}
	}
	// If PurchasedAt explicitly provided (rare), accept it.
	if req.PurchasedAt != nil {
		item.PurchasedAt = req.PurchasedAt
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Delete(id uuid.UUID, userID uuid.UUID) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return errors.ErrForbidden
	}
	return s.repo.Delete(id)
}

func (s *Service) TogglePurchased(id uuid.UUID, userID uuid.UUID) (*ShoppingListItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item.UserID != userID {
		return nil, errors.ErrForbidden
	}
	item.Purchased = !item.Purchased
	if item.Purchased {
		now := time.Now()
		item.PurchasedAt = &now
	} else {
		item.PurchasedAt = nil
	}
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

// ComputeMissingFromInventory adds low-stock inventory items (quantity < 1) to shopping list if missing.
func (s *Service) ComputeMissingFromInventory(userID uuid.UUID, householdID *uuid.UUID) ([]*ShoppingListItem, error) {
	existing, err := s.repo.ListUnpurchasedNamesLower(userID)
	if err != nil {
		return nil, err
	}
	lowStock, err := s.repo.GetLowStockInventoryNames(userID)
	if err != nil {
		return nil, err
	}

	for _, it := range lowStock {
		name := strings.TrimSpace(it.Name)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := existing[key]; ok {
			continue
		}
		existing[key] = struct{}{}

		priority := "medium"
		item := &ShoppingListItem{
			ID:          uuid.New(),
			UserID:      userID,
			HouseholdID: householdID,
			Name:        name,
			Quantity:    1,
			Priority:    priority,
			Purchased:   false,
		}
		if it.Unit.Valid {
			item.Unit = it.Unit.String
		}
		if it.Category.Valid {
			item.Category = it.Category.String
		}
		if err := s.repo.Create(item); err != nil {
			return nil, err
		}
	}

	// Return updated list (unpurchased)
	return s.repo.GetAllByUserID(userID, false)
}

