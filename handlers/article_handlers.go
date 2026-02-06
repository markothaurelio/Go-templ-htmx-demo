/* handlers/article_handle.go */
package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"[repo]/news_article_app/models"
	"[repo]/news_article_app/services"
	"[repo]/news_article_app/templates"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type ArticleHandler struct {
	Service *services.ArticleService
}

func NewArticleHandler(service *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{Service: service}
}

func convertMarkdownToHTML(markdown string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		panic(err)
	}

	return buf.String()
}

func (h *ArticleHandler) CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	markdownContent := r.FormValue("content")

	if title == "" || markdownContent == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Convert Markdown to HTML with proper formatting
	htmlContent := convertMarkdownToHTML(markdownContent)

	fmt.Println("Converted HTML Output:", htmlContent) // Debugging: Check final HTML

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	article := models.Article{
		Title:     title,
		Content:   htmlContent, // Store only sanitized HTML
		AuthorID:  user.ID,
		CreatedAt: time.Now(),
	}

	id, err := h.Service.CreateArticle(article)
	if err != nil {
		log.Error().Err(err).Msg("Error creating article")
		http.Error(w, "Failed to create article", http.StatusInternalServerError)
		return
	}

	log.Info().Int("article ID", id).Msg("Article created")

	templates.ArticlePage(article).Render(r.Context(), w)
}

func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	id_int, err := strconv.Atoi(id)

	article, err := h.Service.GetArticleByID(id_int)
	if err != nil || article.AuthorID == 0 {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	templates.Layout(templates.ArticlePage(article), "Article").Render(r.Context(), w)

}

func (h *ArticleHandler) DeleteArticleByID(w http.ResponseWriter, r *http.Request) {

	if r.Context().Value("authenticated") != true && r.Context().Value("user").(models.User).Role != "admin" {
		return
	}

	id := r.PathValue("id")

	id_int, err := strconv.Atoi(id)

	err = h.Service.DeleteArticle(id_int)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

}
