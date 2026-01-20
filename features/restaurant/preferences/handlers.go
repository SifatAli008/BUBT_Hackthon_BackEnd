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

// Get handles GET /api/v1/restaurant/preferences
// @Summary      Get preferences
// @Description  Get restaurant preferences for the authenticated user
// @Tags         restaurant-preferences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  RestaurantPreferences
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /restaurant/preferences [get]
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
	prefs, err := h.service.GetByUserID(userID)
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

// CreateOrUpdate handles POST /api/v1/restaurant/preferences
// @Summary      Create or update preferences
// @Description  Create or update restaurant preferences
// @Tags         restaurant-preferences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateRestaurantPreferencesRequest  true  "Preferences data"
// @Success      200      {object}  RestaurantPreferences
// @Success      201      {object}  RestaurantPreferences
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /restaurant/preferences [post]
func (h *Handler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateRestaurantPreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	prefs, err := h.service.CreateOrUpdate(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to save preferences", err.Error())
		return
	}
	utils.OKResponse(w, "Preferences saved successfully", prefs)
}
