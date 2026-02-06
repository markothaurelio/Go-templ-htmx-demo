/* services/article_service.go */
package services

import (
	"[repo]/news_article_app/models"
	"[repo]/news_article_app/repositories"
)

type ArticleService struct {
	Repo repositories.ArticleRepository
}

func NewArticleService(repo repositories.ArticleRepository) *ArticleService {
	return &ArticleService{Repo: repo}
}

func (s *ArticleService) CreateArticle(article models.Article) (int, error) {
	return s.Repo.CreateArticle(article)
}

func (s *ArticleService) GetArticleByID(id int) (models.Article, error) {
	return s.Repo.GetArticleByID(id)
}

func (s *ArticleService) GetAllArticles() ([]models.Article, error) {
	return s.Repo.GetAllArticles()
}

func (s *ArticleService) UpdateArticle(article models.Article) error {
	return s.Repo.UpdateArticle(article)
}

func (s *ArticleService) DeleteArticle(id int) error {
	return s.Repo.DeleteArticle(id)
}
