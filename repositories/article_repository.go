/* repositories/article_repository.go */
package repositories

import "[repo]/news_article_app/models"

type ArticleRepository interface {
	CreateArticle(article models.Article) (int, error)
	GetArticleByID(id int) (models.Article, error)
	GetAllArticles() ([]models.Article, error)
	UpdateArticle(article models.Article) error
	DeleteArticle(id int) error
}
