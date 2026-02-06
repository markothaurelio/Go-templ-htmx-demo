/* repositories/postgres_article_repository.go */
package repositories

import (
	"database/sql"
	"errors"

	"[repo]/news_article_app/models"
)

type PostgresArticleRepository struct {
	DB *sql.DB
}

func NewPostgresArticleRepository(db *sql.DB) *PostgresArticleRepository {
	return &PostgresArticleRepository{DB: db}
}

func (r *PostgresArticleRepository) CreateArticle(article models.Article) (int, error) {
	query := `INSERT INTO articles (title, content, author_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	var id int
	err := r.DB.QueryRow(query, article.Title, article.Content, article.AuthorID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresArticleRepository) GetArticleByID(id int) (models.Article, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at FROM articles WHERE id = $1`
	var article models.Article
	err := r.DB.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.AuthorID,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Article{}, errors.New("article not found")
		}
		return models.Article{}, err
	}
	return article, nil
}

func (r *PostgresArticleRepository) GetAllArticles() ([]models.Article, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at FROM articles`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.AuthorID,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (r *PostgresArticleRepository) UpdateArticle(article models.Article) error {
	query := `UPDATE articles SET title = $1, content = $2, updated_at = NOW() WHERE id = $3`
	res, err := r.DB.Exec(query, article.Title, article.Content, article.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("article not found")
	}

	return nil
}

func (r *PostgresArticleRepository) DeleteArticle(id int) error {
	query := `DELETE FROM articles WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("article not found")
	}

	return nil
}
