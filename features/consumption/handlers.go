package consumption

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

// GetAll handles GET /api/v1/consumption
// @Summary      Get consumption logs
// @Description  Get all consumption logs for the authenticated user
// @Tags         consumption
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   ConsumptionLog
// @Failure      401  {object}  errors.AppError
// @Router       /consumption [get]
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
		utils.InternalServerErrorResponse(w, "Failed to retrieve consumption logs", err.Error())
		return
	}
	utils.OKResponse(w, "Consumption logs retrieved successfully", logs)
}

// GetByID handles GET /api/v1/consumption/:id
// @Summary      Get consumption log by ID
// @Description  Get details of a specific consumption log
// @Tags         consumption
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Consumption Log ID"
// @Success      200  {object}  ConsumptionLog
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /consumption/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/consumption/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	log, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve consumption log", err.Error())
		return
	}
	utils.OKResponse(w, "Consumption log retrieved successfully", log)
}

// Create handles POST /api/v1/consumption
// @Summary      Log consumption
// @Description  Create a new consumption log entry
// @Tags         consumption
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateConsumptionLogRequest  true  "Consumption log data"
// @Success      201      {object}  ConsumptionLog
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /consumption [post]
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
	var req CreateConsumptionLogRequest
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
		utils.InternalServerErrorResponse(w, "Failed to create consumption log", err.Error())
		return
	}
	utils.CreatedResponse(w, "Consumption log created successfully", log)
}

// Update handles PUT /api/v1/consumption/:id
// @Summary      Update consumption log
// @Description  Update an existing consumption log
// @Tags         consumption
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                      true  "Consumption Log ID"
// @Param        request  body      UpdateConsumptionLogRequest true  "Consumption log data"
// @Success      200      {object}  ConsumptionLog
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /consumption/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/consumption/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateConsumptionLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	log, err := h.service.Update(id, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update consumption log", err.Error())
		return
	}
	utils.OKResponse(w, "Consumption log updated successfully", log)
}

// Delete handles DELETE /api/v1/consumption/:id
// @Summary      Delete consumption log
// @Description  Delete a consumption log entry
// @Tags         consumption
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Consumption Log ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.AppError
// @Failure      403  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /consumption/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/consumption/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to delete consumption log", err.Error())
		return
	}
	utils.OKResponse(w, "Consumption log deleted successfully", nil)
}

// GetStats handles GET /api/v1/consumption/stats
// @Summary      Get consumption statistics
// @Description  Get consumption statistics for the authenticated user
// @Tags         consumption
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  ConsumptionStats
// @Failure      401  {object}  errors.AppError
// @Router       /consumption/stats [get]
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	stats, err := h.service.GetStats(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve stats", err.Error())
		return
	}
	utils.OKResponse(w, "Stats retrieved successfully", stats)
}
