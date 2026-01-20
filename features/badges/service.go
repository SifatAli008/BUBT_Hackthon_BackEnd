package badges

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

func (s *Service) GetByUserID(userID uuid.UUID) ([]*Badge, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) GetAvailableBadges() []*AvailableBadge {
	// Define available badges
	return []*AvailableBadge{
		{BadgeID: "first-login", Name: "Welcome!", Description: "Logged in for the first time", Icon: "👋", XPReward: 10},
		{BadgeID: "inventory-master", Name: "Inventory Master", Description: "Added 50 items to inventory", Icon: "📦", XPReward: 50},
		{BadgeID: "zero-waste", Name: "Zero Waste Hero", Description: "No waste for 7 days", Icon: "🌱", XPReward: 100},
		{BadgeID: "meal-planner", Name: "Meal Planner", Description: "Created 10 meal plans", Icon: "🍽️", XPReward: 75},
		{BadgeID: "nutrition-tracker", Name: "Nutrition Tracker", Description: "Logged nutrition for 30 days", Icon: "📊", XPReward: 150},
		{BadgeID: "community-helper", Name: "Community Helper", Description: "Shared 5 surplus items", Icon: "🤝", XPReward: 200},
		{BadgeID: "level-10", Name: "Level 10 Achiever", Description: "Reached level 10", Icon: "⭐", XPReward: 500},
		{BadgeID: "level-25", Name: "Level 25 Champion", Description: "Reached level 25", Icon: "🏆", XPReward: 1000},
	}
}

func (s *Service) UnlockBadge(userID uuid.UUID, req *UnlockBadgeRequest) (*Badge, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	
	// Check if badge already exists
	existing, err := s.repo.GetByUserIDAndBadgeID(userID, req.BadgeID)
	if err == nil && existing != nil {
		return existing, errors.ErrAlreadyExists
	}
	
	badge := &Badge{
		ID:          uuid.New(),
		UserID:      userID,
		BadgeID:     req.BadgeID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		UnlockedAt:  time.Now(),
		XPReward:    req.XPReward,
	}
	
	if err := s.repo.Create(badge); err != nil {
		return nil, err
	}
	
	return badge, nil
}
