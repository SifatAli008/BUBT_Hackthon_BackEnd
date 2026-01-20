package leftovers

import (
	"database/sql"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetAll(status string) ([]*LeftoverItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	var query string
	var rows *sql.Rows
	var err error
	if status != "" {
		query = `SELECT id, user_id, user_name, avatar_url, dish_name, description, portions, distance_km, dietary_tags, allergens, pickup_window, status, image, created_at, updated_at FROM leftover_items WHERE status = $1 ORDER BY created_at DESC`
		rows, err = r.db.Query(query, status)
	} else {
		query = `SELECT id, user_id, user_name, avatar_url, dish_name, description, portions, distance_km, dietary_tags, allergens, pickup_window, status, image, created_at, updated_at FROM leftover_items ORDER BY created_at DESC`
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var items []*LeftoverItem
	for rows.Next() {
		item := &LeftoverItem{}
		if err := rows.Scan(&item.ID, &item.UserID, &item.UserName, &item.AvatarURL, &item.DishName, &item.Description, &item.Portions, &item.DistanceKm, pq.Array(&item.DietaryTags), pq.Array(&item.Allergens), &item.PickupWindow, &item.Status, &item.Image, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*LeftoverItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	item := &LeftoverItem{}
	query := `SELECT id, user_id, user_name, avatar_url, dish_name, description, portions, distance_km, dietary_tags, allergens, pickup_window, status, image, created_at, updated_at FROM leftover_items WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.UserID, &item.UserName, &item.AvatarURL, &item.DishName, &item.Description, &item.Portions, &item.DistanceKm, pq.Array(&item.DietaryTags), pq.Array(&item.Allergens), &item.PickupWindow, &item.Status, &item.Image, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return item, nil
}

func (r *Repository) Create(item *LeftoverItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO leftover_items (id, user_id, user_name, avatar_url, dish_name, description, portions, distance_km, dietary_tags, allergens, pickup_window, status, image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id, user_id, user_name, avatar_url, dish_name, description, portions, distance_km, dietary_tags, allergens, pickup_window, status, image, created_at, updated_at`
	now := time.Now()
	return r.db.QueryRow(query, item.ID, item.UserID, item.UserName, item.AvatarURL, item.DishName, item.Description, item.Portions, item.DistanceKm, pq.Array(item.DietaryTags), pq.Array(item.Allergens), item.PickupWindow, item.Status, item.Image, now, now).Scan(&item.ID, &item.UserID, &item.UserName, &item.AvatarURL, &item.DishName, &item.Description, &item.Portions, &item.DistanceKm, pq.Array(&item.DietaryTags), pq.Array(&item.Allergens), &item.PickupWindow, &item.Status, &item.Image, &item.CreatedAt, &item.UpdatedAt)
}

func (r *Repository) Update(item *LeftoverItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE leftover_items SET dish_name=$1, description=$2, portions=$3, distance_km=$4, dietary_tags=$5, allergens=$6, pickup_window=$7, status=$8, image=$9, updated_at=$10 WHERE id=$11 RETURNING id, user_id, user_name, avatar_url, dish_name, description, portions, distance_km, dietary_tags, allergens, pickup_window, status, image, created_at, updated_at`
	return r.db.QueryRow(query, item.DishName, item.Description, item.Portions, item.DistanceKm, pq.Array(item.DietaryTags), pq.Array(item.Allergens), item.PickupWindow, item.Status, item.Image, time.Now(), item.ID).Scan(&item.ID, &item.UserID, &item.UserName, &item.AvatarURL, &item.DishName, &item.Description, &item.Portions, &item.DistanceKm, pq.Array(&item.DietaryTags), pq.Array(&item.Allergens), &item.PickupWindow, &item.Status, &item.Image, &item.CreatedAt, &item.UpdatedAt)
}

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	result, err := r.db.Exec(`DELETE FROM leftover_items WHERE id = $1`, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *Repository) CreateClaim(claim *LeftoverClaim) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO leftover_item_claims (id, leftover_item_id, user_id, user_name, message, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, leftover_item_id, user_id, user_name, message, created_at`
	return r.db.QueryRow(query, claim.ID, claim.LeftoverItemID, claim.UserID, claim.UserName, claim.Message, time.Now()).Scan(&claim.ID, &claim.LeftoverItemID, &claim.UserID, &claim.UserName, &claim.Message, &claim.CreatedAt)
}

func (r *Repository) GetClaimsByLeftoverID(leftoverID uuid.UUID) ([]*LeftoverClaim, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, leftover_item_id, user_id, user_name, message, created_at FROM leftover_item_claims WHERE leftover_item_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, leftoverID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var claims []*LeftoverClaim
	for rows.Next() {
		claim := &LeftoverClaim{}
		if err := rows.Scan(&claim.ID, &claim.LeftoverItemID, &claim.UserID, &claim.UserName, &claim.Message, &claim.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		claims = append(claims, claim)
	}
	return claims, nil
}
