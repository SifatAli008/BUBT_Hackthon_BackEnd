package kitchen_events

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

// GetAll handles GET /api/v1/community/kitchen-events
// @Summary      List kitchen events
// @Description  Get all community kitchen events, optionally filtered by status
// @Tags         community-kitchen-events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status  query     string  false  "Filter by status (upcoming, in-progress, completed)"
// @Success      200     {array}   KitchenEvent
// @Failure      401     {object}  errors.AppError
// @Router       /community/kitchen-events [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	status := r.URL.Query().Get("status")
	events, err := h.service.GetAll(status)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve kitchen events", err.Error())
		return
	}
	utils.OKResponse(w, "Kitchen events retrieved successfully", events)
}

// GetByID handles GET /api/v1/community/kitchen-events/:id
// @Summary      Get kitchen event by ID
// @Description  Get details of a specific kitchen event
// @Tags         community-kitchen-events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Kitchen Event ID"
// @Success      200  {object}  KitchenEvent
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /community/kitchen-events/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/kitchen-events/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	event, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve kitchen event", err.Error())
		return
	}
	utils.OKResponse(w, "Kitchen event retrieved successfully", event)
}

// Create handles POST /api/v1/community/kitchen-events
// @Summary      Create kitchen event
// @Description  Create a new community kitchen event
// @Tags         community-kitchen-events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateKitchenEventRequest  true  "Kitchen event data"
// @Success      201      {object}  KitchenEvent
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/kitchen-events [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	var req CreateKitchenEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	event, err := h.service.Create(&req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create kitchen event", err.Error())
		return
	}
	utils.CreatedResponse(w, "Kitchen event created successfully", event)
}

// Update handles PUT /api/v1/community/kitchen-events/:id
// @Summary      Update kitchen event
// @Description  Update an existing kitchen event
// @Tags         community-kitchen-events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Kitchen Event ID"
// @Param        request  body      UpdateKitchenEventRequest  true  "Kitchen event data"
// @Success      200      {object}  KitchenEvent
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /community/kitchen-events/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/kitchen-events/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateKitchenEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	event, err := h.service.Update(id, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update kitchen event", err.Error())
		return
	}
	utils.OKResponse(w, "Kitchen event updated successfully", event)
}

// Volunteer handles POST /api/v1/community/kitchen-events/:id/volunteer
// @Summary      Volunteer for event
// @Description  Sign up as a volunteer for a kitchen event
// @Tags         community-kitchen-events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string            true  "Kitchen Event ID"
// @Param        request  body      VolunteerRequest  true  "Volunteer data"
// @Success      200      {object}  KitchenEvent
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/kitchen-events/{id}/volunteer [post]
func (h *Handler) Volunteer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, userName, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/kitchen-events/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "volunteer" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	eventID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid event ID", nil)
		return
	}
	var req VolunteerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	event, err := h.service.AddVolunteer(eventID, userID, userName, req.Role)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to volunteer for event", err.Error())
		return
	}
	utils.OKResponse(w, "Volunteered successfully", event)
}
