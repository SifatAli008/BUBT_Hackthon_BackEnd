package badges

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

// GetUserBadges handles GET /api/v1/badges
// @Summary      Get user badges
// @Description  Get all badges for the authenticated user
// @Tags         badges
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   Badge
// @Failure      401  {object}  errors.AppError
// @Router       /badges [get]
func (h *Handler) GetUserBadges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	badges, err := h.service.GetByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve badges", err.Error())
		return
	}
	utils.OKResponse(w, "Badges retrieved successfully", badges)
}

// GetAvailableBadges handles GET /api/v1/badges/available
// @Summary      Get available badges
// @Description  Get list of all available badges
// @Tags         badges
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   AvailableBadge
// @Failure      401  {object}  errors.AppError
// @Router       /badges/available [get]
func (h *Handler) GetAvailableBadges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	badges := h.service.GetAvailableBadges()
	utils.OKResponse(w, "Available badges retrieved successfully", badges)
}

// UnlockBadge handles POST /api/v1/badges/unlock
// @Summary      Unlock badge
// @Description  Unlock a badge for the authenticated user (system endpoint)
// @Tags         badges
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      UnlockBadgeRequest  true  "Badge unlock data"
// @Success      201      {object}  Badge
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      409      {object}  errors.AppError
// @Router       /badges/unlock [post]
func (h *Handler) UnlockBadge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req UnlockBadgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	badge, err := h.service.UnlockBadge(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to unlock badge", err.Error())
		return
	}
	utils.CreatedResponse(w, "Badge unlocked successfully", badge)
}
