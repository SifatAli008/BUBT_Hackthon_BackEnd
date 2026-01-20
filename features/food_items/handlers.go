package food_items

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Handler handles HTTP requests for food items
type Handler struct {
	service *Service
}

// NewHandler creates a new food items handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetAll handles GET /api/v1/food-items
// @Summary      List all food items
// @Description  Get a list of all food items (reference data)
// @Tags         food-items
// @Accept       json
// @Produce      json
// @Success      200  {array}   FoodItem
// @Router       /food-items [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	items, err := h.service.GetAll()
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve food items", err.Error())
		return
	}

	utils.OKResponse(w, "Food items retrieved successfully", items)
}

// GetByID handles GET /api/v1/food-items/:id
// @Summary      Get food item by ID
// @Description  Get details of a specific food item
// @Tags         food-items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Food Item ID"
// @Success      200  {object}  FoodItem
// @Failure      404  {object}  errors.AppError
// @Router       /food-items/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/food-items/")
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
		utils.InternalServerErrorResponse(w, "Failed to retrieve food item", err.Error())
		return
	}

	utils.OKResponse(w, "Food item retrieved successfully", item)
}

// Create handles POST /api/v1/food-items
// @Summary      Create food item
// @Description  Create a new food item (admin only)
// @Tags         food-items
// @Accept       json
// @Produce      json
// @Param        request  body      CreateFoodItemRequest  true  "Food item data"
// @Success      201      {object}  FoodItem
// @Failure      400      {object}  errors.AppError
// @Router       /food-items [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	var req CreateFoodItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}

	item, err := h.service.Create(&req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create food item", err.Error())
		return
	}

	utils.CreatedResponse(w, "Food item created successfully", item)
}

// Update handles PUT /api/v1/food-items/:id
// @Summary      Update food item
// @Description  Update an existing food item (admin only)
// @Tags         food-items
// @Accept       json
// @Produce      json
// @Param        id       path      string                  true  "Food Item ID"
// @Param        request  body      UpdateFoodItemRequest   true  "Food item data"
// @Success      200      {object}  FoodItem
// @Failure      400      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /food-items/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/food-items/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}

	var req UpdateFoodItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}

	item, err := h.service.Update(id, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update food item", err.Error())
		return
	}

	utils.OKResponse(w, "Food item updated successfully", item)
}

// Delete handles DELETE /api/v1/food-items/:id
// @Summary      Delete food item
// @Description  Delete a food item (admin only)
// @Tags         food-items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Food Item ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  errors.AppError
// @Router       /food-items/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/food-items/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to delete food item", err.Error())
		return
	}

	utils.OKResponse(w, "Food item deleted successfully", nil)
}
