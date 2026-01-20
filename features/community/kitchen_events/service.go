package kitchen_events

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

func (s *Service) GetAll(status string) ([]*KitchenEvent, error) {
	return s.repo.GetAll(status)
}

func (s *Service) GetByID(id uuid.UUID) (*KitchenEvent, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(req *CreateKitchenEventRequest) (*KitchenEvent, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	event := &KitchenEvent{
		ID:              uuid.New(),
		Title:           req.Title,
		Description:     req.Description,
		Date:            req.Date,
		Time:            req.Time,
		Location:        req.Location,
		Tags:            req.Tags,
		VolunteersNeeded: req.VolunteersNeeded,
		Volunteers:      JSONB{},
		FoodSavedKg:     0,
		Status:          "upcoming",
		Image:           req.Image,
	}
	if err := s.repo.Create(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) Update(id uuid.UUID, req *UpdateKitchenEventRequest) (*KitchenEvent, error) {
	event, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Date != nil {
		event.Date = *req.Date
	}
	if req.Time != "" {
		event.Time = req.Time
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.Tags != nil {
		event.Tags = req.Tags
	}
	if req.VolunteersNeeded != nil {
		event.VolunteersNeeded = *req.VolunteersNeeded
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if req.FoodSavedKg != nil {
		event.FoodSavedKg = *req.FoodSavedKg
	}
	if req.Image != "" {
		event.Image = req.Image
	}
	if err := s.repo.Update(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) AddVolunteer(eventID uuid.UUID, userID uuid.UUID, userName string, role string) (*KitchenEvent, error) {
	event, err := s.repo.GetByID(eventID)
	if err != nil {
		return nil, err
	}
	volunteers := []map[string]interface{}{}
	if event.Volunteers != nil {
		if v, ok := event.Volunteers["volunteers"].([]interface{}); ok {
			for _, vol := range v {
				if vm, ok := vol.(map[string]interface{}); ok {
					volunteers = append(volunteers, vm)
				}
			}
		}
	}
	volunteer := map[string]interface{}{
		"id":        uuid.New().String(),
		"userId":    userID.String(),
		"name":      userName,
		"role":      role,
		"avatarUrl": "",
	}
	volunteers = append(volunteers, volunteer)
	event.Volunteers = JSONB{"volunteers": volunteers}
	if err := s.repo.Update(event); err != nil {
		return nil, err
	}
	return event, nil
}
