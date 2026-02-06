/* repositories/login_repository.go */
package repositories

import "[repo]/news_article_app/models"

type LoginRepository interface {
	CreateUser(user models.User) (int, error)
	DeleteUser(username string)
}
