package auth

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for authentication
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(user *User) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		INSERT INTO users (id, email, name, password_hash, household_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, email, name, password_hash, household_id, role, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.HouseholdID,
		user.Role,
		now,
		now,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.HouseholdID,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.ErrAlreadyExists
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(email string) (*User, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	user := &User{}
	query := `
		SELECT id, email, name, password_hash, household_id, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.HouseholdID,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(id uuid.UUID) (*User, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	user := &User{}
	query := `
		SELECT id, email, name, password_hash, household_id, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.HouseholdID,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (r *Repository) UpdateUser(user *User) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		UPDATE users
		SET name = $1, household_id = $2, role = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, email, name, password_hash, household_id, role, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Name,
		user.HouseholdID,
		user.Role,
		time.Now(),
		user.ID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.HouseholdID,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}

	return nil
}

// EmailExists checks if an email already exists
func (r *Repository) EmailExists(email string) (bool, error) {
	if r.db == nil {
		return false, errors.ErrDatabase
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, errors.WrapError(err, errors.ErrDatabase)
	}

	return exists, nil
}
