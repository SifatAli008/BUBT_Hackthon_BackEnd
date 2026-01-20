package staff

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) GetAllTasks(userID uuid.UUID) ([]*StaffTask, error) {
	return s.repo.GetAllTasksByUserID(userID)
}

func (s *Service) GetTaskByID(id uuid.UUID) (*StaffTask, error) {
	return s.repo.GetTaskByID(id)
}

func (s *Service) CreateTask(userID uuid.UUID, req *CreateStaffTaskRequest) (*StaffTask, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}
	task := &StaffTask{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Assignee:    req.Assignee,
		Shift:       req.Shift,
		Completed:   false,
		Priority:    priority,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) UpdateTask(id uuid.UUID, userID uuid.UUID, req *UpdateStaffTaskRequest) (*StaffTask, error) {
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	if task.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Assignee != "" {
		task.Assignee = req.Assignee
	}
	if req.Shift != "" {
		task.Shift = req.Shift
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if err := s.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) GetAllShifts(userID uuid.UUID) ([]*ShiftSchedule, error) {
	return s.repo.GetAllShiftsByUserID(userID)
}

func (s *Service) CreateShift(userID uuid.UUID, req *CreateShiftScheduleRequest) (*ShiftSchedule, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	shift := &ShiftSchedule{
		ID:     uuid.New(),
		UserID: userID,
		Role:   req.Role,
		Staff:  req.Staff,
		Time:   req.Time,
		Notes:  req.Notes,
	}
	if err := s.repo.CreateShift(shift); err != nil {
		return nil, err
	}
	return shift, nil
}
