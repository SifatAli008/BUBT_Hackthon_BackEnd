package food_items

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for food items
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new food items repository
func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

// GetAll retrieves all food items
func (r *Repository) GetAll() ([]*FoodItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	query := `
		SELECT id, name, category, typical_expiry_days, storage_tips, created_at, updated_at
		FROM food_items
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var items []*FoodItem
	for rows.Next() {
		item := &FoodItem{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Category,
			&item.TypicalExpiryDays,
			&item.StorageTips,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetByID retrieves a food item by ID
func (r *Repository) GetByID(id uuid.UUID) (*FoodItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	item := &FoodItem{}
	query := `
		SELECT id, name, category, typical_expiry_days, storage_tips, created_at, updated_at
		FROM food_items
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Category,
		&item.TypicalExpiryDays,
		&item.StorageTips,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}

	return item, nil
}

// Create creates a new food item
func (r *Repository) Create(item *FoodItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		INSERT INTO food_items (id, name, category, typical_expiry_days, storage_tips, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, category, typical_expiry_days, storage_tips, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		item.ID,
		item.Name,
		item.Category,
		item.TypicalExpiryDays,
		item.StorageTips,
		now,
		now,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Category,
		&item.TypicalExpiryDays,
		&item.StorageTips,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}

	return nil
}

// Update updates a food item
func (r *Repository) Update(item *FoodItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		UPDATE food_items
		SET name = $1, category = $2, typical_expiry_days = $3, storage_tips = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, name, category, typical_expiry_days, storage_tips, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		item.Name,
		item.Category,
		item.TypicalExpiryDays,
		item.StorageTips,
		time.Now(),
		item.ID,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Category,
		&item.TypicalExpiryDays,
		&item.StorageTips,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}

	return nil
}

// Delete deletes a food item
func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `DELETE FROM food_items WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}
