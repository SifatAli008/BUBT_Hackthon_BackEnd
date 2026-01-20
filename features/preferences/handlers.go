package preferences

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"

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

// Get handles GET /api/v1/preferences
// @Summary      Get family preferences
// @Description  Get family preferences for the authenticated user's household
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  FamilyPreferences
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /preferences [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	// Use user's household_id or user_id as household_id
	householdID := userID
	if user, ok := r.Context().Value("user").(*auth.User); ok && user.HouseholdID != nil {
		householdID = *user.HouseholdID
	}
	prefs, err := h.service.GetByHouseholdID(householdID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve preferences", err.Error())
		return
	}
	utils.OKResponse(w, "Preferences retrieved successfully", prefs)
}

// CreateOrUpdate handles POST/PUT /api/v1/preferences
// @Summary      Create or update preferences
// @Description  Create or update family preferences
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreatePreferencesRequest  true  "Preferences data"
// @Success      200      {object}  FamilyPreferences
// @Success      201      {object}  FamilyPreferences
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /preferences [post]
// @Router       /preferences [put]
func (h *Handler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	// Set household_id from user if not provided
	if req.HouseholdID == uuid.Nil {
		if user, ok := r.Context().Value("user").(*auth.User); ok && user.HouseholdID != nil {
			req.HouseholdID = *user.HouseholdID
		} else {
			req.HouseholdID = userID
		}
	}
	prefs, err := h.service.CreateOrUpdate(&req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to save preferences", err.Error())
		return
	}
	if r.Method == http.MethodPost {
		utils.CreatedResponse(w, "Preferences created successfully", prefs)
	} else {
		utils.OKResponse(w, "Preferences updated successfully", prefs)
	}
}
