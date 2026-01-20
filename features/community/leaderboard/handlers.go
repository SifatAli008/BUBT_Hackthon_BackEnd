package leaderboard

import (
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

// GetLeaderboard handles GET /api/v1/community/leaderboard
// @Summary      Get leaderboard
// @Description  Get community leaderboard by type
// @Tags         community-leaderboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        type  query     string  false  "Leaderboard type (top-sharers, zero-waste, volunteer-stars, building-impact, weekly-xp)"
// @Success      200   {object}  Leaderboard
// @Failure      401   {object}  errors.AppError
// @Router       /community/leaderboard [get]
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	leaderboardType := r.URL.Query().Get("type")
	leaderboard, err := h.service.GetLeaderboard(leaderboardType)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve leaderboard", err.Error())
		return
	}
	utils.OKResponse(w, "Leaderboard retrieved successfully", leaderboard)
}

// GetImpact handles GET /api/v1/community/impact
// @Summary      Get community impact
// @Description  Get overall community impact statistics
// @Tags         community-impact
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  Impact
// @Failure      401  {object}  errors.AppError
// @Router       /community/impact [get]
func (h *Handler) GetImpact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	impact, err := h.service.GetImpact()
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve impact", err.Error())
		return
	}
	utils.OKResponse(w, "Impact retrieved successfully", impact)
}

// GetPersonalImpact handles GET /api/v1/community/impact/personal
// @Summary      Get personal impact
// @Description  Get personal contribution statistics
// @Tags         community-impact
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  Impact
// @Failure      401  {object}  errors.AppError
// @Router       /community/impact/personal [get]
func (h *Handler) GetPersonalImpact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	impact, err := h.service.GetPersonalImpact(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve personal impact", err.Error())
		return
	}
	utils.OKResponse(w, "Personal impact retrieved successfully", impact)
}
