package middleware

import (
	"context"
	"net/http"
	"os"

	"[repo]/news_article_app/models"
	"github.com/rs/zerolog/log"
)

// Load JWT secret key from environment variable, fail early if not set
var jwtKey = []byte(os.Getenv("jwtkey")) // Ensure jwtKey is []byte

// assigned to all unauthenticated users of the web application
var null_user = models.User{
	ID:   -1,
	Name: "null",
	Role: "null",
}

// AuthMiddleware validates JWT from cookies and injects user data into the request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the auth_token cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			log.Trace().Msg("No auth_token cookie found")
			ctx := context.WithValue(r.Context(), "authenticated", false)
			ctx = context.WithValue(ctx, "user", null_user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Validate JWT and get user data
		user, err := validateJWT(cookie.Value)
		if err != nil {
			ctx := context.WithValue(r.Context(), "authenticated", false)
			ctx = context.WithValue(ctx, "user", null_user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		log.Debug().Str("username", user.Name).Str("role", user.Role).Msg("User authenticated")

		// Update context with user details
		ctx := context.WithValue(r.Context(), "authenticated", true)
		ctx = context.WithValue(ctx, "user", *user)

		// Proceed to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
