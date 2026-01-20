package profiles

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

func (s *Service) GetByUserID(userID uuid.UUID) (*CommunityProfile, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) GetByUsername(username string) (*CommunityProfile, error) {
	return s.repo.GetByUsername(username)
}

func (s *Service) Create(userID uuid.UUID, req *CreateProfileRequest) (*CommunityProfile, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	notificationsEnabled := true
	if req.NotificationsEnabled != nil {
		notificationsEnabled = *req.NotificationsEnabled
	}
	notifyOnClaim := true
	if req.NotifyOnClaim != nil {
		notifyOnClaim = *req.NotifyOnClaim
	}
	notifyOnMessages := true
	if req.NotifyOnMessages != nil {
		notifyOnMessages = *req.NotifyOnMessages
	}
	profile := &CommunityProfile{
		ID:                  uuid.New(),
		UserID:              userID,
		Username:            req.Username,
		AvatarURL:           req.AvatarURL,
		CommunityRole:       req.CommunityRole,
		Bio:                 req.Bio,
		PreferredItems:      req.PreferredItems,
		AvoidItems:          req.AvoidItems,
		DietaryRestrictions: req.DietaryRestrictions,
		Allergens:           req.Allergens,
		AcceptsHotMeals:     req.AcceptsHotMeals,
		DistancePreference:  req.DistancePreference,
		Visibility:          req.Visibility,
		NotificationsEnabled: notificationsEnabled,
		NotifyOnClaim:       notifyOnClaim,
		NotifyOnMessages:    notifyOnMessages,
	}
	if profile.CommunityRole == "" {
		profile.CommunityRole = "member"
	}
	if profile.Visibility == "" {
		profile.Visibility = "public"
	}
	if err := s.repo.Create(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *Service) Update(userID uuid.UUID, req *UpdateProfileRequest) (*CommunityProfile, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.AvatarURL != "" {
		profile.AvatarURL = req.AvatarURL
	}
	if req.CommunityRole != "" {
		profile.CommunityRole = req.CommunityRole
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}
	if req.PreferredItems != nil {
		profile.PreferredItems = req.PreferredItems
	}
	if req.AvoidItems != nil {
		profile.AvoidItems = req.AvoidItems
	}
	if req.DietaryRestrictions != nil {
		profile.DietaryRestrictions = req.DietaryRestrictions
	}
	if req.Allergens != nil {
		profile.Allergens = req.Allergens
	}
	if req.AcceptsHotMeals != nil {
		profile.AcceptsHotMeals = *req.AcceptsHotMeals
	}
	if req.DistancePreference != "" {
		profile.DistancePreference = req.DistancePreference
	}
	if req.Visibility != "" {
		profile.Visibility = req.Visibility
	}
	if req.NotificationsEnabled != nil {
		profile.NotificationsEnabled = *req.NotificationsEnabled
	}
	if req.NotifyOnClaim != nil {
		profile.NotifyOnClaim = *req.NotifyOnClaim
	}
	if req.NotifyOnMessages != nil {
		profile.NotifyOnMessages = *req.NotifyOnMessages
	}
	if err := s.repo.Update(profile); err != nil {
		return nil, err
	}
	return profile, nil
}
