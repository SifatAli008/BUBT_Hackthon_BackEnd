package meal_plans

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

func (h *Handler) getUserAndHousehold(r *http.Request) (uuid.UUID, *uuid.UUID, error) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		return uuid.Nil, nil, errors.ErrUnauthorized
	}
	userID := user.ID
	if user.HouseholdID != nil {
		return userID, user.HouseholdID, nil
	}
	hh := userID
	return userID, &hh, nil
}

// GetWeekly handles GET /api/v1/meal-plans/weekly
func (h *Handler) GetWeekly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	startDate := r.URL.Query().Get("startDate")
	plan, err := h.service.GetWeeklyByUserID(userID, startDate)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve meal plan", err.Error())
		return
	}
	utils.OKResponse(w, "Meal plan retrieved successfully", plan)
}

// Upsert handles POST /api/v1/meal-plans
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, householdID, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req UpsertMealPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	item, err := h.service.Upsert(userID, householdID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to save meal plan", err.Error())
		return
	}
	utils.OKResponse(w, "Meal plan saved successfully", item)
}

// Delete handles DELETE /api/v1/meal-plans/:id
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/meal-plans/")
	id, err := uuid.Parse(strings.Split(idStr, "/")[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to delete meal plan", err.Error())
		return
	}
	utils.OKResponse(w, "Meal plan deleted successfully", nil)
}

