package meal_plans

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

func (r *Repository) GetByID(id uuid.UUID) (*MealPlan, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	mp := &MealPlan{}
	var ingredients pq.StringArray

	query := `
		SELECT id, user_id, household_id, date, meal_type, name, description, ingredients, servings, created_at, updated_at
		FROM meal_plans
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&mp.ID,
		&mp.UserID,
		&mp.HouseholdID,
		&mp.Date,
		&mp.MealType,
		&mp.Name,
		&mp.Description,
		pq.Array(&ingredients),
		&mp.Servings,
		&mp.CreatedAt,
		&mp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	mp.Ingredients = []string(ingredients)
	return mp, nil
}

func (r *Repository) GetRangeByUserID(userID uuid.UUID, start time.Time, end time.Time) ([]*MealPlan, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `
		SELECT id, user_id, household_id, date, meal_type, name, description, ingredients, servings, created_at, updated_at
		FROM meal_plans
		WHERE user_id = $1
		  AND date >= $2
		  AND date < $3
		ORDER BY date ASC, meal_type ASC
	`
	rows, err := r.db.Query(query, userID, start, end)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()

	var out []*MealPlan
	for rows.Next() {
		mp := &MealPlan{}
		var ingredients pq.StringArray
		if err := rows.Scan(
			&mp.ID,
			&mp.UserID,
			&mp.HouseholdID,
			&mp.Date,
			&mp.MealType,
			&mp.Name,
			&mp.Description,
			pq.Array(&ingredients),
			&mp.Servings,
			&mp.CreatedAt,
			&mp.UpdatedAt,
		); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		mp.Ingredients = []string(ingredients)
		out = append(out, mp)
	}
	return out, nil
}

func (r *Repository) FindByUserDateType(userID uuid.UUID, date time.Time, mealType string) (*uuid.UUID, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	var id uuid.UUID
	err := r.db.QueryRow(`SELECT id FROM meal_plans WHERE user_id=$1 AND date=$2 AND meal_type=$3 LIMIT 1`, userID, date, mealType).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return &id, nil
}

func (r *Repository) Create(mp *MealPlan) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	now := time.Now()
	var ingredients pq.StringArray = mp.Ingredients

	query := `
		INSERT INTO meal_plans (id, user_id, household_id, date, meal_type, name, description, ingredients, servings, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, user_id, household_id, date, meal_type, name, description, ingredients, servings, created_at, updated_at
	`
	var outIngredients pq.StringArray
	err := r.db.QueryRow(
		query,
		mp.ID,
		mp.UserID,
		mp.HouseholdID,
		mp.Date,
		mp.MealType,
		mp.Name,
		mp.Description,
		pq.Array(&ingredients),
		mp.Servings,
		now,
		now,
	).Scan(
		&mp.ID,
		&mp.UserID,
		&mp.HouseholdID,
		&mp.Date,
		&mp.MealType,
		&mp.Name,
		&mp.Description,
		pq.Array(&outIngredients),
		&mp.Servings,
		&mp.CreatedAt,
		&mp.UpdatedAt,
	)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	mp.Ingredients = []string(outIngredients)
	return nil
}

func (r *Repository) Update(mp *MealPlan) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	var ingredients pq.StringArray = mp.Ingredients

	query := `
		UPDATE meal_plans
		SET name=$1, description=$2, ingredients=$3, servings=$4, updated_at=$5
		WHERE id=$6
		RETURNING id, user_id, household_id, date, meal_type, name, description, ingredients, servings, created_at, updated_at
	`
	var outIngredients pq.StringArray
	err := r.db.QueryRow(
		query,
		mp.Name,
		mp.Description,
		pq.Array(&ingredients),
		mp.Servings,
		time.Now(),
		mp.ID,
	).Scan(
		&mp.ID,
		&mp.UserID,
		&mp.HouseholdID,
		&mp.Date,
		&mp.MealType,
		&mp.Name,
		&mp.Description,
		pq.Array(&outIngredients),
		&mp.Servings,
		&mp.CreatedAt,
		&mp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	mp.Ingredients = []string(outIngredients)
	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	res, err := r.db.Exec(`DELETE FROM meal_plans WHERE id=$1`, id)
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

