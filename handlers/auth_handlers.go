/* handlers/standard_handlers.go */

package handlers

import (
	"net/http"
	"time"

	"[repo]/news_article_app/services"
	"[repo]/news_article_app/templates"
)

type BaseAuthHandler struct {
	Service *services.BaseAuthService
}

type CodeAuthHandler struct {
	Service *services.CodeAuthService
}

func NewBaseAuthHandler(service *services.BaseAuthService) *BaseAuthHandler {
	return &BaseAuthHandler{Service: service}
}

func NewCodeAuthHandler(service *services.CodeAuthService) *CodeAuthHandler {
	return &CodeAuthHandler{Service: service}
}

func (h *BaseAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	success := h.Service.Login(w, r)

	if success == true {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

}

func (h *BaseAuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	h.Service.Logout(w, r)

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (h *CodeAuthHandler) EmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Example email (you'll likely retrieve this from the request in practice)
	email := r.FormValue("email")

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "email",                        // Name of the cookie
		Value:    email,                          // Email value
		Path:     "/",                            // Cookie applies to the entire domain
		Expires:  time.Now().Add(24 * time.Hour), // Expiry time (24 hours here)
		HttpOnly: true,                           // Prevent JavaScript from accessing the cookie
		Secure:   true,                           // Set to true in production (requires HTTPS)
	})

	h.Service.SendCode(w, r)
	templates.VerifyLoginForm("").Render(r.Context(), w) // Serve the form
}

func (h *CodeAuthHandler) VerifyLoginHandler(w http.ResponseWriter, r *http.Request) {

	// Use the email value in your logic
	h.Service.VerifyCode(w, r)

}
