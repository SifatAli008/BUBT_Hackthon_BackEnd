package profiles

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) getUserID(r *http.Request) (uuid.UUID, error) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return user.ID, nil
}

// GetProfile handles GET /api/v1/community/profile
// @Summary      Get user's community profile
// @Description  Get the authenticated user's community profile
// @Tags         community-profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  CommunityProfile
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /community/profile [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	profile, err := h.service.GetByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve profile", err.Error())
		return
	}
	utils.OKResponse(w, "Profile retrieved successfully", profile)
}

// GetProfileByUsername handles GET /api/v1/community/profile/:username
// @Summary      Get profile by username
// @Description  Get a community profile by username
// @Tags         community-profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username  path      string  true  "Username"
// @Success      200       {object}  CommunityProfile
// @Failure      401       {object}  errors.AppError
// @Failure      404       {object}  errors.AppError
// @Router       /community/profile/{username} [get]
func (h *Handler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	username := strings.TrimPrefix(r.URL.Path, "/api/v1/community/profile/")
	profile, err := h.service.GetByUsername(username)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve profile", err.Error())
		return
	}
	utils.OKResponse(w, "Profile retrieved successfully", profile)
}

// CreateProfile handles POST /api/v1/community/profile
// @Summary      Create community profile
// @Description  Create a new community profile for the authenticated user
// @Tags         community-profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateProfileRequest  true  "Profile data"
// @Success      201      {object}  CommunityProfile
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/profile [post]
func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	profile, err := h.service.Create(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create profile", err.Error())
		return
	}
	utils.CreatedResponse(w, "Profile created successfully", profile)
}

// UpdateProfile handles PUT /api/v1/community/profile
// @Summary      Update community profile
// @Description  Update the authenticated user's community profile
// @Tags         community-profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      UpdateProfileRequest  true  "Profile data"
// @Success      200      {object}  CommunityProfile
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /community/profile [put]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	profile, err := h.service.Update(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update profile", err.Error())
		return
	}
	utils.OKResponse(w, "Profile updated successfully", profile)
}
