package authentication

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// Get JWT secret key from environment variable
var jwtKey = []byte(os.Getenv("jwtkey")) // Ensure jwtKey is []byte

// GenerateJWT creates a JWT token for a user
func generateJWT(userID int, username, role string) (string, error) {

	log.Debug().Bytes("key", jwtKey).Msg("jwt")

	claims := jwt.MapClaims{
		"id":       userID, // User's unique identifier
		"username": username,
		"role":     role,
		"iss":      "news-article-app",
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),                     // Issued at timestamp
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with the secret key
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}
