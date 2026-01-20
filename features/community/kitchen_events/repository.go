package kitchen_events

import (
	"database/sql"
	"encoding/json"
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

func (r *Repository) GetAll(status string) ([]*KitchenEvent, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	var query string
	var rows *sql.Rows
	var err error
	if status != "" {
		query = `SELECT id, title, description, date, time, location, tags, volunteers_needed, volunteers, food_saved_kg, status, image, created_at, updated_at FROM community_kitchen_events WHERE status = $1 ORDER BY date ASC, time ASC`
		rows, err = r.db.Query(query, status)
	} else {
		query = `SELECT id, title, description, date, time, location, tags, volunteers_needed, volunteers, food_saved_kg, status, image, created_at, updated_at FROM community_kitchen_events ORDER BY date ASC, time ASC`
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var events []*KitchenEvent
	for rows.Next() {
		event := &KitchenEvent{}
		var volunteersJSON []byte
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Time, &event.Location, pq.Array(&event.Tags), &event.VolunteersNeeded, &volunteersJSON, &event.FoodSavedKg, &event.Status, &event.Image, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		if len(volunteersJSON) > 0 {
			json.Unmarshal(volunteersJSON, &event.Volunteers)
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*KitchenEvent, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	event := &KitchenEvent{}
	var volunteersJSON []byte
	query := `SELECT id, title, description, date, time, location, tags, volunteers_needed, volunteers, food_saved_kg, status, image, created_at, updated_at FROM community_kitchen_events WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Time, &event.Location, pq.Array(&event.Tags), &event.VolunteersNeeded, &volunteersJSON, &event.FoodSavedKg, &event.Status, &event.Image, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(volunteersJSON) > 0 {
		json.Unmarshal(volunteersJSON, &event.Volunteers)
	}
	return event, nil
}

func (r *Repository) Create(event *KitchenEvent) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	volunteersJSON, _ := json.Marshal(event.Volunteers)
	query := `INSERT INTO community_kitchen_events (id, title, description, date, time, location, tags, volunteers_needed, volunteers, food_saved_kg, status, image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, title, description, date, time, location, tags, volunteers_needed, volunteers, food_saved_kg, status, image, created_at, updated_at`
	now := time.Now()
	var volunteersJSONOut []byte
	err := r.db.QueryRow(query, event.ID, event.Title, event.Description, event.Date, event.Time, event.Location, pq.Array(event.Tags), event.VolunteersNeeded, volunteersJSON, event.FoodSavedKg, event.Status, event.Image, now, now).Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Time, &event.Location, pq.Array(&event.Tags), &event.VolunteersNeeded, &volunteersJSONOut, &event.FoodSavedKg, &event.Status, &event.Image, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(volunteersJSONOut) > 0 {
		json.Unmarshal(volunteersJSONOut, &event.Volunteers)
	}
	return nil
}

func (r *Repository) Update(event *KitchenEvent) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	volunteersJSON, _ := json.Marshal(event.Volunteers)
	query := `UPDATE community_kitchen_events SET title=$1, description=$2, date=$3, time=$4, location=$5, tags=$6, volunteers_needed=$7, volunteers=$8, food_saved_kg=$9, status=$10, image=$11, updated_at=$12 WHERE id=$13 RETURNING id, title, description, date, time, location, tags, volunteers_needed, volunteers, food_saved_kg, status, image, created_at, updated_at`
	var volunteersJSONOut []byte
	err := r.db.QueryRow(query, event.Title, event.Description, event.Date, event.Time, event.Location, pq.Array(event.Tags), event.VolunteersNeeded, volunteersJSON, event.FoodSavedKg, event.Status, event.Image, time.Now(), event.ID).Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Time, &event.Location, pq.Array(&event.Tags), &event.VolunteersNeeded, &volunteersJSONOut, &event.FoodSavedKg, &event.Status, &event.Image, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(volunteersJSONOut) > 0 {
		json.Unmarshal(volunteersJSONOut, &event.Volunteers)
	}
	return nil
}
