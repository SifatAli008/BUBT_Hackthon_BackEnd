package xp

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"
	"strconv"

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

// GetUserXP handles GET /api/v1/xp
// @Summary      Get user XP
// @Description  Get XP and level information for the authenticated user
// @Tags         xp
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UserXP
// @Failure      401  {object}  errors.AppError
// @Router       /xp [get]
func (h *Handler) GetUserXP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	xp, err := h.service.GetByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve XP", err.Error())
		return
	}
	utils.OKResponse(w, "XP retrieved successfully", xp)
}

// AddXP handles POST /api/v1/xp/add
// @Summary      Add XP
// @Description  Add XP to the authenticated user (system endpoint)
// @Tags         xp
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      AddXPRequest  true  "XP data"
// @Success      200      {object}  UserXP
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /xp/add [post]
func (h *Handler) AddXP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req AddXPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	if validationErrors := utils.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.BadRequestResponse(w, "Validation failed: "+validationErrors[0], nil)
		return
	}
	xp, err := h.service.AddXP(userID, req.Amount)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to add XP", err.Error())
		return
	}
	utils.OKResponse(w, "XP added successfully", xp)
}

// GetLeaderboard handles GET /api/v1/xp/leaderboard
// @Summary      Get leaderboard
// @Description  Get XP leaderboard
// @Tags         xp
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query     int  false  "Number of entries (default: 10, max: 100)"
// @Success      200    {array}   LeaderboardEntry
// @Failure      401    {object}  errors.AppError
// @Router       /xp/leaderboard [get]
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	leaderboard, err := h.service.GetLeaderboard(limit)
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
