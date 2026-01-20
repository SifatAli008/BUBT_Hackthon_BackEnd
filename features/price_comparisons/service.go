package price_comparisons

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

func (s *Service) GetAll() ([]*PriceComparison, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id uuid.UUID) (*PriceComparison, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(req *CreatePriceComparisonRequest) (*PriceComparison, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	storesJSON := make([]interface{}, len(req.Stores))
	for i, store := range req.Stores {
		storesJSON[i] = store
	}
	c := &PriceComparison{
		ID:        uuid.New(),
		ItemName:  req.ItemName,
		Category:  req.Category,
		Stores:    JSONB{"stores": storesJSON},
		BestPrice: JSONB(req.BestPrice),
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) Update(id uuid.UUID, req *UpdatePriceComparisonRequest) (*PriceComparison, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.ItemName != "" {
		c.ItemName = req.ItemName
	}
	if req.Category != "" {
		c.Category = req.Category
	}
	if req.Stores != nil {
		storesJSON := make([]interface{}, len(req.Stores))
		for i, store := range req.Stores {
			storesJSON[i] = store
		}
		c.Stores = JSONB{"stores": storesJSON}
	}
	if req.BestPrice != nil {
		c.BestPrice = JSONB(req.BestPrice)
	}
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}
