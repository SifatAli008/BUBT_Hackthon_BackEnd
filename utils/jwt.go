package utils

import (
	"errors"
	"foodlink_backend/config"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret []byte

// InitJWT initializes JWT with secret from config
func InitJWT(cfg *config.Config) {
	jwtSecret = []byte(cfg.JWTSecret)
}

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID uuid.UUID, email, role string, expiry time.Duration) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not initialized")
	}

	expirationTime := time.Now().Add(expiry)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
			Issuer:    "foodlink-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ParseExpiry parses expiry duration from string.
// Supports Go durations like "24h", "1h30m" and also day/week shorthands like "1d", "2d", "1w".
func ParseExpiry(expiryStr string) (time.Duration, error) {
	s := strings.TrimSpace(expiryStr)
	if s == "" {
		return 0, errors.New("empty expiry")
	}

	// Support "Nd" / "Nw" (Go's time.ParseDuration does not support days/weeks)
	if strings.HasSuffix(s, "d") || strings.HasSuffix(s, "w") {
		unit := s[len(s)-1:]
		numStr := strings.TrimSpace(s[:len(s)-1])
		n, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			switch unit {
			case "d":
				return time.Duration(n * float64(24*time.Hour)), nil
			case "w":
				return time.Duration(n * float64(7*24*time.Hour)), nil
			}
		}
		// Fall through to the standard parser on parse failure
	}

	return time.ParseDuration(s)
}
