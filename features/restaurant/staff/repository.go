package staff

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetAllTasksByUserID(userID uuid.UUID) ([]*StaffTask, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, title, description, assignee, shift, completed, priority, created_at FROM restaurant_staff_tasks WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var tasks []*StaffTask
	for rows.Next() {
		task := &StaffTask{}
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Assignee, &task.Shift, &task.Completed, &task.Priority, &task.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) GetTaskByID(id uuid.UUID) (*StaffTask, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	task := &StaffTask{}
	query := `SELECT id, user_id, title, description, assignee, shift, completed, priority, created_at FROM restaurant_staff_tasks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Assignee, &task.Shift, &task.Completed, &task.Priority, &task.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return task, nil
}

func (r *Repository) CreateTask(task *StaffTask) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO restaurant_staff_tasks (id, user_id, title, description, assignee, shift, completed, priority, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, user_id, title, description, assignee, shift, completed, priority, created_at`
	return r.db.QueryRow(query, task.ID, task.UserID, task.Title, task.Description, task.Assignee, task.Shift, task.Completed, task.Priority, time.Now()).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Assignee, &task.Shift, &task.Completed, &task.Priority, &task.CreatedAt)
}

func (r *Repository) UpdateTask(task *StaffTask) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE restaurant_staff_tasks SET title=$1, description=$2, assignee=$3, shift=$4, completed=$5, priority=$6 WHERE id=$7 RETURNING id, user_id, title, description, assignee, shift, completed, priority, created_at`
	return r.db.QueryRow(query, task.Title, task.Description, task.Assignee, task.Shift, task.Completed, task.Priority, task.ID).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Assignee, &task.Shift, &task.Completed, &task.Priority, &task.CreatedAt)
}

func (r *Repository) GetAllShiftsByUserID(userID uuid.UUID) ([]*ShiftSchedule, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, role, staff, time, notes, created_at FROM restaurant_shift_schedule WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var shifts []*ShiftSchedule
	for rows.Next() {
		shift := &ShiftSchedule{}
		if err := rows.Scan(&shift.ID, &shift.UserID, &shift.Role, &shift.Staff, &shift.Time, &shift.Notes, &shift.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (r *Repository) CreateShift(shift *ShiftSchedule) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO restaurant_shift_schedule (id, user_id, role, staff, time, notes, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, role, staff, time, notes, created_at`
	return r.db.QueryRow(query, shift.ID, shift.UserID, shift.Role, shift.Staff, shift.Time, shift.Notes, time.Now()).Scan(&shift.ID, &shift.UserID, &shift.Role, &shift.Staff, &shift.Time, &shift.Notes, &shift.CreatedAt)
}
