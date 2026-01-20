package inventory

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for inventory
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new inventory repository
func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

// GetAllByUserID retrieves all inventory items for a user
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]*InventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	query := `
		SELECT id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at
		FROM inventory_items
		WHERE user_id = $1
		ORDER BY expiry_date NULLS LAST, created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var items []*InventoryItem
	for rows.Next() {
		item := &InventoryItem{}
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Quantity,
			&item.Unit,
			&item.ExpiryDate,
			&item.Category,
			&item.Location,
			&item.FoodItemID,
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

// GetByID retrieves an inventory item by ID
func (r *Repository) GetByID(id uuid.UUID) (*InventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	item := &InventoryItem{}
	query := `
		SELECT id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at
		FROM inventory_items
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.UserID,
		&item.Name,
		&item.Quantity,
		&item.Unit,
		&item.ExpiryDate,
		&item.Category,
		&item.Location,
		&item.FoodItemID,
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

// Create creates a new inventory item
func (r *Repository) Create(item *InventoryItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		INSERT INTO inventory_items (id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		item.ID,
		item.UserID,
		item.Name,
		item.Quantity,
		item.Unit,
		item.ExpiryDate,
		item.Category,
		item.Location,
		item.FoodItemID,
		now,
		now,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.Name,
		&item.Quantity,
		&item.Unit,
		&item.ExpiryDate,
		&item.Category,
		&item.Location,
		&item.FoodItemID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}

	return nil
}

// Update updates an inventory item
func (r *Repository) Update(item *InventoryItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		UPDATE inventory_items
		SET name = $1, quantity = $2, unit = $3, expiry_date = $4, category = $5, location = $6, food_item_id = $7, updated_at = $8
		WHERE id = $9
		RETURNING id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		item.Name,
		item.Quantity,
		item.Unit,
		item.ExpiryDate,
		item.Category,
		item.Location,
		item.FoodItemID,
		time.Now(),
		item.ID,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.Name,
		&item.Quantity,
		&item.Unit,
		&item.ExpiryDate,
		&item.Category,
		&item.Location,
		&item.FoodItemID,
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

// Delete deletes an inventory item
func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `DELETE FROM inventory_items WHERE id = $1`
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

// GetExpiring retrieves items expiring within specified days
func (r *Repository) GetExpiring(userID uuid.UUID, days int) ([]*InventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	query := `
		SELECT id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at
		FROM inventory_items
		WHERE user_id = $1
		AND expiry_date IS NOT NULL
		AND expiry_date BETWEEN CURRENT_TIMESTAMP AND CURRENT_TIMESTAMP + INTERVAL '1 day' * $2
		ORDER BY expiry_date ASC
	`

	rows, err := r.db.Query(query, userID, days)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var items []*InventoryItem
	for rows.Next() {
		item := &InventoryItem{}
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Quantity,
			&item.Unit,
			&item.ExpiryDate,
			&item.Category,
			&item.Location,
			&item.FoodItemID,
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

// GetExpired retrieves expired items
func (r *Repository) GetExpired(userID uuid.UUID) ([]*InventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	query := `
		SELECT id, user_id, name, quantity, unit, expiry_date, category, location, food_item_id, created_at, updated_at
		FROM inventory_items
		WHERE user_id = $1
		AND expiry_date IS NOT NULL
		AND expiry_date < CURRENT_TIMESTAMP
		ORDER BY expiry_date ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var items []*InventoryItem
	for rows.Next() {
		item := &InventoryItem{}
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Quantity,
			&item.Unit,
			&item.ExpiryDate,
			&item.Category,
			&item.Location,
			&item.FoodItemID,
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
