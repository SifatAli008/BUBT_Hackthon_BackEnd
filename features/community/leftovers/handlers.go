package leftovers

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

func (h *Handler) getUserID(r *http.Request) (uuid.UUID, string, error) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		return uuid.Nil, "", errors.ErrUnauthorized
	}
	return user.ID, user.Name, nil
}

// GetAll handles GET /api/v1/community/leftovers
// @Summary      List leftover items
// @Description  Get all leftover items, optionally filtered by status
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status  query     string  false  "Filter by status (available, claimed)"
// @Success      200     {array}   LeftoverItem
// @Failure      401     {object}  errors.AppError
// @Router       /community/leftovers [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	status := r.URL.Query().Get("status")
	items, err := h.service.GetAll(status)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve leftover items", err.Error())
		return
	}
	utils.OKResponse(w, "Leftover items retrieved successfully", items)
}

// GetByID handles GET /api/v1/community/leftovers/:id
// @Summary      Get leftover item by ID
// @Description  Get details of a specific leftover item
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Leftover Item ID"
// @Success      200  {object}  LeftoverItem
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /community/leftovers/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/leftovers/")
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
		utils.InternalServerErrorResponse(w, "Failed to retrieve leftover item", err.Error())
		return
	}
	utils.OKResponse(w, "Leftover item retrieved successfully", item)
}

// Create handles POST /api/v1/community/leftovers
// @Summary      Create leftover item
// @Description  Create a new leftover item post
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateLeftoverItemRequest  true  "Leftover item data"
// @Success      201      {object}  LeftoverItem
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/leftovers [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, userName, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateLeftoverItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	item, err := h.service.Create(userID, userName, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create leftover item", err.Error())
		return
	}
	utils.CreatedResponse(w, "Leftover item created successfully", item)
}

// Update handles PUT /api/v1/community/leftovers/:id
// @Summary      Update leftover item
// @Description  Update an existing leftover item
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Leftover Item ID"
// @Param        request  body      UpdateLeftoverItemRequest true  "Leftover item data"
// @Success      200      {object}  LeftoverItem
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /community/leftovers/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/leftovers/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateLeftoverItemRequest
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
		utils.InternalServerErrorResponse(w, "Failed to update leftover item", err.Error())
		return
	}
	utils.OKResponse(w, "Leftover item updated successfully", item)
}

// Delete handles DELETE /api/v1/community/leftovers/:id
// @Summary      Delete leftover item
// @Description  Delete a leftover item
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Leftover Item ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.AppError
// @Failure      403  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /community/leftovers/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/leftovers/")
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
		utils.InternalServerErrorResponse(w, "Failed to delete leftover item", err.Error())
		return
	}
	utils.OKResponse(w, "Leftover item deleted successfully", map[string]string{"message": "Deleted"})
}

// CreateClaim handles POST /api/v1/community/leftovers/:id/claim
// @Summary      Claim leftover
// @Description  Claim a leftover item
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Leftover Item ID"
// @Param        request  body      CreateLeftoverClaimRequest true  "Claim data"
// @Success      201      {object}  LeftoverClaim
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/leftovers/{id}/claim [post]
func (h *Handler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, userName, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/leftovers/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "claim" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	leftoverID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid leftover ID", nil)
		return
	}
	var req CreateLeftoverClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	claim, err := h.service.CreateClaim(leftoverID, userID, userName, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create claim", err.Error())
		return
	}
	utils.CreatedResponse(w, "Claim created successfully", claim)
}

// GetClaims handles GET /api/v1/community/leftovers/:id/claims
// @Summary      Get claims
// @Description  Get all claims for a leftover item
// @Tags         community-leftovers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Leftover Item ID"
// @Success      200  {array}   LeftoverClaim
// @Failure      401  {object}  errors.AppError
// @Router       /community/leftovers/{id}/claims [get]
func (h *Handler) GetClaims(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/leftovers/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "claims" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	leftoverID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid leftover ID", nil)
		return
	}
	claims, err := h.service.GetClaimsByLeftoverID(leftoverID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve claims", err.Error())
		return
	}
	utils.OKResponse(w, "Claims retrieved successfully", claims)
}
