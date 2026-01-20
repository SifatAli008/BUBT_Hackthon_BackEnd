package leftovers

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

func (s *Service) GetAll(status string) ([]*LeftoverItem, error) {
	return s.repo.GetAll(status)
}

func (s *Service) GetByID(id uuid.UUID) (*LeftoverItem, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, userName string, req *CreateLeftoverItemRequest) (*LeftoverItem, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	item := &LeftoverItem{
		ID:           uuid.New(),
		UserID:       userID,
		UserName:     userName,
		DishName:     req.DishName,
		Description:  req.Description,
		Portions:     req.Portions,
		DistanceKm:   req.DistanceKm,
		DietaryTags:  req.DietaryTags,
		Allergens:    req.Allergens,
		PickupWindow: req.PickupWindow,
		Status:       "available",
		Image:        req.Image,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateLeftoverItemRequest) (*LeftoverItem, error) {
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
	if req.DishName != "" {
		item.DishName = req.DishName
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Portions != nil {
		item.Portions = *req.Portions
	}
	if req.DistanceKm != nil {
		item.DistanceKm = *req.DistanceKm
	}
	if req.DietaryTags != nil {
		item.DietaryTags = req.DietaryTags
	}
	if req.Allergens != nil {
		item.Allergens = req.Allergens
	}
	if req.PickupWindow != "" {
		item.PickupWindow = req.PickupWindow
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Image != "" {
		item.Image = req.Image
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

func (s *Service) CreateClaim(leftoverID uuid.UUID, userID uuid.UUID, userName string, req *CreateLeftoverClaimRequest) (*LeftoverClaim, error) {
	_, err := s.repo.GetByID(leftoverID)
	if err != nil {
		return nil, err
	}
	claim := &LeftoverClaim{
		ID:            uuid.New(),
		LeftoverItemID: leftoverID,
		UserID:        userID,
		UserName:      userName,
		Message:       req.Message,
	}
	if err := s.repo.CreateClaim(claim); err != nil {
		return nil, err
	}
	return claim, nil
}

func (s *Service) GetClaimsByLeftoverID(leftoverID uuid.UUID) ([]*LeftoverClaim, error) {
	return s.repo.GetClaimsByLeftoverID(leftoverID)
}
