package surplus

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

func (r *Repository) GetAll(status string) ([]*SurplusPost, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	var query string
	var rows *sql.Rows
	var err error
	if status != "" {
		query = `SELECT id, user_id, user_name, avatar_url, title, description, category, tags, quantity, unit, pickup_window, pickup_location, distance_km, image, status, expires_at, created_at, updated_at FROM community_surplus_posts WHERE status = $1 ORDER BY created_at DESC`
		rows, err = r.db.Query(query, status)
	} else {
		query = `SELECT id, user_id, user_name, avatar_url, title, description, category, tags, quantity, unit, pickup_window, pickup_location, distance_km, image, status, expires_at, created_at, updated_at FROM community_surplus_posts ORDER BY created_at DESC`
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var posts []*SurplusPost
	for rows.Next() {
		post := &SurplusPost{}
		var pickupWindowJSON []byte
		if err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.AvatarURL, &post.Title, &post.Description, &post.Category, pq.Array(&post.Tags), &post.Quantity, &post.Unit, &pickupWindowJSON, &post.PickupLocation, &post.DistanceKm, &post.Image, &post.Status, &post.ExpiresAt, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		if len(pickupWindowJSON) > 0 {
			json.Unmarshal(pickupWindowJSON, &post.PickupWindow)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*SurplusPost, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	post := &SurplusPost{}
	var pickupWindowJSON []byte
	query := `SELECT id, user_id, user_name, avatar_url, title, description, category, tags, quantity, unit, pickup_window, pickup_location, distance_km, image, status, expires_at, created_at, updated_at FROM community_surplus_posts WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.UserName, &post.AvatarURL, &post.Title, &post.Description, &post.Category, pq.Array(&post.Tags), &post.Quantity, &post.Unit, &pickupWindowJSON, &post.PickupLocation, &post.DistanceKm, &post.Image, &post.Status, &post.ExpiresAt, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(pickupWindowJSON) > 0 {
		json.Unmarshal(pickupWindowJSON, &post.PickupWindow)
	}
	return post, nil
}

func (r *Repository) Create(post *SurplusPost) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	pickupWindowJSON, _ := json.Marshal(post.PickupWindow)
	query := `INSERT INTO community_surplus_posts (id, user_id, user_name, avatar_url, title, description, category, tags, quantity, unit, pickup_window, pickup_location, distance_km, image, status, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id, user_id, user_name, avatar_url, title, description, category, tags, quantity, unit, pickup_window, pickup_location, distance_km, image, status, expires_at, created_at, updated_at`
	now := time.Now()
	var pickupWindowJSONOut []byte
	err := r.db.QueryRow(query, post.ID, post.UserID, post.UserName, post.AvatarURL, post.Title, post.Description, post.Category, pq.Array(post.Tags), post.Quantity, post.Unit, pickupWindowJSON, post.PickupLocation, post.DistanceKm, post.Image, post.Status, post.ExpiresAt, now, now).Scan(&post.ID, &post.UserID, &post.UserName, &post.AvatarURL, &post.Title, &post.Description, &post.Category, pq.Array(&post.Tags), &post.Quantity, &post.Unit, &pickupWindowJSONOut, &post.PickupLocation, &post.DistanceKm, &post.Image, &post.Status, &post.ExpiresAt, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(pickupWindowJSONOut) > 0 {
		json.Unmarshal(pickupWindowJSONOut, &post.PickupWindow)
	}
	return nil
}

func (r *Repository) Update(post *SurplusPost) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	pickupWindowJSON, _ := json.Marshal(post.PickupWindow)
	query := `UPDATE community_surplus_posts SET title=$1, description=$2, category=$3, tags=$4, quantity=$5, unit=$6, pickup_window=$7, pickup_location=$8, distance_km=$9, image=$10, status=$11, updated_at=$12 WHERE id=$13 RETURNING id, user_id, user_name, avatar_url, title, description, category, tags, quantity, unit, pickup_window, pickup_location, distance_km, image, status, expires_at, created_at, updated_at`
	var pickupWindowJSONOut []byte
	err := r.db.QueryRow(query, post.Title, post.Description, post.Category, pq.Array(post.Tags), post.Quantity, post.Unit, pickupWindowJSON, post.PickupLocation, post.DistanceKm, post.Image, post.Status, time.Now(), post.ID).Scan(&post.ID, &post.UserID, &post.UserName, &post.AvatarURL, &post.Title, &post.Description, &post.Category, pq.Array(&post.Tags), &post.Quantity, &post.Unit, &pickupWindowJSONOut, &post.PickupLocation, &post.DistanceKm, &post.Image, &post.Status, &post.ExpiresAt, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(pickupWindowJSONOut) > 0 {
		json.Unmarshal(pickupWindowJSONOut, &post.PickupWindow)
	}
	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	result, err := r.db.Exec(`DELETE FROM community_surplus_posts WHERE id = $1`, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *Repository) CreateRequest(req *SurplusRequest) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO surplus_requests (id, post_id, user_id, user_name, message, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, post_id, user_id, user_name, message, status, created_at`
	return r.db.QueryRow(query, req.ID, req.PostID, req.UserID, req.UserName, req.Message, req.Status, time.Now()).Scan(&req.ID, &req.PostID, &req.UserID, &req.UserName, &req.Message, &req.Status, &req.CreatedAt)
}

func (r *Repository) GetRequestsByPostID(postID uuid.UUID) ([]*SurplusRequest, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, post_id, user_id, user_name, message, status, created_at FROM surplus_requests WHERE post_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var requests []*SurplusRequest
	for rows.Next() {
		req := &SurplusRequest{}
		if err := rows.Scan(&req.ID, &req.PostID, &req.UserID, &req.UserName, &req.Message, &req.Status, &req.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *Repository) GetRequestByID(id uuid.UUID) (*SurplusRequest, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	req := &SurplusRequest{}
	query := `SELECT id, post_id, user_id, user_name, message, status, created_at FROM surplus_requests WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&req.ID, &req.PostID, &req.UserID, &req.UserName, &req.Message, &req.Status, &req.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return req, nil
}

func (r *Repository) UpdateRequest(req *SurplusRequest) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE surplus_requests SET status=$1 WHERE id=$2 RETURNING id, post_id, user_id, user_name, message, status, created_at`
	return r.db.QueryRow(query, req.Status, req.ID).Scan(&req.ID, &req.PostID, &req.UserID, &req.UserName, &req.Message, &req.Status, &req.CreatedAt)
}

func (r *Repository) CreateComment(comment *SurplusComment) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO surplus_comments (id, post_id, user_id, user_name, message, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, post_id, user_id, user_name, message, created_at`
	return r.db.QueryRow(query, comment.ID, comment.PostID, comment.UserID, comment.UserName, comment.Message, time.Now()).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.UserName, &comment.Message, &comment.CreatedAt)
}

func (r *Repository) GetCommentsByPostID(postID uuid.UUID) ([]*SurplusComment, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, post_id, user_id, user_name, message, created_at FROM surplus_comments WHERE post_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var comments []*SurplusComment
	for rows.Next() {
		comment := &SurplusComment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.UserName, &comment.Message, &comment.CreatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
