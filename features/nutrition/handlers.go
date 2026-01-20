package nutrition

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// GetAll handles GET /api/v1/nutrition
// @Summary      Get nutrition data
// @Description  Get nutrition data for the authenticated user within a date range
// @Tags         nutrition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date  query     string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "End date (YYYY-MM-DD)"
// @Success      200         {array}   NutritionData
// @Failure      401         {object}  errors.AppError
// @Router       /nutrition [get]
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
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()
	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = t
		}
	}
	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = t
		}
	}
	data, err := h.service.GetByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve nutrition data", err.Error())
		return
	}
	utils.OKResponse(w, "Nutrition data retrieved successfully", data)
}

// GetToday handles GET /api/v1/nutrition/today
// @Summary      Get today's nutrition
// @Description  Get today's nutrition data for the authenticated user
// @Tags         nutrition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  NutritionData
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /nutrition/today [get]
func (h *Handler) GetToday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	data, err := h.service.GetToday(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve today's nutrition", err.Error())
		return
	}
	utils.OKResponse(w, "Today's nutrition retrieved successfully", data)
}

// GetByID handles GET /api/v1/nutrition/:id
// @Summary      Get nutrition data by ID
// @Description  Get details of a specific nutrition data entry
// @Tags         nutrition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Nutrition Data ID"
// @Success      200  {object}  NutritionData
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /nutrition/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/nutrition/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	data, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve nutrition data", err.Error())
		return
	}
	utils.OKResponse(w, "Nutrition data retrieved successfully", data)
}

// Create handles POST /api/v1/nutrition
// @Summary      Log nutrition data
// @Description  Create a new nutrition data entry
// @Tags         nutrition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateNutritionDataRequest  true  "Nutrition data"
// @Success      201      {object}  NutritionData
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /nutrition [post]
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
	var req CreateNutritionDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	data, err := h.service.Create(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create nutrition data", err.Error())
		return
	}
	utils.CreatedResponse(w, "Nutrition data created successfully", data)
}

// Update handles PUT /api/v1/nutrition/:id
// @Summary      Update nutrition data
// @Description  Update an existing nutrition data entry
// @Tags         nutrition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Nutrition Data ID"
// @Param        request  body      UpdateNutritionDataRequest true  "Nutrition data"
// @Success      200      {object}  NutritionData
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /nutrition/{id} [put]
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
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/nutrition/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateNutritionDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	data, err := h.service.Update(id, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update nutrition data", err.Error())
		return
	}
	utils.OKResponse(w, "Nutrition data updated successfully", data)
}

// GetStats handles GET /api/v1/nutrition/stats
// @Summary      Get nutrition statistics
// @Description  Get nutrition statistics for the authenticated user
// @Tags         nutrition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Number of days (default: 30)"
// @Success      200   {object}  NutritionStats
// @Failure      401   {object}  errors.AppError
// @Router       /nutrition/stats [get]
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
	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}
	stats, err := h.service.GetStats(userID, days)
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
