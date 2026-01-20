package inventory

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Handler handles HTTP requests for inventory
type Handler struct {
	service *Service
}

// NewHandler creates a new inventory handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// getUserIDFromContext extracts user ID from request context
func (h *Handler) getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return user.ID, nil
}

// GetAll handles GET /api/v1/inventory
// @Summary      Get user inventory
// @Description  Get all inventory items for the authenticated user
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   InventoryItem
// @Failure      401  {object}  errors.AppError
// @Router       /inventory [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}

	items, err := h.service.GetAllByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve inventory", err.Error())
		return
	}

	utils.OKResponse(w, "Inventory retrieved successfully", items)
}

// GetByID handles GET /api/v1/inventory/:id
// @Summary      Get inventory item by ID
// @Description  Get details of a specific inventory item
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Inventory Item ID"
// @Success      200  {object}  InventoryItem
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /inventory/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/inventory/")
	id, err := uuid.Parse(idStr)
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
		utils.InternalServerErrorResponse(w, "Failed to retrieve inventory item", err.Error())
		return
	}

	utils.OKResponse(w, "Inventory item retrieved successfully", item)
}

// Create handles POST /api/v1/inventory
// @Summary      Add inventory item
// @Description  Add a new item to user's inventory
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateInventoryItemRequest  true  "Inventory item data"
// @Success      201      {object}  InventoryItem
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /inventory [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}

	var req CreateInventoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}

	item, err := h.service.Create(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create inventory item", err.Error())
		return
	}

	utils.CreatedResponse(w, "Inventory item created successfully", item)
}

// Update handles PUT /api/v1/inventory/:id
// @Summary      Update inventory item
// @Description  Update an existing inventory item
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Inventory Item ID"
// @Param        request  body      UpdateInventoryItemRequest true  "Inventory item data"
// @Success      200      {object}  InventoryItem
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /inventory/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/inventory/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}

	var req UpdateInventoryItemRequest
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
		utils.InternalServerErrorResponse(w, "Failed to update inventory item", err.Error())
		return
	}

	utils.OKResponse(w, "Inventory item updated successfully", item)
}

// Delete handles DELETE /api/v1/inventory/:id
// @Summary      Delete inventory item
// @Description  Delete an inventory item
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Inventory Item ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.AppError
// @Failure      403  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /inventory/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/inventory/")
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
		utils.InternalServerErrorResponse(w, "Failed to delete inventory item", err.Error())
		return
	}

	utils.OKResponse(w, "Inventory item deleted successfully", nil)
}

// GetExpiring handles GET /api/v1/inventory/expiring
// @Summary      Get expiring items
// @Description  Get inventory items expiring within specified days (default 7)
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Number of days (default: 7)"
// @Success      200   {array}   InventoryItem
// @Failure      401   {object}  errors.AppError
// @Router       /inventory/expiring [get]
func (h *Handler) GetExpiring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}

	days := 7 // default
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	items, err := h.service.GetExpiring(userID, days)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve expiring items", err.Error())
		return
	}

	utils.OKResponse(w, "Expiring items retrieved successfully", items)
}

// GetExpired handles GET /api/v1/inventory/expired
// @Summary      Get expired items
// @Description  Get all expired inventory items
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   InventoryItem
// @Failure      401  {object}  errors.AppError
// @Router       /inventory/expired [get]
func (h *Handler) GetExpired(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}

	items, err := h.service.GetExpired(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve expired items", err.Error())
		return
	}

	utils.OKResponse(w, "Expired items retrieved successfully", items)
}
