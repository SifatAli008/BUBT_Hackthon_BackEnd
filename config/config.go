package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	DatabaseURL string

	// JWT
	JWTSecret        string
	JWTExpiry        string
	JWTRefreshSecret string
	JWTRefreshExpiry string

	// DB pool (optional)
	DBMaxPoolSize      int
	DBMaxIdleConns     int
	DBIdleTimeout      time.Duration
	DBConnMaxLifetime  time.Duration

	// Logging
	LogLevel string
}

func Load() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Println("Warning: JWT_SECRET not set, using default (not secure for production)")
		jwtSecret = "your-secret-key-change-in-production"
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getFirstEnv([]string{"ENVIRONMENT", "NODE_ENV"}, "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// JWT (support both JWT_EXPIRY and JWT_EXPIRES_IN)
		JWTSecret:        jwtSecret,
		JWTExpiry:        getFirstEnv([]string{"JWT_EXPIRY", "JWT_EXPIRES_IN"}, "24h"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
		JWTRefreshExpiry: getFirstEnv([]string{"JWT_REFRESH_EXPIRY", "JWT_REFRESH_EXPIRES_IN"}, "48h"),

		// DB pool (support common .env keys)
		DBMaxPoolSize:     getEnvInt("DB_MAX_POOL_SIZE", 25),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBIdleTimeout:     getEnvDurationMS("DB_IDLE_TIMEOUT_MS", 0),
		DBConnMaxLifetime: getEnvDurationMS("DB_CONN_MAX_LIFETIME_MS", 5*60*1000), // default 5 minutes

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getFirstEnv(keys []string, defaultValue string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("Warning: %s is not a valid int (%q); using default %d", key, raw, defaultValue)
		return defaultValue
	}
	return v
}

func getEnvDurationMS(key string, defaultValueMS int) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return time.Duration(defaultValueMS) * time.Millisecond
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("Warning: %s is not a valid ms int (%q); using default %dms", key, raw, defaultValueMS)
		return time.Duration(defaultValueMS) * time.Millisecond
	}
	return time.Duration(v) * time.Millisecond
}
