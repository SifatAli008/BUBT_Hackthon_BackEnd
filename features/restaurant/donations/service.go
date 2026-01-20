package donations

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

func (s *Service) GetAllByUserID(userID uuid.UUID) ([]*DonationLog, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *Service) Create(userID uuid.UUID, req *CreateDonationLogRequest) (*DonationLog, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	log := &DonationLog{
		ID:            uuid.New(),
		UserID:        userID,
		Date:          req.Date,
		RecipientType: req.RecipientType,
		RecipientName: req.RecipientName,
		Items:         req.Items,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		MealsProvided: req.MealsProvided,
		CO2SavedKg:    req.CO2SavedKg,
		Notes:         req.Notes,
	}
	if err := s.repo.Create(log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *Service) GetImpact(userID uuid.UUID) (*ImpactMetrics, error) {
	return s.repo.GetImpactByUserID(userID)
}
