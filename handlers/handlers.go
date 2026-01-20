package handlers

import (
	"foodlink_backend/database"
	"foodlink_backend/utils"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Server is running"`
	Database string `json:"database,omitempty" example:"connected"`
}

type APIResponse struct {
	Message string `json:"message" example:"Welcome to Foodlink API v1"`
}

// HealthCheck handles the health check endpoint
// @Summary      Health check
// @Description  Check if the server is running and database connectivity
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Message: "Server is running",
	}

	// Check database connection if available
	if database.GetDB() != nil {
		if err := database.HealthCheck(); err == nil {
			response.Database = "connected"
		} else {
			response.Database = "disconnected"
		}
	}

	utils.OKResponse(w, "Server is healthy", response)
}

// APIV1 handles API v1 routes
// @Summary      API v1 welcome
// @Description  Welcome message for API v1
// @Tags         api
// @Accept       json
// @Produce      json
// @Success      200  {object}  APIResponse
// @Router       /api/v1/ [get]
func APIV1(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Message: "Welcome to Foodlink API v1",
	}

	utils.OKResponse(w, "API v1", response)
}
