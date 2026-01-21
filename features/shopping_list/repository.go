package shopping_list

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetAllByUserID(userID uuid.UUID, includePurchased bool) ([]*ShoppingListItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	query := `
		SELECT id, user_id, household_id, name, quantity, unit, category, priority, purchased, purchased_at, estimated_price, created_at, updated_at
		FROM shopping_list_items
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	if !includePurchased {
		query += " AND purchased = FALSE"
	}
	query += " ORDER BY purchased ASC, priority DESC, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var items []*ShoppingListItem
	for rows.Next() {
		item := &ShoppingListItem{}
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.HouseholdID,
			&item.Name,
			&item.Quantity,
			&item.Unit,
			&item.Category,
			&item.Priority,
			&item.Purchased,
			&item.PurchasedAt,
			&item.EstimatedPrice,
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

func (r *Repository) GetByID(id uuid.UUID) (*ShoppingListItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	item := &ShoppingListItem{}
	query := `
		SELECT id, user_id, household_id, name, quantity, unit, category, priority, purchased, purchased_at, estimated_price, created_at, updated_at
		FROM shopping_list_items
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.UserID,
		&item.HouseholdID,
		&item.Name,
		&item.Quantity,
		&item.Unit,
		&item.Category,
		&item.Priority,
		&item.Purchased,
		&item.PurchasedAt,
		&item.EstimatedPrice,
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

func (r *Repository) Create(item *ShoppingListItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	now := time.Now()
	query := `
		INSERT INTO shopping_list_items (id, user_id, household_id, name, quantity, unit, category, priority, purchased, purchased_at, estimated_price, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, user_id, household_id, name, quantity, unit, category, priority, purchased, purchased_at, estimated_price, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		item.ID,
		item.UserID,
		item.HouseholdID,
		item.Name,
		item.Quantity,
		item.Unit,
		item.Category,
		item.Priority,
		item.Purchased,
		item.PurchasedAt,
		item.EstimatedPrice,
		now,
		now,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.HouseholdID,
		&item.Name,
		&item.Quantity,
		&item.Unit,
		&item.Category,
		&item.Priority,
		&item.Purchased,
		&item.PurchasedAt,
		&item.EstimatedPrice,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	return nil
}

func (r *Repository) Update(item *ShoppingListItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}

	query := `
		UPDATE shopping_list_items
		SET name=$1, quantity=$2, unit=$3, category=$4, priority=$5, purchased=$6, purchased_at=$7, estimated_price=$8, updated_at=$9
		WHERE id=$10
		RETURNING id, user_id, household_id, name, quantity, unit, category, priority, purchased, purchased_at, estimated_price, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		item.Name,
		item.Quantity,
		item.Unit,
		item.Category,
		item.Priority,
		item.Purchased,
		item.PurchasedAt,
		item.EstimatedPrice,
		time.Now(),
		item.ID,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.HouseholdID,
		&item.Name,
		&item.Quantity,
		&item.Unit,
		&item.Category,
		&item.Priority,
		&item.Purchased,
		&item.PurchasedAt,
		&item.EstimatedPrice,
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

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	res, err := r.db.Exec(`DELETE FROM shopping_list_items WHERE id=$1`, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if ra == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *Repository) ListUnpurchasedNamesLower(userID uuid.UUID) (map[string]struct{}, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	rows, err := r.db.Query(`SELECT name FROM shopping_list_items WHERE user_id=$1 AND purchased=FALSE`, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	set := map[string]struct{}{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		set[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}
	return set, nil
}

func (r *Repository) GetLowStockInventoryNames(userID uuid.UUID) ([]struct {
	Name     string
	Unit     sql.NullString
	Category sql.NullString
}, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}

	rows, err := r.db.Query(`SELECT name, unit, category FROM inventory_items WHERE user_id=$1 AND quantity < 1`, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var out []struct {
		Name     string
		Unit     sql.NullString
		Category sql.NullString
	}
	for rows.Next() {
		var r0 struct {
			Name     string
			Unit     sql.NullString
			Category sql.NullString
		}
		if err := rows.Scan(&r0.Name, &r0.Unit, &r0.Category); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		out = append(out, r0)
	}
	return out, nil
}

