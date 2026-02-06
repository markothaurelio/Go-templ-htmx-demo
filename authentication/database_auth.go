package authentication

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"[repo]/news_article_app/templates"
	"github.com/rs/zerolog/log"

	"golang.org/x/crypto/bcrypt"
)

// AuthSys manages authentication with PostgreSQL
type AuthSys struct {
	db *sql.DB
}

// NewAuthSys initializes the authentication system with the database
func NewAuthSys(db *sql.DB) *AuthSys {
	return &AuthSys{db: db}
}

// Login authenticates a user and sets a JWT token in a cookie
func (a *AuthSys) Login(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return false
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		templates.LoginForm("Invalid form submission").Render(r.Context(), w)
		return false
	}

	// Get user input
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Retrieve the user from the database
	var (
		userID     int
		role       string
		storedHash string
	)

	err := a.db.QueryRowContext(
		context.Background(),
		"SELECT id, role, password_hash FROM users WHERE username = $1",
		username,
	).Scan(&userID, &role, &storedHash)

	if err == sql.ErrNoRows {
		templates.LoginForm("Invalid credentials. Please try again.").Render(r.Context(), w)
		return false
	}
	if err != nil {
		log.Error().Err(err).Msg("Database query error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return false
	}

	// NEVER log passwords or password hashes.
	// Compare submitted password with stored bcrypt hash.
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		templates.LoginForm("Invalid credentials. Please try again.").Render(r.Context(), w)
		return false
	}

	// Generate JWT
	token, err := generateJWT(userID, username, role)
	if err != nil {
		log.Error().Err(err).Msg("Error generating JWT")
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return false
	}

	// Set JWT in a cookie
	// In production you want Secure=true (HTTPS) and SameSite=Lax/Strict.
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // safer default; requires HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	log.Info().Int("user_id", userID).Str("user", username).Msg("User logged in successfully")
	return true
}

// Logout clears the authentication cookie
func (a *AuthSys) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	log.Info().Msg("User logged out successfully")
}
