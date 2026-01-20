package auth

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"net/http"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Register a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration request"
// @Success      201      {object}  AuthResponse
// @Failure      400      {object}  errors.AppError
// @Failure      409      {object}  errors.AppError
// @Router       /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}

	response, err := h.service.Register(&req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to register user", err.Error())
		return
	}

	utils.CreatedResponse(w, "User registered successfully", response)
}

// Login handles user login
// @Summary      Login user
// @Description  Authenticate user and get access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login request"
// @Success      200      {object}  AuthResponse
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.UnauthorizedResponse(w, "Invalid credentials")
		return
	}

	utils.OKResponse(w, "Login successful", response)
}

// Logout handles user logout
// @Summary      Logout user
// @Description  Logout user (token invalidation handled client-side)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  errors.AppError
// @Router       /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	// In a stateless JWT system, logout is handled client-side
	// Server can maintain a blacklist if needed
	utils.OKResponse(w, "Logged out successfully", map[string]string{
		"message": "Logged out successfully",
	})
}

// RefreshToken handles token refresh
// @Summary      Refresh access token
// @Description  Refresh the access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  AuthResponse
// @Failure      401      {object}  errors.AppError
// @Router       /auth/refresh [post]
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		utils.UnauthorizedResponse(w, "Invalid token")
		return
	}

	// Generate new token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, h.service.jwtExpiry)
	if err != nil {
		utils.InternalServerErrorResponse(w, "Failed to generate token", err.Error())
		return
	}

	response := &AuthResponse{
		User:        user,
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.service.jwtExpiry.Seconds()),
	}

	utils.OKResponse(w, "Token refreshed successfully", response)
}

// GetMe handles getting current user info
// @Summary      Get current user
// @Description  Get authenticated user information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  UserResponse
// @Failure      401      {object}  errors.AppError
// @Router       /auth/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		utils.UnauthorizedResponse(w, "Invalid token")
		return
	}

	utils.OKResponse(w, "User retrieved successfully", user.ToUserResponse())
}
