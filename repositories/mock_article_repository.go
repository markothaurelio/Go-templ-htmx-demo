/* repositories/mock_article_repository.go */
package repositories

import (
	"[repo]/news_article_app/models"
)

type MockArticleRepository struct {
	Articles []models.Article
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{}
}

func (m *MockArticleRepository) CreateArticle(article models.Article) (int, error) {
	article.ID = len(m.Articles) + 1
	m.Articles = append(m.Articles, article)
	return article.ID, nil
}

func (m *MockArticleRepository) GetArticleByID(id int) (models.Article, error) {
	for _, article := range m.Articles {
		if article.ID == id {
			return article, nil
		}
	}
	return models.Article{}, nil
}

func (m *MockArticleRepository) GetAllArticles() ([]models.Article, error) {
	return m.Articles, nil
}

func (m *MockArticleRepository) UpdateArticle(article models.Article) error {
	for i, a := range m.Articles {
		if a.ID == article.ID {
			m.Articles[i] = article
			return nil
		}
	}
	return nil
}

func (m *MockArticleRepository) DeleteArticle(id int) error {
	for i, a := range m.Articles {
		if a.ID == id {
			m.Articles = append(m.Articles[:i], m.Articles[i+1:]...)
			return nil
		}
	}
	return nil
}
