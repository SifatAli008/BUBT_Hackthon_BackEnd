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

func (s *Service) GetByUserID(userID uuid.UUID) (*RestaurantPreferences, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) CreateOrUpdate(userID uuid.UUID, req *CreateRestaurantPreferencesRequest) (*RestaurantPreferences, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	prefs := &RestaurantPreferences{
		ID:                 uuid.New(),
		UserID:             userID,
		CuisineType:        req.CuisineType,
		OperatingHours:     req.OperatingHours,
		DonationPreferences: req.DonationPreferences,
	}
	if err := s.repo.CreateOrUpdate(prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}
