package middleware

import (
	"[repo]/news_article_app/models"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// validateJWT parses and validates a JWT token, returning user claims
func validateJWT(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Trace().Err(err).Msg("Invalid JWT token")
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Trace().Msg("Invalid JWT claims")
		return nil, err
	}

	// Extract user details
	user := models.User{}

	if idClaim, ok := claims["id"]; ok {
		if idFloat, ok := idClaim.(float64); ok {
			user.ID = int(idFloat) // Convert float64 to int
		}
	}

	if usernameClaim, ok := claims["username"]; ok {
		if username, ok := usernameClaim.(string); ok {
			user.Name = username
		}
	}

	if roleClaim, ok := claims["role"]; ok {
		if role, ok := roleClaim.(string); ok {
			user.Role = role
		}
	}

	return &user, nil
}
