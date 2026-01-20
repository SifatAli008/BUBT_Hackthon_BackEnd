package surplus

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

func (s *Service) GetAll(status string) ([]*SurplusPost, error) {
	return s.repo.GetAll(status)
}

func (s *Service) GetByID(id uuid.UUID) (*SurplusPost, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, userName string, avatarURL string, req *CreateSurplusPostRequest) (*SurplusPost, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	post := &SurplusPost{
		ID:            uuid.New(),
		UserID:        userID,
		UserName:      userName,
		AvatarURL:     avatarURL,
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Tags:          req.Tags,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		PickupWindow:  JSONB(req.PickupWindow),
		PickupLocation: req.PickupLocation,
		DistanceKm:    req.DistanceKm,
		Image:         req.Image,
		Status:        "available",
		ExpiresAt:     req.ExpiresAt,
	}
	if err := s.repo.Create(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateSurplusPostRequest) (*SurplusPost, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if post.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Description != "" {
		post.Description = req.Description
	}
	if req.Category != "" {
		post.Category = req.Category
	}
	if req.Tags != nil {
		post.Tags = req.Tags
	}
	if req.Quantity != nil {
		post.Quantity = *req.Quantity
	}
	if req.Unit != "" {
		post.Unit = req.Unit
	}
	if req.PickupWindow != nil {
		post.PickupWindow = JSONB(req.PickupWindow)
	}
	if req.PickupLocation != "" {
		post.PickupLocation = req.PickupLocation
	}
	if req.Status != "" {
		post.Status = req.Status
	}
	if req.Image != "" {
		post.Image = req.Image
	}
	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) Delete(id uuid.UUID, userID uuid.UUID) error {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.ErrForbidden
	}
	return s.repo.Delete(id)
}

func (s *Service) CreateRequest(postID uuid.UUID, userID uuid.UUID, userName string, req *CreateSurplusRequestRequest) (*SurplusRequest, error) {
	_, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	request := &SurplusRequest{
		ID:       uuid.New(),
		PostID:   postID,
		UserID:   userID,
		UserName: userName,
		Message:  req.Message,
		Status:   "pending",
	}
	if err := s.repo.CreateRequest(request); err != nil {
		return nil, err
	}
	return request, nil
}

func (s *Service) GetRequestsByPostID(postID uuid.UUID, userID uuid.UUID) ([]*SurplusRequest, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	if post.UserID != userID {
		return nil, errors.ErrForbidden
	}
	return s.repo.GetRequestsByPostID(postID)
}

func (s *Service) UpdateRequest(requestID uuid.UUID, postID uuid.UUID, userID uuid.UUID, req *UpdateSurplusRequestRequest) (*SurplusRequest, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	if post.UserID != userID {
		return nil, errors.ErrForbidden
	}
	request, err := s.repo.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}
	if request.PostID != postID {
		return nil, errors.ErrNotFound
	}
	request.Status = req.Status
	if err := s.repo.UpdateRequest(request); err != nil {
		return nil, err
	}
	return request, nil
}

func (s *Service) CreateComment(postID uuid.UUID, userID uuid.UUID, userName string, req *CreateSurplusCommentRequest) (*SurplusComment, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	_, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	comment := &SurplusComment{
		ID:       uuid.New(),
		PostID:   postID,
		UserID:   userID,
		UserName: userName,
		Message:  req.Message,
	}
	if err := s.repo.CreateComment(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *Service) GetCommentsByPostID(postID uuid.UUID) ([]*SurplusComment, error) {
	return s.repo.GetCommentsByPostID(postID)
}
