package main

import (
	"net/http"
	"os"

	"[repo]/news_article_app/appdata"
	"[repo]/news_article_app/authentication"
	"[repo]/news_article_app/handlers"
	"[repo]/news_article_app/middleware"
	"[repo]/news_article_app/models"
	"[repo]/news_article_app/postgres_db"
	"[repo]/news_article_app/repositories"
	"[repo]/news_article_app/services"
	"[repo]/news_article_app/templates"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GeoJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Load JSON data from the file
	data, err := os.ReadFile("states.geojson")
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	// Set JSON content type and write the data
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {

	// Initialize
	initialize_logging()

	// Connect to PostgreSQL
	db := postgres_db.ConnectDB()
	defer postgres_db.CloseDB()

	// Register article repo with service and service wth handlers
	articleRepo := repositories.NewPostgresArticleRepository(db)
	articleService := services.NewArticleService(articleRepo)
	article_handler := handlers.NewArticleHandler(articleService)

	baseAuthSys := authentication.NewAuthSys(db)
	baseAuthService := services.NewBaseAuthService(baseAuthSys)
	baseAuthHandler := handlers.NewBaseAuthHandler(baseAuthService)

	// NAV HANDLERS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Layout(templates.Index(), "Index Page").Render(r.Context(), w)
	})

	// Data endpoints (renamed from /geojson/* to /data/*)
	for path, data := range appdata.GeoJSONDataSources {
		http.HandleFunc("/data/"+path, data)
	}

	// Renamed from /analyses to /articles
	http.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {

		articles, err := articleService.GetAllArticles()
		if err != nil {
			log.Error().Err(err).Msg("Could not get articles")
		}

		templates.Layout(templates.ArticleList(articles), "Article Feed").Render(r.Context(), w)
	})

	http.HandleFunc("/maps", func(w http.ResponseWriter, r *http.Request) {
		templates.Layout(templates.MainMapPage(), "").Render(r.Context(), w)
	})

	// Register map handlers dynamically using a loop
	// Renamed path suffix from "-demographics" to "-dataset"
	for category, yearData := range appdata.MapCategories {
		http.HandleFunc("/maps/"+category+"-dataset", handlers.MapHandler(category, yearData))
	}

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		templates.Layout(templates.About(), "About").Render(r.Context(), w)
	})

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {

		if r.Context().Value("authenticated") == false {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		templates.Layout(templates.ProfilePage(r.Context().Value("user").(models.User)), "Profile").Render(r.Context(), w)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {

		if r.Context().Value("authenticated") == false {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		if r.Method == http.MethodGet {
			templates.Layout(templates.ArticleSubmission(), "Submit").Render(r.Context(), w)
		} else if r.Method == http.MethodPost {
			article_handler.CreateArticleHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			templates.Layout(templates.LoginForm(""), "Login").Render(r.Context(), w)
		} else if r.Method == http.MethodPost {
			baseAuthHandler.LoginHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		baseAuthHandler.LogoutHandler(w, r)
	})

	// Renamed article routes from /analyses/{id} to /articles/{id}
	http.HandleFunc("/articles/{id}", article_handler.GetArticleByID)

	// Renamed delete routes accordingly
	http.HandleFunc("/articles/{id}/delete", article_handler.DeleteArticleByID)

	// Serve static files from the "static" folder
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Info().Str("port", ":3000").Msg("starting application")
	http.ListenAndServe("0.0.0.0:3000", middleware.GetAllGlobalMiddleware())
}

func initialize_logging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	level := "debug"
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
