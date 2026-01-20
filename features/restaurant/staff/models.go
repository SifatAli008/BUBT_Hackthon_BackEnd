package staff

import (
	"time"

	"github.com/google/uuid"
)

type StaffTask struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Assignee    string    `json:"assignee" db:"assignee"`
	Shift       string    `json:"shift" db:"shift"`
	Completed   bool      `json:"completed" db:"completed"`
	Priority    string    `json:"priority" db:"priority"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ShiftSchedule struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"`
	Staff     string    `json:"staff" db:"staff"`
	Time      string    `json:"time" db:"time"`
	Notes     string    `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateStaffTaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty"`
	Assignee    string `json:"assignee" validate:"required,min=1"`
	Shift       string `json:"shift" validate:"required,min=1"`
	Priority    string `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
}

type UpdateStaffTaskRequest struct {
	Title       string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description string `json:"description,omitempty"`
	Assignee    string `json:"assignee,omitempty" validate:"omitempty,min=1"`
	Shift       string `json:"shift,omitempty" validate:"omitempty,min=1"`
	Completed   *bool  `json:"completed,omitempty"`
	Priority    string `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
}

type CreateShiftScheduleRequest struct {
	Role  string `json:"role" validate:"required,min=1"`
	Staff string `json:"staff" validate:"required,min=1"`
	Time  string `json:"time" validate:"required,min=1"`
	Notes string `json:"notes,omitempty"`
}
