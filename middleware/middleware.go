package middleware

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	visitorStore = make(map[string]bool)
	visitorLock  sync.Mutex
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("method", r.Method).
			Str("path", r.URL.Path).
			Bool("authenticated", r.Context().Value("authenticated").(bool)).
			Any("user", r.Context().Value("user")).
			Msg("user request")

		next.ServeHTTP(w, r) // Call the next handler
	})
}

func securityMiddleware(next http.Handler) http.Handler { // TODO work out what to do here
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("security stuff")
		next.ServeHTTP(w, r)
	})
}

func chainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func uniqueVisitorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("visitor_id")

		visitorLock.Lock()
		defer visitorLock.Unlock()

		if err != nil || cookie.Value == "" {
			// New visitor
			newID := uuid.NewString()
			http.SetCookie(w, &http.Cookie{
				Name:     "visitor_id",
				Value:    newID,
				Path:     "/",
				MaxAge:   86400 * 365,          // 1 year
				Secure:   true,                 // 🔐 Ensures cookie is only sent over HTTPS
				HttpOnly: true,                 // Optional: prevents access via JavaScript (recommended for security)
				SameSite: http.SameSiteLaxMode, // Optional: prevents some CSRF issues
			})

			visitorStore[newID] = true

			log.Info().
				Int("unique_visitor_count", len(visitorStore)).
				Str("visitor_id", newID).
				Msg("New unique visitor")
		} else {
			// Existing visitor
			if _, exists := visitorStore[cookie.Value]; !exists {
				visitorStore[cookie.Value] = true

				log.Info().
					Int("unique_visitor_count", len(visitorStore)).
					Str("visitor_id", cookie.Value).
					Msg("Returning visitor counted as unique (not seen before)")
			}
		}

		next.ServeHTTP(w, r)
	})
}

func GetAllGlobalMiddleware() http.Handler {
	// Chain middleware
	return chainMiddleware(
		http.DefaultServeMux,
		loggingMiddleware,
		uniqueVisitorMiddleware,
		AuthMiddleware)
}
