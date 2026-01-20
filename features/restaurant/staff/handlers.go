package staff

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

// GetAllTasks handles GET /api/v1/restaurant/tasks
// @Summary      List tasks
// @Description  Get all staff tasks for the authenticated restaurant
// @Tags         restaurant-staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   StaffTask
// @Failure      401  {object}  errors.AppError
// @Router       /restaurant/tasks [get]
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	tasks, err := h.service.GetAllTasks(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve tasks", err.Error())
		return
	}
	utils.OKResponse(w, "Tasks retrieved successfully", tasks)
}

// CreateTask handles POST /api/v1/restaurant/tasks
// @Summary      Create task
// @Description  Create a new staff task
// @Tags         restaurant-staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateStaffTaskRequest  true  "Task data"
// @Success      201      {object}  StaffTask
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /restaurant/tasks [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateStaffTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	task, err := h.service.CreateTask(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create task", err.Error())
		return
	}
	utils.CreatedResponse(w, "Task created successfully", task)
}

// UpdateTask handles PUT /api/v1/restaurant/tasks/:id
// @Summary      Update task
// @Description  Update an existing staff task
// @Tags         restaurant-staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "Task ID"
// @Param        request  body      UpdateStaffTaskRequest true  "Task data"
// @Success      200      {object}  StaffTask
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /restaurant/tasks/{id} [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/restaurant/tasks/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateStaffTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	task, err := h.service.UpdateTask(id, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update task", err.Error())
		return
	}
	utils.OKResponse(w, "Task updated successfully", task)
}

// GetAllShifts handles GET /api/v1/restaurant/shifts
// @Summary      List shifts
// @Description  Get all shift schedules for the authenticated restaurant
// @Tags         restaurant-staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   ShiftSchedule
// @Failure      401  {object}  errors.AppError
// @Router       /restaurant/shifts [get]
func (h *Handler) GetAllShifts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	shifts, err := h.service.GetAllShifts(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve shifts", err.Error())
		return
	}
	utils.OKResponse(w, "Shifts retrieved successfully", shifts)
}

// CreateShift handles POST /api/v1/restaurant/shifts
// @Summary      Create shift
// @Description  Create a new shift schedule
// @Tags         restaurant-staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateShiftScheduleRequest  true  "Shift data"
// @Success      201      {object}  ShiftSchedule
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /restaurant/shifts [post]
func (h *Handler) CreateShift(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateShiftScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	shift, err := h.service.CreateShift(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create shift", err.Error())
		return
	}
	utils.CreatedResponse(w, "Shift created successfully", shift)
}
