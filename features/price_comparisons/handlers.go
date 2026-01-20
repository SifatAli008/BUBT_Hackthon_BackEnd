package price_comparisons

import (
	"encoding/json"
	"foodlink_backend/errors"
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

// GetAll handles GET /api/v1/price-comparisons
// @Summary      Get price comparisons
// @Description  Get all price comparisons
// @Tags         price-comparisons
// @Accept       json
// @Produce      json
// @Success      200  {array}   PriceComparison
// @Router       /price-comparisons [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	comparisons, err := h.service.GetAll()
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve price comparisons", err.Error())
		return
	}
	utils.OKResponse(w, "Price comparisons retrieved successfully", comparisons)
}

// GetByID handles GET /api/v1/price-comparisons/:id
// @Summary      Get price comparison by ID
// @Description  Get details of a specific price comparison
// @Tags         price-comparisons
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Price Comparison ID"
// @Success      200  {object}  PriceComparison
// @Failure      404  {object}  errors.AppError
// @Router       /price-comparisons/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/price-comparisons/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	comparison, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve price comparison", err.Error())
		return
	}
	utils.OKResponse(w, "Price comparison retrieved successfully", comparison)
}

// Create handles POST /api/v1/price-comparisons
// @Summary      Create price comparison
// @Description  Create a new price comparison
// @Tags         price-comparisons
// @Accept       json
// @Produce      json
// @Param        request  body      CreatePriceComparisonRequest  true  "Price comparison data"
// @Success      201      {object}  PriceComparison
// @Failure      400      {object}  errors.AppError
// @Router       /price-comparisons [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	var req CreatePriceComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	comparison, err := h.service.Create(&req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create price comparison", err.Error())
		return
	}
	utils.CreatedResponse(w, "Price comparison created successfully", comparison)
}

// Update handles PUT /api/v1/price-comparisons/:id
// @Summary      Update price comparison
// @Description  Update an existing price comparison
// @Tags         price-comparisons
// @Accept       json
// @Produce      json
// @Param        id       path      string                      true  "Price Comparison ID"
// @Param        request  body      UpdatePriceComparisonRequest true  "Price comparison data"
// @Success      200      {object}  PriceComparison
// @Failure      400      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /price-comparisons/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/price-comparisons/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdatePriceComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	comparison, err := h.service.Update(id, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update price comparison", err.Error())
		return
	}
	utils.OKResponse(w, "Price comparison updated successfully", comparison)
}
