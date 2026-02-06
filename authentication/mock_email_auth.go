package authentication

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"[repo]/news_article_app/mailer"
	"[repo]/news_article_app/repositories"
	"github.com/rs/zerolog/log"
)

type MockCodeAuthSys struct {
	Repo   *repositories.MockStorage
	Mailer mailer.GomailMailer
}

func NewMockCodeAuthSys(repo *repositories.MockStorage, mailer *mailer.GomailMailer) *MockCodeAuthSys {
	return &MockCodeAuthSys{Repo: repo, Mailer: *mailer}
}

func (h *MockCodeAuthSys) SendCode(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")

	// Generate a 6-digit random code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store the code in the database
	h.Repo.SaveLoginCode(email, code, time.Now().Add(time.Minute*20))

	// Send the code via email
	err := h.Mailer.SendEmail(email, "Your Login Code", fmt.Sprintf("Your login code is: %s", code))

	if err != nil {
		log.Error().AnErr("err", err).Msg("mailer")
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MockCodeAuthSys) VerifyCode(w http.ResponseWriter, r *http.Request) bool {

	// Get the email cookie
	cookie, err := r.Cookie("email")
	if err != nil {
		// Handle error if the cookie is not present
		http.Error(w, "Email cookie not found", http.StatusUnauthorized)
		return false
	}

	cookieEmail := cookie.Value // Extract the email value from the cookie

	code := r.FormValue("code")

	login, yes := h.Repo.GetLoginCode(code)

	log.Debug().Str("email_1", cookieEmail).Str("email_2", login.Email)

	if cookieEmail == login.Email && yes {

		// Generate JWT
		token, err := generateJWT(0, login.Email, "nan") // TODO this needs to be reworked
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
	} else {
		http.ResponseWriter.Write(w, []byte("Invalid code"))
		return false
	}

}
