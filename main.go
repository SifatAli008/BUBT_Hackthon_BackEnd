package main

import (
	_ "foodlink_backend/docs" // Import docs for Swagger
	"foodlink_backend/config"
	"foodlink_backend/database"
	"foodlink_backend/routes"
	"foodlink_backend/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title           Foodlink Backend API
// @version         1.0
// @description     A Go backend API server for Foodlink application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @schemes   http https

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize JWT
	utils.InitJWT(cfg)

	// Initialize database connection
	if cfg.DatabaseURL != "" {
		if err := database.Init(cfg); err != nil {
			log.Printf("Warning: Failed to initialize database: %v", err)
			log.Println("Server will start without database connection")
		} else {
			// Initialize schema if database is connected
			if err := database.InitSchema(); err != nil {
				log.Printf("Warning: Failed to initialize schema: %v", err)
			}
		}
		defer database.Close()
	} else {
		log.Println("Warning: DATABASE_URL not set, database features will be unavailable")
	}

	// Setup routes with middleware
	router := routes.SetupRoutes(cfg)

	// Start server
	serverAddr := ":" + cfg.Port
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Server running at http://localhost%s", serverAddr)
	log.Printf("Swagger UI: http://localhost%s/swagger/index.html", serverAddr)
	log.Printf("Swagger JSON: http://localhost%s/swagger/doc.json", serverAddr)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start: ", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down server...")

	// Close database connection
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server stopped")
}
