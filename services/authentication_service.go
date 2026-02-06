package services

import (
	"net/http"

	"[repo]/news_article_app/authentication"
	"github.com/rs/zerolog/log"
)

type BaseAuthService struct {
	AuthSys authentication.BaseAuthSystem
}

type CodeAuthService struct {
	AuthSys authentication.CodeAuthSystem
}

// Constructor for AuthService
func NewBaseAuthService(authsys authentication.BaseAuthSystem) *BaseAuthService {
	return &BaseAuthService{AuthSys: authsys}
}

// Constructor for AuthService
func NewCodeAuthService(authsys authentication.CodeAuthSystem) *CodeAuthService {
	return &CodeAuthService{AuthSys: authsys}
}

// Login method for AuthService
func (s *BaseAuthService) Login(w http.ResponseWriter, r *http.Request) (success bool) {
	success = s.AuthSys.Login(w, r) // No return needed

	log.Info().
		Str("ip", r.RemoteAddr).
		Str("username", r.FormValue("username")).
		Bool("success", success).
		Msg("Login Attempt")

	return success

}

// Logout method for AuthService
func (s *BaseAuthService) Logout(w http.ResponseWriter, r *http.Request) {
	s.AuthSys.Logout(w, r) // No return needed
}

// Login method for AuthService
func (s *CodeAuthService) SendCode(w http.ResponseWriter, r *http.Request) {
	s.AuthSys.SendCode(w, r) // No return needed
}

// Logout method for AuthService
func (s *CodeAuthService) VerifyCode(w http.ResponseWriter, r *http.Request) {
	success := s.AuthSys.VerifyCode(w, r) // No return needed

	email, err := r.Cookie("email")

	if err != nil {
		log.Error().AnErr("err", err).Msg("code verification")
	}

	log.Info().
		Str("ip", r.RemoteAddr).
		Str("email", email.Value).
		Bool("success", success).
		Msg("Login Attempt")

}
