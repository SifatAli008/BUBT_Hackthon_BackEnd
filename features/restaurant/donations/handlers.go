package donations

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

// GetAll handles GET /api/v1/restaurant/donations
// @Summary      List donations
// @Description  Get all donation logs for the authenticated restaurant
// @Tags         restaurant-donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   DonationLog
// @Failure      401  {object}  errors.AppError
// @Router       /restaurant/donations [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	logs, err := h.service.GetAllByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve donations", err.Error())
		return
	}
	utils.OKResponse(w, "Donations retrieved successfully", logs)
}

// Create handles POST /api/v1/restaurant/donations
// @Summary      Log donation
// @Description  Create a new donation log entry
// @Tags         restaurant-donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateDonationLogRequest  true  "Donation log data"
// @Success      201      {object}  DonationLog
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /restaurant/donations [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateDonationLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	log, err := h.service.Create(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create donation log", err.Error())
		return
	}
	utils.CreatedResponse(w, "Donation logged successfully", log)
}

// GetImpact handles GET /api/v1/restaurant/impact
// @Summary      Get impact metrics
// @Description  Get impact metrics for the authenticated restaurant
// @Tags         restaurant-impact
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  ImpactMetrics
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /restaurant/impact [get]
func (h *Handler) GetImpact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	impact, err := h.service.GetImpact(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve impact metrics", err.Error())
		return
	}
	utils.OKResponse(w, "Impact metrics retrieved successfully", impact)
}
