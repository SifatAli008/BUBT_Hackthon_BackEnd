package consumption

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

func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]*ConsumptionLog, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, inventory_item_id, food_name, quantity, unit, category, consumed_at, was_wasted, notes, created_at, updated_at FROM consumption_logs WHERE user_id = $1 ORDER BY consumed_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var logs []*ConsumptionLog
	for rows.Next() {
		log := &ConsumptionLog{}
		if err := rows.Scan(&log.ID, &log.UserID, &log.InventoryItemID, &log.FoodName, &log.Quantity, &log.Unit, &log.Category, &log.ConsumedAt, &log.WasWasted, &log.Notes, &log.CreatedAt, &log.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*ConsumptionLog, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	log := &ConsumptionLog{}
	query := `SELECT id, user_id, inventory_item_id, food_name, quantity, unit, category, consumed_at, was_wasted, notes, created_at, updated_at FROM consumption_logs WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&log.ID, &log.UserID, &log.InventoryItemID, &log.FoodName, &log.Quantity, &log.Unit, &log.Category, &log.ConsumedAt, &log.WasWasted, &log.Notes, &log.CreatedAt, &log.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return log, nil
}

func (r *Repository) Create(log *ConsumptionLog) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO consumption_logs (id, user_id, inventory_item_id, food_name, quantity, unit, category, consumed_at, was_wasted, notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, user_id, inventory_item_id, food_name, quantity, unit, category, consumed_at, was_wasted, notes, created_at, updated_at`
	now := time.Now()
	return r.db.QueryRow(query, log.ID, log.UserID, log.InventoryItemID, log.FoodName, log.Quantity, log.Unit, log.Category, log.ConsumedAt, log.WasWasted, log.Notes, now, now).Scan(&log.ID, &log.UserID, &log.InventoryItemID, &log.FoodName, &log.Quantity, &log.Unit, &log.Category, &log.ConsumedAt, &log.WasWasted, &log.Notes, &log.CreatedAt, &log.UpdatedAt)
}

func (r *Repository) Update(log *ConsumptionLog) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE consumption_logs SET food_name=$1, quantity=$2, unit=$3, category=$4, consumed_at=$5, was_wasted=$6, notes=$7, updated_at=$8 WHERE id=$9 RETURNING id, user_id, inventory_item_id, food_name, quantity, unit, category, consumed_at, was_wasted, notes, created_at, updated_at`
	return r.db.QueryRow(query, log.FoodName, log.Quantity, log.Unit, log.Category, log.ConsumedAt, log.WasWasted, log.Notes, time.Now(), log.ID).Scan(&log.ID, &log.UserID, &log.InventoryItemID, &log.FoodName, &log.Quantity, &log.Unit, &log.Category, &log.ConsumedAt, &log.WasWasted, &log.Notes, &log.CreatedAt, &log.UpdatedAt)
}

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	result, err := r.db.Exec(`DELETE FROM consumption_logs WHERE id = $1`, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *Repository) GetStats(userID uuid.UUID) (*ConsumptionStats, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	stats := &ConsumptionStats{}
	query := `SELECT COALESCE(SUM(quantity), 0), COALESCE(SUM(CASE WHEN was_wasted THEN quantity ELSE 0 END), 0), COUNT(*) FROM consumption_logs WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&stats.TotalConsumed, &stats.TotalWasted, &stats.TotalLogs)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if stats.TotalConsumed > 0 {
		stats.WastePercentage = (stats.TotalWasted / stats.TotalConsumed) * 100
	}
	return stats, nil
}
