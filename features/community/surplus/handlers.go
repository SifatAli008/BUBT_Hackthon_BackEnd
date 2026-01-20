package surplus

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) getUserID(r *http.Request) (uuid.UUID, string, string, error) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		return uuid.Nil, "", "", errors.ErrUnauthorized
	}
	return user.ID, user.Name, "", nil // AvatarURL will be handled via community profiles later
}

// GetAll handles GET /api/v1/community/surplus
// @Summary      List surplus posts
// @Description  Get all community surplus posts, optionally filtered by status
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status  query     string  false  "Filter by status (available, claimed, expired)"
// @Success      200     {array}   SurplusPost
// @Failure      401     {object}  errors.AppError
// @Router       /community/surplus [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	status := r.URL.Query().Get("status")
	posts, err := h.service.GetAll(status)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve surplus posts", err.Error())
		return
	}
	utils.OKResponse(w, "Surplus posts retrieved successfully", posts)
}

// GetByID handles GET /api/v1/community/surplus/:id
// @Summary      Get surplus post by ID
// @Description  Get details of a specific surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Surplus Post ID"
// @Success      200  {object}  SurplusPost
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /community/surplus/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	post, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve surplus post", err.Error())
		return
	}
	utils.OKResponse(w, "Surplus post retrieved successfully", post)
}

// Create handles POST /api/v1/community/surplus
// @Summary      Create surplus post
// @Description  Create a new community surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateSurplusPostRequest  true  "Surplus post data"
// @Success      201      {object}  SurplusPost
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/surplus [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, userName, avatarURL, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateSurplusPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	post, err := h.service.Create(userID, userName, avatarURL, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create surplus post", err.Error())
		return
	}
	utils.CreatedResponse(w, "Surplus post created successfully", post)
}

// Update handles PUT /api/v1/community/surplus/:id
// @Summary      Update surplus post
// @Description  Update an existing surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Surplus Post ID"
// @Param        request  body      UpdateSurplusPostRequest  true  "Surplus post data"
// @Success      200      {object}  SurplusPost
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /community/surplus/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateSurplusPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	post, err := h.service.Update(id, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update surplus post", err.Error())
		return
	}
	utils.OKResponse(w, "Surplus post updated successfully", post)
}

// Delete handles DELETE /api/v1/community/surplus/:id
// @Summary      Delete surplus post
// @Description  Delete a surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Surplus Post ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.AppError
// @Failure      403  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /community/surplus/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to delete surplus post", err.Error())
		return
	}
	utils.OKResponse(w, "Surplus post deleted successfully", map[string]string{"message": "Deleted"})
}

// CreateRequest handles POST /api/v1/community/surplus/:id/request
// @Summary      Request surplus
// @Description  Create a request for a surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                        true  "Surplus Post ID"
// @Param        request  body      CreateSurplusRequestRequest   true  "Request data"
// @Success      201      {object}  SurplusRequest
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/surplus/{id}/request [post]
func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, userName, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "request" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	postID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid post ID", nil)
		return
	}
	var req CreateSurplusRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	request, err := h.service.CreateRequest(postID, userID, userName, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create request", err.Error())
		return
	}
	utils.CreatedResponse(w, "Request created successfully", request)
}

// GetRequests handles GET /api/v1/community/surplus/:id/requests
// @Summary      Get surplus requests
// @Description  Get all requests for a surplus post (post owner only)
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Surplus Post ID"
// @Success      200  {array}   SurplusRequest
// @Failure      401  {object}  errors.AppError
// @Failure      403  {object}  errors.AppError
// @Router       /community/surplus/{id}/requests [get]
func (h *Handler) GetRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "requests" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	postID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid post ID", nil)
		return
	}
	requests, err := h.service.GetRequestsByPostID(postID, userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve requests", err.Error())
		return
	}
	utils.OKResponse(w, "Requests retrieved successfully", requests)
}

// UpdateRequest handles PUT /api/v1/community/surplus/:id/requests/:requestId
// @Summary      Approve/decline request
// @Description  Approve or decline a surplus request (post owner only)
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      string                        true  "Surplus Post ID"
// @Param        requestId  path      string                        true  "Request ID"
// @Param        request    body      UpdateSurplusRequestRequest   true  "Request status"
// @Success      200        {object}  SurplusRequest
// @Failure      400        {object}  errors.AppError
// @Failure      401        {object}  errors.AppError
// @Failure      403        {object}  errors.AppError
// @Router       /community/surplus/{id}/requests/{requestId} [put]
func (h *Handler) UpdateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, _, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/"), "/")
	if len(pathParts) < 3 || pathParts[1] != "requests" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	postID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid post ID", nil)
		return
	}
	requestID, err := uuid.Parse(pathParts[2])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid request ID", nil)
		return
	}
	var req UpdateSurplusRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	request, err := h.service.UpdateRequest(requestID, postID, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update request", err.Error())
		return
	}
	utils.OKResponse(w, "Request updated successfully", request)
}

// CreateComment handles POST /api/v1/community/surplus/:id/comments
// @Summary      Add comment
// @Description  Add a comment to a surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                      true  "Surplus Post ID"
// @Param        request  body      CreateSurplusCommentRequest true  "Comment data"
// @Success      201      {object}  SurplusComment
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /community/surplus/{id}/comments [post]
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, userName, _, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "comments" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	postID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid post ID", nil)
		return
	}
	var req CreateSurplusCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	comment, err := h.service.CreateComment(postID, userID, userName, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create comment", err.Error())
		return
	}
	utils.CreatedResponse(w, "Comment created successfully", comment)
}

// GetComments handles GET /api/v1/community/surplus/:id/comments
// @Summary      Get comments
// @Description  Get all comments for a surplus post
// @Tags         community-surplus
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Surplus Post ID"
// @Success      200  {array}   SurplusComment
// @Failure      401  {object}  errors.AppError
// @Router       /community/surplus/{id}/comments [get]
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "comments" {
		utils.BadRequestResponse(w, "Invalid path", nil)
		return
	}
	postID, err := uuid.Parse(pathParts[0])
	if err != nil {
		utils.BadRequestResponse(w, "Invalid post ID", nil)
		return
	}
	comments, err := h.service.GetCommentsByPostID(postID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve comments", err.Error())
		return
	}
	utils.OKResponse(w, "Comments retrieved successfully", comments)
}
