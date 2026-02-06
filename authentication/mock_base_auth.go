package authentication

import (
	"net/http"
	"time"

	"[repo]/news_article_app/templates"
	"github.com/rs/zerolog/log"
)

type MockAuthSys struct {
}

func NewMockBaseAuthSys() *MockAuthSys {
	return &MockAuthSys{}
}

func (a MockAuthSys) Login(w http.ResponseWriter, r *http.Request) bool {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return false
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		templates.LoginForm("Invalid form submission").Render(r.Context(), w)
	}

	// Get username and password from form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Simple hardcoded authentication (PLACEHOLDER VALUES - configure in your local/dev env only)
	if username != "YOUR_DEV_USERNAME" || password != "YOUR_DEV_PASSWORD" {
		templates.LoginForm("Invalid credentials. Please try again.").Render(r.Context(), w)
		return false
	}
	// Generate JWT
	token, err := generateJWT(0, username, "nan") // TODO this needs to be reworked
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return false
	}

	// Set JWT as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,  // Prevent JavaScript access
		Secure:   false, // Set true in production (requires HTTPS)
		Path:     "/",
	})

	return true

}

func (a MockAuthSys) Logout(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("do this")
}
