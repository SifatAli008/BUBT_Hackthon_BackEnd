package price_comparisons

import (
	"database/sql"
	"encoding/json"
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

func (r *Repository) GetAll() ([]*PriceComparison, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, item_name, category, stores, best_price, updated_at FROM price_comparisons ORDER BY updated_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var comparisons []*PriceComparison
	for rows.Next() {
		c := &PriceComparison{}
		var storesJSON, bestPriceJSON []byte
		if err := rows.Scan(&c.ID, &c.ItemName, &c.Category, &storesJSON, &bestPriceJSON, &c.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		if len(storesJSON) > 0 {
			json.Unmarshal(storesJSON, &c.Stores)
		}
		if len(bestPriceJSON) > 0 {
			json.Unmarshal(bestPriceJSON, &c.BestPrice)
		}
		comparisons = append(comparisons, c)
	}
	return comparisons, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*PriceComparison, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	c := &PriceComparison{}
	var storesJSON, bestPriceJSON []byte
	query := `SELECT id, item_name, category, stores, best_price, updated_at FROM price_comparisons WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.ItemName, &c.Category, &storesJSON, &bestPriceJSON, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(storesJSON) > 0 {
		json.Unmarshal(storesJSON, &c.Stores)
	}
	if len(bestPriceJSON) > 0 {
		json.Unmarshal(bestPriceJSON, &c.BestPrice)
	}
	return c, nil
}

func (r *Repository) Create(c *PriceComparison) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	storesJSON, _ := json.Marshal(c.Stores)
	bestPriceJSON, _ := json.Marshal(c.BestPrice)
	query := `INSERT INTO price_comparisons (id, item_name, category, stores, best_price, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, item_name, category, stores, best_price, updated_at`
	now := time.Now()
	var storesJSONOut, bestPriceJSONOut []byte
	err := r.db.QueryRow(query, c.ID, c.ItemName, c.Category, storesJSON, bestPriceJSON, now).Scan(&c.ID, &c.ItemName, &c.Category, &storesJSONOut, &bestPriceJSONOut, &c.UpdatedAt)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(storesJSONOut) > 0 {
		json.Unmarshal(storesJSONOut, &c.Stores)
	}
	if len(bestPriceJSONOut) > 0 {
		json.Unmarshal(bestPriceJSONOut, &c.BestPrice)
	}
	return nil
}

func (r *Repository) Update(c *PriceComparison) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	storesJSON, _ := json.Marshal(c.Stores)
	bestPriceJSON, _ := json.Marshal(c.BestPrice)
	query := `UPDATE price_comparisons SET item_name=$1, category=$2, stores=$3, best_price=$4, updated_at=$5 WHERE id=$6 RETURNING id, item_name, category, stores, best_price, updated_at`
	var storesJSONOut, bestPriceJSONOut []byte
	err := r.db.QueryRow(query, c.ItemName, c.Category, storesJSON, bestPriceJSON, time.Now(), c.ID).Scan(&c.ID, &c.ItemName, &c.Category, &storesJSONOut, &bestPriceJSONOut, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(storesJSONOut) > 0 {
		json.Unmarshal(storesJSONOut, &c.Stores)
	}
	if len(bestPriceJSONOut) > 0 {
		json.Unmarshal(bestPriceJSONOut, &c.BestPrice)
	}
	return nil
}
