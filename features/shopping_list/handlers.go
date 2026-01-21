package shopping_list

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
	// fall back to user id as household
	hh := userID
	return userID, &hh, nil
}

// GetAll handles GET /api/v1/shopping-list
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	includePurchased := r.URL.Query().Get("include_purchased") == "1" || strings.ToLower(r.URL.Query().Get("include_purchased")) == "true"
	items, err := h.service.GetAllByUserID(userID, includePurchased)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve shopping list", err.Error())
		return
	}
	utils.OKResponse(w, "Shopping list retrieved successfully", items)
}

// GetByID handles GET /api/v1/shopping-list/:id
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/shopping-list/")
	id, err := uuid.Parse(strings.Split(idStr, "/")[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	item, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve shopping list item", err.Error())
		return
	}
	utils.OKResponse(w, "Shopping list item retrieved successfully", item)
}

// Create handles POST /api/v1/shopping-list
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, householdID, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateShoppingListItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	item, err := h.service.Create(userID, householdID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create shopping list item", err.Error())
		return
	}
	utils.CreatedResponse(w, "Shopping list item created successfully", item)
}

// Update handles PUT /api/v1/shopping-list/:id
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/shopping-list/")
	id, err := uuid.Parse(strings.Split(idStr, "/")[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateShoppingListItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	item, err := h.service.Update(id, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update shopping list item", err.Error())
		return
	}
	utils.OKResponse(w, "Shopping list item updated successfully", item)
}

// Delete handles DELETE /api/v1/shopping-list/:id
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
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/shopping-list/")
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
		utils.InternalServerErrorResponse(w, "Failed to delete shopping list item", err.Error())
		return
	}
	utils.OKResponse(w, "Shopping list item deleted successfully", nil)
}

// Toggle handles PUT /api/v1/shopping-list/:id/toggle
func (h *Handler) Toggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch && r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/shopping-list/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "toggle" {
		utils.NotFoundResponse(w, "Not found")
		return
	}
	id, err := uuid.Parse(parts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	item, err := h.service.TogglePurchased(id, userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to toggle shopping list item", err.Error())
		return
	}
	utils.OKResponse(w, "Shopping list item toggled successfully", item)
}

// ComputeMissing handles POST /api/v1/shopping-list/compute-missing
func (h *Handler) ComputeMissing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, householdID, err := h.getUserAndHousehold(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	items, err := h.service.ComputeMissingFromInventory(userID, householdID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to compute missing items", err.Error())
		return
	}
	utils.OKResponse(w, "Shopping list updated successfully", items)
}

